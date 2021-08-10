package models

import (
	"fmt"
	"regexp"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

var SendQQ func(int64, interface{})
var SendQQGroup func(int64, interface{})
var ListenQQPrivateMessage = func(uid int64, msg string) {
	SendQQ(uid, handleMessage(msg, "qq", uid))
}

var ListenQQGroupMessage = func(gid int64, uid int64, msg string) {
	if gid == Config.QQGroupID {
		if Config.QbotPublicMode {
			SendQQGroup(gid, handleMessage(msg, "qqg", gid, uid))
		} else {
			SendQQ(uid, handleMessage(msg, "qq", uid))
		}
	}
}

var handleMessage = func(msgs ...interface{}) interface{} {
	switch msgs[0].(string) {
	case "status", "状态":
		return Count()
	case "qrcode", "扫码", "二维码":
		url := fmt.Sprintf("http://127.0.0.1:%d/api/login/qrcode.png?%vid=%v", web.BConfig.Listen.HTTPPort, msgs[1], msgs[2])
		rsp, err := httplib.Get(url).Response()
		if err != nil {
			return nil
		}
		return rsp
	default:
		ss := regexp.MustCompile(`pt_key=([^;=\s]+);pt_pin=([^;=\s]+)`).FindAllStringSubmatch(msgs[0].(string), -1)
		if len(ss) > 0 {
			for _, s := range ss {
				ck := JdCookie{
					PtKey: s[1],
					PtPin: s[2],
				}
				if CookieOK(&ck) {
					if nck := GetJdCookie(ck.PtPin); nck != nil {
						ck.ToPool(ck.PtKey)
						msg := fmt.Sprintf("更新账号，%s", ck.PtPin)
						(&JdCookie{}).Push(msg)
						logs.Info(msg)
					} else {
						NewJdCookie(ck)
						msg := fmt.Sprintf("添加账号，%s", ck.PtPin)
						(&JdCookie{}).Push(msg)
						logs.Info(msg)
					}
				}
			}
			go func() {
				Save <- &JdCookie{}
			}()
		}

	}
	return nil
}
