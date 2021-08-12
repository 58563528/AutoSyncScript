package models

import "github.com/beego/beego/v2/client/httplib"

func pushPlus(token string, content string) {
	httplib.Get("http://pushplus.hxtrip.com/send?token=" + token + "&title=小滴滴&content=" + content + "&template=html").Response()
}
