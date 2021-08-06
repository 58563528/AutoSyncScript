package models

import (
	"io"
	"os"
	"regexp"
	"runtime"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
)

var version = "2021080601"

func initVersion() {
	logs.Info("检查更新" + version)
	value, err := httplib.Get("https://ghproxy.com/https://raw.githubusercontent.com/cdle/jd_study/main/xdd/models/version.go").String()
	if err != nil {
		logs.Info("更新User-Agent失败")
	} else {
		name := "jdc_" + runtime.GOOS + "_" + runtime.GOARCH
		if match := regexp.MustCompile(`var version = "(\d{10})"`).FindStringSubmatch(value); len(match) != 0 {
			if match[1] > version {
				logs.Warn("版本过低，下载更新")
				rsp, err := httplib.Get("https://ghproxy.com/https://github.com/cdle/jd_study/releases/download/main/" + name).Response()
				if err != nil {
					logs.Warn("无法下载更新")
					return
				}
				filename := ExecPath + "/" + name
				f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
				if err != nil {
					logs.Warn("无法打开文件" + filename)
					return
				}
				defer f.Close()
				_, err = io.Copy(f, rsp.Body)
				if err != nil {
					logs.Warn("下载更新失败")
					return
				}
				logs.Info("更新下载至"+filename, "请手动启动")
			}
		}
	}
}
