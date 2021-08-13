package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
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
	if _, ok := replies["壁纸"]; !ok {
		replies["壁纸"] = "https://acg.toubiec.cn/random.php"
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
				xyb := 0
				for _, s := range ss {
					ck := JdCookie{
						PtKey: s[1],
						PtPin: s[2],
					}
					if CookieOK(&ck) {
						xyb++
						if tp == "qq" {
							ck.QQ = id

						} else if tp == "tg" {
							ck.Telegram = id
						} else if tp == "qqg" {
							ck.QQ = msgs[3].(int)
						}
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
				return fmt.Sprintf("许愿币+%d", xyb)
			}
		}
		{
			s := regexp.MustCompile(`([^\s]+)\s+(.*)`).FindStringSubmatch(msg)
			if len(s) > 0 {
				v := s[2]
				switch s[1] {
				case "查询", "query":
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
				case "许愿":
					if tp == "qqg" {
						id = msgs[3].(int)
					}
					b := 0
					for _, ck := range GetJdCookies() {
						if id == ck.QQ || id == ck.Telegram {
							b++
						}
					}
					if b <= 0 {
						return "许愿币不足"
					} else {
						(&JdCookie{}).Push(fmt.Sprintf("%d许愿%s，许愿币余额%d。", id, v, b))
						return "收到许愿"
					}
				case "扣除许愿币":
					id, _ := strconv.Atoi(v)
					b := 0
					k := 0
					for _, ck := range GetJdCookies() {
						if id == ck.QQ || id == ck.Telegram {
							if k <= 5 {
								ck.Updates(map[string]interface{}{
									QQ:       0,
									Telegram: 0,
								})
								k++
							} else {
								b++
							}
						}
					}
					return fmt.Sprintf("操作成功，%d剩余许愿币%d", id, b)
				}

			}
		}
		for k, v := range replies {
			if regexp.MustCompile(k).FindString(msg) != "" {
				if regexp.MustCompile(`^https{0,1}://[^\x{4e00}-\x{9fa5}\n\r\s]{3,}$`).FindString(v) != "" {
					url := v
					rsp, err := httplib.Get(url).Response()
					if err != nil {
						return nil
					}
					return rsp
				}
				return v
			}
		}
	}
	return nil
}
