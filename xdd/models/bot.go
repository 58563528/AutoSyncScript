package models

import (
	"fmt"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/server/web"
)

var NotifyQQ func(int64, interface{})
var ListenQQ = func(uid int64, msg string) {
	NotifyQQ(uid, handleMessage(msg, "qq", uid))
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
	}
	return nil
}
