package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/cdle/jd_study/xdd/models"

	"github.com/beego/beego/v2/client/httplib"
	qrcode "github.com/skip2/go-qrcode"
)

//LoginController 主页控制器
type LoginController struct {
	BaseController
}

type StepOne struct {
	SToken string `json:"s_token"`
}

type StepTwo struct {
	Token string `json:"token"`
}

type StepThree struct {
	CheckIP int    `json:"check_ip"`
	Errcode int    `json:"errcode"`
	Message string `json:"message"`
}

var JdCookieRunners sync.Map
var jdua = models.GetUserAgent

func (c *LoginController) GetQrcode() {
	if v := c.GetSession("jd_token"); v != nil {
		token := v.(string)
		if v, ok := JdCookieRunners.Load(token); ok {
			if len(v.([]string)) == 2 {
				var url = `https://plogin.m.jd.com/cgi-bin/m/tmauth?appid=300&client_type=m&token=` + token
				data, _ := qrcode.Encode(url, qrcode.Medium, 256)
				c.Ctx.WriteString(`{"url":"` + url + `","img":"` + base64.StdEncoding.EncodeToString(data) + `"}`)
				return
			}
		}
	}
	var state = time.Now().Unix()
	var url = fmt.Sprintf(`https://plogin.m.jd.com/cgi-bin/mm/new_login_entrance?lang=chs&appid=300&returnurl=https://wq.jd.com/passport/LoginRedirect?state=%d&returnurl=https://home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&/myJd/home.action&source=wq_passport`,
		state)
	req := httplib.Get(url)
	req.Header("Connection", "Keep-Alive")
	req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.Header("Accept", "application/json, text/plain, */*")
	req.Header("Accept-Language", "zh-cn")
	req.Header("Referer", url)
	req.Header("User-Agent", jdua())
	req.Header("Host", "plogin.m.jd.com")
	rsp, err := req.Response()
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
	data, err := ioutil.ReadAll(rsp.Body)
	so := StepOne{}
	err = json.Unmarshal(data, &so)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
	cookies := strings.Join(rsp.Header.Values("Set-Cookie"), " ")
	var cookie = strings.Join([]string{
		"guid=" + FetchJdCookieValue("guid", cookies),
		"lang=chs",
		"lsid=" + FetchJdCookieValue("lsid", cookies),
		"lstoken=" + FetchJdCookieValue("lstoken", cookies),
	}, ";")
	state = time.Now().Unix()
	req = httplib.Post(
		fmt.Sprintf(`https://plogin.m.jd.com/cgi-bin/m/tmauthreflogurl?s_token=%s&v=%d&remember=true`,
			so.SToken,
			state),
	)
	req.Header("Connection", "Keep-Alive")
	req.Header("Content-Type", "application/x-www-form-urlencoded; Charset=UTF-8")
	req.Header("Accept", "application/json, text/plain, */*")
	req.Header("Cookie", cookie)
	req.Header("Referer", fmt.Sprintf(`https://plogin.m.jd.com/login/login?appid=300&returnurl=https://wqlogin2.jd.com/passport/LoginRedirect?state=%d&returnurl=//home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&/myJd/home.action&source=wq_passport`,
		state),
	)
	req.Header("User-Agent", jdua())
	req.Header("Host", "plogin.m.jd.com")
	req.Body(fmt.Sprintf(`{
		'lang': 'chs',
		'appid': 300,
		'returnurl': 'https://wqlogin2.jd.com/passport/LoginRedirect?state=%dreturnurl=//home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&/myJd/home.action&source=wq_passport',
	 }`, state))
	rsp, err = req.Response()
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
	data, err = ioutil.ReadAll(rsp.Body)
	st := StepTwo{}
	err = json.Unmarshal(data, &st)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
	url = `https://plogin.m.jd.com/cgi-bin/m/tmauth?client_type=m&appid=300&token=` + st.Token
	cookies = strings.Join(rsp.Header.Values("Set-Cookie"), " ")

	c.SetSession("jd_token", st.Token)
	c.SetSession("jd_cookie", cookie)
	okl_token := FetchJdCookieValue("okl_token", cookies)
	c.SetSession("jd_okl_token", okl_token)
	data, _ = qrcode.Encode(url, qrcode.Medium, 256)
	// fmt.Println(st.Token, cookie, okl_token)
	JdCookieRunners.Store(st.Token, []string{cookie, okl_token})
	c.Ctx.WriteString(`{"url":"` + url + `","img":"` + base64.StdEncoding.EncodeToString(data) + `"}`) //"data:image/png;base64," +
}

