package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/astaxie/beego/httplib"
	"github.com/beego/beego/v2/server/web/context"

	"github.com/beego/beego/v2/server/web"
	"github.com/cdle/jd_study/jdc/controllers"
	"github.com/cdle/jd_study/jdc/models"
)

func main() {
	models.Save <- &models.JdCookie{}
	web.Get("/", func(ctx *context.Context) {
		if models.Config.Qrcode != "" {
			if strings.Contains(models.Config.Qrcode, "http") {
				s, _ := httplib.Get(models.Config.Qrcode).String()
				ctx.WriteString(s)
				return
			} else {
				f, err := os.Open(models.Config.Qrcode)
				if err == nil {
					d, _ := ioutil.ReadAll(f)
					ctx.WriteString(string(d))
					return
				}
			}
		}
		ctx.WriteString(models.Qrocde)
	})
	web.Router("/api/login/qrcode", &controllers.LoginController{}, "get:GetQrcode")
	web.Router("/api/login/query", &controllers.LoginController{}, "get:Query")
	web.Router("/api/account", &controllers.AccountController{}, "get:List")
	web.Router("/api/account", &controllers.AccountController{}, "post:CreateOrUpdate")
	web.Router("/admin", &controllers.AccountController{}, "get:Admin")
	web.BConfig.AppName = "jdc"
	web.BConfig.WebConfig.AutoRender = false
	web.BConfig.CopyRequestBody = true
	web.BConfig.WebConfig.Session.SessionOn = true
	web.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600
	web.BConfig.WebConfig.Session.SessionName = "jdc"
	web.Run()
}
