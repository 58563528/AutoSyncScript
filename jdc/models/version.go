package models

import (
	"fmt"
	"regexp"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
)

var version = "2021080201"

func initVersion() {
	logs.Info("检查更新" + version)
	value, err := httplib.Get("https://ghproxy.com/https://raw.githubusercontent.com/cdle/jd_study/main/jdc/models/version.go").String()
	if err != nil {
		logs.Info("更新User-Agent失败")
	} else {
		if match := regexp.MustCompile(`var version = "(\d{10})"`).FindStringSubmatch(value); len(match) != 0 {
			fmt.Println(match)
		}
	}
}