func init() {
	go func() {
		for {
			time.Sleep(time.Second)
			JdCookieRunners.Range(func(k, v interface{}) bool {
				jd_token := k.(string)
				vv := v.([]string)
				if len(vv) == 2 {
					cookie := vv[0]
					okl_token := vv[1]
					// fmt.Println(jd_token, cookie, okl_token)
					result := CheckLogin(jd_token, cookie, okl_token)
					// fmt.Println(result)
					switch result {
					case "成功":

					case "授权登录未确认":
					default: //失效

					}
				}
				return true
			})
		}
	}()
}

//Query 查询
func (c *LoginController) Query() {
	if v := c.GetSession("jd_token"); v == nil {
		c.Ctx.WriteString("重新获取二维码")
		return
	} else {
		token := v.(string)
		if v, ok := JdCookieRunners.Load(token); !ok {
			c.Ctx.WriteString("重新获取二维码")
			return
		} else {
			if len(v.([]string)) == 2 {
				c.Ctx.WriteString("授权登录未确认")
				return
			} else {
				pin := v.([]string)[0]
				c.SetSession("pin", pin)
				if note := c.GetString("note"); note != "" {
					if ck := models.GetJdCookie(pin); ck != nil {
						ck.Updates(models.Note, note)
					}
				}
				if strings.Contains(models.Config.Master, pin) {
					c.Ctx.WriteString("登录")
				} else {
					c.Ctx.WriteString("成功")
				}
				return
			}
		}
	}
}

func CheckLogin(token, cookie, okl_token string) string {
	state := time.Now().Unix()
	req := httplib.Post(
		fmt.Sprintf(`https://plogin.m.jd.com/cgi-bin/m/tmauthchecktoken?&token=%s&ou_state=0&okl_token=%s`,
			token,
			okl_token,
		),
	)
	req.Header("Referer", fmt.Sprintf(`https://plogin.m.jd.com/login/login?appid=300&returnurl=https://wqlogin2.jd.com/passport/LoginRedirect?state=%d&returnurl=//home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&/myJd/home.action&source=wq_passport`,
		state),
	)
	req.Header("Cookie", cookie)
	req.Header("Connection", "Keep-Alive")
	req.Header("Content-Type", "application/x-www-form-urlencoded; Charset=UTF-8")
	req.Header("Accept", "application/json, text/plain, */*")
	req.Header("User-Agent", jdua())
	req.Header("Host", "plogin.m.jd.com")

	req.Param("lang", "chs")
	req.Param("appid", "300")
	req.Param("returnurl", fmt.Sprintf("https://wqlogin2.jd.com/passport/LoginRedirect?state=%d&returnurl=//home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&/myJd/home.action", state))
	req.Param("source", "wq_passport")

	rsp, err := req.Response()
	if err != nil {
		return err.Error()
	}
	data, err := ioutil.ReadAll(rsp.Body)
	sth := StepThree{}
	err = json.Unmarshal(data, &sth)
	if err != nil {
		return err.Error()
	}
	// fmt.Println(sth)
	switch sth.Errcode {
	case 0:
		cookies := strings.Join(rsp.Header.Values("Set-Cookie"), " ")
		pt_key := FetchJdCookieValue("pt_key", cookies)
		pt_pin := FetchJdCookieValue("pt_pin", cookies)
		if pt_pin == "" {
			JdCookieRunners.Delete(token)
			return sth.Message
		}
		go func() {
			ck := models.JdCookie{
				PtKey: pt_key,
				PtPin: pt_pin,
			}
			if nck := models.GetJdCookie(ck.PtPin); nck != nil {
				ck.ToPool(ck.PtKey)
				msg := fmt.Sprintf("更新账号，%s", ck.PtPin)
				models.QywxNotify(&models.QywxConfig{
					Content: msg,
				})
				logs.Info(msg)
			} else {
				models.NewJdCookie(ck)
				msg := &models.QywxConfig{
					Content: fmt.Sprintf("添加账号，%s", ck.PtPin),
				}
				models.QywxNotify(msg)
				logs.Info(msg)
			}
			go func() {
				models.Save <- &ck
			}()
		}()
		JdCookieRunners.Store(token, []string{pt_pin})
		return "成功"
	case 19: //Token无效，请退出重试
		JdCookieRunners.Delete(token)
		return sth.Message
	case 21: //Token不存在，请退出重试
		JdCookieRunners.Delete(token)
		return sth.Message
	case 176: //授权登录未确认
		return sth.Message
	default:
		JdCookieRunners.Delete(token)
	}
	return ""
}

func FetchJdCookieValue(key string, cookies string) string {
	match := regexp.MustCompile(key + `=([^;]*);{0,1}\s`).FindStringSubmatch(cookies)
	if len(match) == 2 {
		return match[1]
	} else {
		return ""
	}
}
