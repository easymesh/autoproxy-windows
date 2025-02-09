package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	mathrand "math/rand"

	"github.com/astaxie/beego/logs"
)

func SaveToFile(name string, body []byte) error {
	return os.WriteFile(name, body, 0664)
}

func GetToken(length int) string {
	token := make([]byte, length)
	bytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!#$%^&*"
	for i := 0; i < length; i++ {
		token[i] = bytes[mathrand.Int()%len(bytes)]
	}
	return string(token)
}

func GetUser(length int) string {
	token := make([]byte, length)
	bytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	for i := 0; i < length; i++ {
		token[i] = bytes[mathrand.Int()%len(bytes)]
	}
	return string(token)
}

func CapSignal(proc func()) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		proc()
		logs.Error("recv signcal %s, ready to exit", sig.String())
		os.Exit(-1)
	}()
}

func InterfaceAddsGet(iface *net.Interface) ([]net.IP, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	ips := make([]net.IP, 0)
	for _, v := range addrs {
		ipone, _, err := net.ParseCIDR(v.String())
		if err != nil {
			continue
		}
		if len(ipone) > 0 {
			ips = append(ips, ipone)
		}
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("interface not any address.")
	}
	return ips, nil
}

func ByteViewLite(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%db", size)
	} else if size < (1024 * 1024) {
		return fmt.Sprintf("%.1fKb", float64(size)/float64(1024))
	} else {
		return fmt.Sprintf("%.1fMb", float64(size)/float64(1024*1024))
	}
}

func ByteView(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	} else if size < (1024 * 1024) {
		return fmt.Sprintf("%.1fKB", float64(size)/float64(1024))
	} else if size < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.1fMB", float64(size)/float64(1024*1024))
	} else if size < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.1fGB", float64(size)/float64(1024*1024*1024))
	} else {
		return fmt.Sprintf("%.1fTB", float64(size)/float64(1024*1024*1024*1024))
	}
}

func StringList(list []string) string {
	var body string
	for idx, v := range list {
		if idx == len(list)-1 {
			body += fmt.Sprintf("%s", v)
		} else {
			body += fmt.Sprintf("%s;", v)
		}
	}
	return body
}

func StringDiff(oldlist []string, newlist []string) ([]string, []string) {
	del := make([]string, 0)
	add := make([]string, 0)
	for _, v1 := range oldlist {
		flag := false
		for _, v2 := range newlist {
			if v1 == v2 {
				flag = true
				break
			}
		}
		if !flag {
			del = append(del, v1)
		}
	}
	for _, v1 := range newlist {
		flag := false
		for _, v2 := range oldlist {
			if v1 == v2 {
				flag = true
				break
			}
		}
		if !flag {
			add = append(add, v1)
		}
	}
	return del, add
}

func StringClone(list []string) []string {
	output := make([]string, len(list))
	copy(output, list)
	return output
}
