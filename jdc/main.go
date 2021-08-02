package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"

	"github.com/beego/beego/v2/server/web"
	"github.com/cdle/jd_study/jdc/controllers"
	"github.com/cdle/jd_study/jdc/models"
)

var qrcode = ""
var name = "jdc"

func main() {
	go func() {
		models.Save <- &models.JdCookie{}
	}()
	web.Get("/", func(ctx *context.Context) {
		if models.Config.Qrcode == "" {
			models.Config.Qrcode = "https://ghproxy.com/https://raw.githubusercontent.com/cdle/jd_study/main/jdc/theme/bidong.html"
		}
		if qrcode != "" {
			ctx.WriteString(qrcode)
			return
		}
		if strings.Contains(models.Config.Qrcode, "http") {
			logs.Info("下载最新主题")
			s, _ := httplib.Get(models.Config.Qrcode).String()
			qrcode = s
			ctx.WriteString(s)
			return
		} else {
			f, err := os.Open(models.Config.Qrcode)
			if err == nil {
				d, _ := ioutil.ReadAll(f)
				qrcode = string(d)
				ctx.WriteString(string(d))
				return
			}
		}
	})
	web.Router("/api/login/qrcode", &controllers.LoginController{}, "get:GetQrcode")
	web.Router("/api/login/query", &controllers.LoginController{}, "get:Query")
	web.Router("/api/account", &controllers.AccountController{}, "get:List")
	web.Router("/api/account", &controllers.AccountController{}, "post:CreateOrUpdate")
	web.Router("/admin", &controllers.AccountController{}, "get:Admin")
	if models.Config.Static == "" {
		models.Config.Static = "./static"
	}
	web.BConfig.WebConfig.StaticDir["/static"] = models.Config.Static
	web.BConfig.AppName = name
	web.BConfig.WebConfig.AutoRender = false
	web.BConfig.CopyRequestBody = true
	web.BConfig.WebConfig.Session.SessionOn = true
	web.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600
	web.BConfig.WebConfig.Session.SessionName = name
	// go func() {
	// 	time.Sleep(time.Second)
	// 	killp()
	// }()
	web.Run()
}
