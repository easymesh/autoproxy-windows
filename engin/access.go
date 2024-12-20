package engin

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/astaxie/beego/logs"
)

func PublicFailDelay() {
	time.Sleep(5 * time.Second) // 防DOS攻击延时
}

type StatInfo struct {
	ForwardSize int64
	RequestCnt  int64
	SessionCnt  int64
}

type HttpAccess struct {
	Timeout    int
	Address    string
	httpserver *http.Server
	sync.WaitGroup

	stat           StatInfo
	authHandler    func(auth *AuthInfo) bool
	forwardHandler func(address string, r *http.Request) Forward
	defaultForward Forward
}

type Access interface {
	StatGet() StatInfo
	Shutdown() error
	AuthHandlerSet(func(*AuthInfo) bool)
	ForwardHandlerSet(func(address string, r *http.Request) Forward)
}

func AuthFailHandler(w http.ResponseWriter, r *http.Request) {
	PublicFailDelay()
	logs.Warn("Request authentication failed. RemoteAddr: ", r.RemoteAddr)
	w.Header().Add("Proxy-Authenticate", "Basic realm=\"Access to internal site\"")
	http.Error(w,
		"Request authentication failed.",
		http.StatusProxyAuthRequired)
}

func AuthInfoParse(r *http.Request) *AuthInfo {
	value := r.Header.Get("Proxy-Authorization")
	if value == "" {
		return nil
	}
	body, err := base64.StdEncoding.DecodeString(value[6:])
	if err != nil {
		return nil
	}
	ctx := strings.Split(string(body), ":")
	if len(ctx) != 2 {
		return nil
	}
	return &AuthInfo{User: ctx[0], Token: ctx[1]}
}

func (acc *HttpAccess) NoProxyHandler(w http.ResponseWriter, r *http.Request) {
	logs.Warn("request is illegal. RemoteAddr: ", r.RemoteAddr)
	http.Error(w,
		"This is a proxy server. Does not respond to non-proxy requests.",
		http.StatusInternalServerError)
}

func (acc *HttpAccess) AuthHandlerSet(handler func(auth *AuthInfo) bool) {
	acc.authHandler = handler
}

func (acc *HttpAccess) ForwardHandlerSet(handler func(address string, r *http.Request) Forward) {
	acc.forwardHandler = handler
}

func (acc *HttpAccess) AuthHttp(r *http.Request) bool {
	if acc.authHandler == nil {
		return true
	}
	if AuthCache(r) == true {
		return true
	}
	auth := acc.authHandler(AuthInfoParse(r))
	if auth == true {
		AuthLogin(r)
	}
	return auth
}

func (acc *HttpAccess) StatGet() StatInfo {
	return acc.stat
}

func (acc *HttpAccess) Shutdown() error {
	context, cencel := context.WithTimeout(context.Background(), 15*time.Second)
	err := acc.httpserver.Shutdown(context)
	cencel()
	if err != nil {
		logs.Error("http access ready to shut down fail, %s", err.Error())
	}
	acc.Wait()
	return err
}

func DebugReqeust(r *http.Request) {
	var headers string
	for key, value := range r.Header {
		headers += fmt.Sprintf("[%s:%s]", key, value)
	}
	logs.Info("%s %s %s %s %s %s", r.RemoteAddr, r.Host, r.URL.Scheme, r.Method, r.URL.String(), headers)
}

func (acc *HttpAccess) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	DebugReqeust(r)

	atomic.AddInt64(&acc.stat.RequestCnt, 1)

	if acc.AuthHttp(r) == false {
		AuthFailHandler(w, r)
		return
	}

	if r.Method == "CONNECT" {
		acc.HttpsRoundTripper(w, r)
		return
	}

	var rsp *http.Response
	var err error

	if !r.URL.IsAbs() {
		r.URL.Host = r.Host
		r.URL.Scheme = "http"
		rsp, err = acc.defaultForward.Http(r)
	} else {
		removeProxyHeaders(r)
		rsp, err = acc.HttpRoundTripper(r)
	}

	if err != nil {
		errStr := fmt.Sprintf("transport %s %s failed! %s", r.Host, r.URL.String(), err.Error())
		logs.Warn(errStr)
		http.Error(w, errStr, http.StatusInternalServerError)
		return
	}

	if rsp == nil {
		errStr := fmt.Sprintf("transport %s read response failed!", r.URL.Host)
		logs.Warn(errStr)
		http.Error(w, errStr, http.StatusInternalServerError)
		return
	}

	atomic.AddInt64(&acc.stat.SessionCnt, 1)
	defer atomic.AddInt64(&acc.stat.SessionCnt, -1)

	copyHeaders(w.Header(), rsp.Header)
	w.WriteHeader(rsp.StatusCode)

	size, err := io.Copy(w, rsp.Body)
	if size == 0 && err != nil {
		logs.Warn("io copy fail", err.Error())
	}
	rsp.Body.Close()

	if size > 0 {
		atomic.AddInt64(&acc.stat.ForwardSize, size)
	}

}

func copyHeaders(dst, src http.Header) {
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}

func removeProxyHeaders(r *http.Request) {
	r.RequestURI = ""
	r.Header.Del("Proxy-Connection")
	r.Header.Del("Proxy-Authenticate")
	r.Header.Del("Proxy-Authorization")
}

func NewHttpsAccess(addr string, timeout int, tlsEnable bool, certfile, keyfile string) (Access, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logs.Error("listen address fail", addr)
		return nil, err
	}

	var config *tls.Config
	if tlsEnable {
		config, err = TlsConfigServer(certfile, keyfile)
		if err != nil {
			logs.Error("make tls config server fail, %s", err.Error())
			return nil, err
		}
		lis = tls.NewListener(lis, config)
	}

	acc := new(HttpAccess)
	acc.Address = addr
	acc.Timeout = timeout
	acc.defaultForward = NewDefault(timeout)

	tmout := time.Duration(timeout) * time.Second

	httpserver := &http.Server{
		Handler:      acc,
		ReadTimeout:  tmout,
		WriteTimeout: tmout,
		TLSConfig:    config,
	}

	acc.httpserver = httpserver

	acc.Add(1)

	go func() {
		defer acc.Done()
		err = httpserver.Serve(lis)
		if err != nil {
			logs.Error("http server ", err.Error())
		}
	}()

	if config == nil {
		logs.Info("access http start success.")
	} else {
		logs.Info("access https start success.")
	}

	return acc, nil
}
