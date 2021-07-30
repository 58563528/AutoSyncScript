package models

import (
	"io/ioutil"

	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/yaml.v2"
)

type Container struct {
	Type      string
	Name      string
	Default   bool
	Address   string
	Username  string
	Password  string
	Path      string
	Version   string
	Token     string
	Available bool
	Delete    []string
	Weigth    int
	First     bool
}
type Yaml struct {
	Containers []Container
	Qrcode     string
	Master     string
	Mode       string
	Static     string
}

var Balance = "balance"
var Parallel = "parallel"

var Config Yaml

func initConfig() {
	content, err := ioutil.ReadFile(ExecPath + "/conf/config.yaml")
	if err != nil {
		logs.Warn("解析config.yaml读取错误: %v", err)
	}
	if yaml.Unmarshal(content, &Config) != nil {
		logs.Warn("解析config.yaml出错: %v", err)
	}
	if Config.Master == "" {
		Config.Master = "xxxx"
	}
	if Config.Mode != Parallel {
		Config.Mode = Balance
	}
}
