package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/astaxie/beego/logs"
)

type Remote struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Protocol string `json:"protocol"`
	Auth     bool   `json:"auth"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Config struct {
	Address    string   `json:"address"`
	Port       int      `json:"port"`
	Mode       string   `json:"mode"`
	TestUrl    string   `json:"testUrl"`
	RemoteName string   `json:"remote"`
	RemoteList []Remote `json:"remoteList"`
	DomainList []string `json:"domainList"`
}

var configCache = Config{
	Address:    "0.0.0.0",
	Port:       8080,
	Mode:       "local",
	TestUrl:    "https://google.com/",
	RemoteName: "",
	DomainList: make([]string, 0),
	RemoteList: make([]Remote, 0),
}

var configFilePath string
var configLock sync.Mutex

func configSyncToFile() error {
	configLock.Lock()
	defer configLock.Unlock()

	value, err := json.MarshalIndent(configCache, "\t", " ")
	if err != nil {
		logs.Error("json marshal config fail, %s", err.Error())
		return err
	}
	return SaveToFile(configFilePath, value)
}

func ConfigGet() *Config {
	return &configCache
}

func DomainListSave(domain []string) error {
	configCache.DomainList = domain
	return configSyncToFile()
}

func RemoteListSave(remote []Remote) error {
	configCache.RemoteList = remote
	return configSyncToFile()
}

func ModeSave(mode string) error {
	configCache.Mode = mode
	return configSyncToFile()
}

func TestUrlSave(test string) error {
	configCache.TestUrl = test
	return configSyncToFile()
}

func RemoteSave(remote string) error {
	configCache.RemoteName = remote
	return configSyncToFile()
}

func ListenAddressSave(addr string) error {
	configCache.Address = addr
	return configSyncToFile()
}

func ListenPortSave(port int) error {
	configCache.Port = port
	return configSyncToFile()
}

func ConfigInit() error {
	configFilePath = fmt.Sprintf("%s%c%s", ConfigDirGet(), os.PathSeparator, "config.json")

	_, err := os.Stat(configFilePath)
	if err != nil {
		value, err := json.Marshal(configCache)
		if err != nil {
			logs.Error("json marshal config fail, %s", err.Error())
			return err
		}

		err = SaveToFile(configFilePath, value)
		if err != nil {
			logs.Error("default config save to app data dir fail, %s", err.Error())
			return err
		}
	}

	value, err := os.ReadFile(configFilePath)
	if err != nil {
		logs.Error("read config file from app data dir fail, %s", err.Error())
		return err
	}

	err = json.Unmarshal(value, &configCache)
	if err != nil {
		logs.Error("json unmarshal config fail, %s", err.Error())
		return err
	}

	return nil
}
