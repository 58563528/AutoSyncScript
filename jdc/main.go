package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/server/web/context"

	"github.com/beego/beego/v2/server/web"
	"github.com/cdle/jd_study/jdc/controllers"
	"github.com/cdle/jd_study/jdc/models"
)

var help = "-p 运行端口\n-qla 青龙登录地址\n-qlu 青龙登录用户名\n-qlp 青龙登录密码\n-v4 配置文件路径"

func main() {
	l := len(os.Args)
	if l == 0 {
		fmt.Println(help)
		return
	}
	for i, arg := range os.Args {
		if i+1 <= l-1 {
			v := os.Args[i+1]
			switch arg {
			case "-h":
				fmt.Println(help)
				return
			case "-p":
				p, _ := strconv.Atoi(v)
				web.BConfig.Listen.HTTPPort = p
			case "-qla":
				vv := regexp.MustCompile(`^(https?://[\.\w]+:?\d*)`).FindStringSubmatch(v)
				if len(vv) == 2 {
					models.QlAddress = vv[1]
				}
			case "-qlu":
				models.QlUserName = v
			case "-qlp":
				models.QlPassword = v
			case "-qlv":
				models.QlVersion = v
			case "-v4":
				models.V4Config = v
			}
		}
	}
	if models.V4Config != "" {
		f, err := os.Open(models.V4Config)
		if err != nil {
			logs.Warn("无法打开V4配置文件，请检查路径是否正确")
			return
		}
		models.V4Handle(&models.JdCookie{})
		f.Close()
	} else {
		if models.QlAddress == "" {
			logs.Warn("未指定青龙登录地址")
			return
		}
		if models.QlUserName == "" {
			logs.Warn("未指定青龙登录用户名")
			return
		}
		if models.QlPassword == "" {
			logs.Warn("未指定青龙登录密码")
			return
		}
		if models.GetToken(); models.Token == "" {
			logs.Warn("JDC无法与青龙面板取得联系，请检查账号")
			return
		} else {
			logs.Info("JDC成功接入青龙")
		}
	}

	web.Get("/", func(ctx *context.Context) {
		ctx.WriteString(models.Qrocde)
	})
	web.Router("/api/login/qrcode", &controllers.LoginController{}, "get:GetQrcode")
	web.Router("/api/login/query", &controllers.LoginController{}, "get:Query")
	web.Router("/api/account", &controllers.AccountController{}, "get:List")
	web.Router("/api/account", &controllers.AccountController{}, "post:CreateOrUpdate")
	web.BConfig.AppName = "jdc"
	web.BConfig.WebConfig.AutoRender = false
	web.BConfig.CopyRequestBody = true
	web.BConfig.WebConfig.Session.SessionOn = true
	web.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600
	web.BConfig.WebConfig.Session.SessionName = "jdc"
	web.Run()
}
