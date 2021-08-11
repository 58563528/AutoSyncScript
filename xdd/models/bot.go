package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

var SendQQ func(int64, interface{})
var SendQQGroup func(int64, interface{})
var ListenQQPrivateMessage = func(uid int64, msg string) {
	SendQQ(uid, handleMessage(msg, "qq", int(uid)))
}

var ListenQQGroupMessage = func(gid int64, uid int64, msg string) {
	if gid == Config.QQGroupID {
		if Config.QbotPublicMode {
			SendQQGroup(gid, handleMessage(msg, "qqg", int(gid), int(uid)))
		} else {
			SendQQ(uid, handleMessage(msg, "qq", int(uid)))
		}
	}
}

var replies = map[string]string{}

func InitReplies() {
	f, err := os.Open(ExecPath + "/conf/reply.php")
	if err == nil {
		defer f.Close()
		data, _ := ioutil.ReadAll(f)
		ss := regexp.MustCompile("`([^`]+)`\\s*=>\\s*`([^`]+)`").FindAllStringSubmatch(string(data), -1)
		for _, s := range ss {
			replies[s[1]] = s[2]
		}
	}
}

var handleMessage = func(msgs ...interface{}) interface{} {
	msg := msgs[0].(string)
	tp := msgs[1].(string)
	id := msgs[2].(int)
	switch msg {
	case "status", "状态":
		return Count()
	case "qrcode", "扫码", "二维码":
		url := fmt.Sprintf("http://127.0.0.1:%d/api/login/qrcode.png?%vid=%v", web.BConfig.Listen.HTTPPort, tp, id)
		rsp, err := httplib.Get(url).Response()
		if err != nil {
			return nil
		}
		return rsp
	case "查询", "query":
		if tp == "qq" {
			cks := GetJdCookies()
			for _, ck := range cks {
				if ck.QQ == id {
					SendQQ(int64(id), ck.Query())
				}
			}
		}
		return nil
	default:

		{ //
			ss := regexp.MustCompile(`pt_key=([^;=\s]+);pt_pin=([^;=\s]+)`).FindAllStringSubmatch(msg, -1)
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
				return nil
			}
		}

		{
			s := regexp.MustCompile(`(查询|query)\s+(.*)`).FindStringSubmatch(msg)
			if len(s) > 0 {
				cks := GetJdCookies()
				a := s[2]
				{
					if s := strings.Split(a, "-"); len(s) == 2 {
						for i, ck := range cks {
							if i+1 >= Int(s[0]) && i+1 <= Int(s[1]) {
								switch tp {
								case "tg":
									tgBotNotify(ck.Query())
								case "qq":
									if id == ck.QQ {
										SendQQ(int64(id), ck.Query())
									} else {
										SendQQ(Config.QQID, ck.Query())
									}
								case "qqg":
									uid := msgs[3].(int)
									if uid == ck.QQ || uid == int(Config.QQID) {
										SendQQGroup(int64(id), ck.Query())
									}
								}
							}
						}
						return nil
					}

				}
				{
					// if x := regexp.MustCompile(`^\d+$`).FindString(a); x != "" {
					// 	id := Int(x)
					// 	for i, ck := range cks {
					// 		if i+1 == id {
					// 			switch tp {
					// 			case "tg":
					// 				tgBotNotify(ck.Query())
					// 			case "qq":
					// 				if id == ck.QQ {
					// 					SendQQ(int64(id), ck.Query())
					// 				} else {
					// 					SendQQ(Config.QQID, ck.Query())
					// 				}
					// 			case "qqg":
					// 				uid := msgs[3].(int)
					// 				if uid == ck.QQ || uid == int(Config.QQID) {
					// 					SendQQGroup(int64(id), ck.Query())
					// 				}
					// 			}
					// 		}
					// 	}
					// 	return nil
					// }
					if x := regexp.MustCompile(`^[\s\d,]+$`).FindString(a); x != "" {
						xx := regexp.MustCompile(`(\d+)`).FindAllStringSubmatch(a, -1)
						for i, ck := range cks {
							for _, x := range xx {
								if fmt.Sprint(i+1) == x[1] {
									switch tp {
									case "tg":
										tgBotNotify(ck.Query())
									case "qq":
										if id == ck.QQ {
											SendQQ(int64(id), ck.Query())
										} else {
											SendQQ(Config.QQID, ck.Query())
										}
									case "qqg":
										uid := msgs[3].(int)
										if uid == ck.QQ || uid == int(Config.QQID) {
											SendQQGroup(int64(id), ck.Query())
										}
									}
								}
							}

						}
						return nil
					}
				}
				{
					a = strings.Replace(a, " ", "", -1)
					for _, ck := range cks {
						if strings.Contains(ck.Note, a) || strings.Contains(ck.Nickname, a) || strings.Contains(ck.PtPin, a) {
							switch tp {
							case "tg":
								tgBotNotify(ck.Query())
							case "qq":
								if id == ck.QQ {
									SendQQ(int64(id), ck.Query())
								} else {
									SendQQ(Config.QQID, ck.Query())
								}
							case "qqg":
								uid := msgs[3].(int)
								if uid == ck.QQ || uid == int(Config.QQID) {
									SendQQGroup(int64(id), ck.Query())
								}
							}
						}
					}
					return nil
				}

			}
		}
		for k, v := range replies {
			if regexp.MustCompile(k).FindString(msg) != "" {
				return v
			}
		}

	}
	return nil
}
