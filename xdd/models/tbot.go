package models

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	tb "gopkg.in/tucnak/telebot.v2"
)

var b *tb.Bot

func initTgBot() {
	go func() {
		if Config.TelegramBotToken == "" {
			return
		}
		var err error
		b, err = tb.NewBot(tb.Settings{
			Token:  Config.TelegramBotToken,
			Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		})
		if err != nil {
			logs.Warn("监听tgbot失败")
			return
		}
		b.Handle(tb.OnText, func(m *tb.Message) {
			if Config.TelegramUserID == 0 {
				Config.TelegramUserID = m.Sender.ID
				msg := fmt.Sprintf("tgbot自动绑定用户%d", m.Sender.ID)
				logs.Warn(msg)
				b.Send(m.Sender, msg)
			}
			if m.Text == "status" || m.Text == "状态" {
				tgBotNotify(Count())
			} else if m.Text == "qrcode" || m.Text == "扫码" {
				url := fmt.Sprintf("http://127.0.0.1:%d/api/login/qrcode.png?tgid=%d", web.BConfig.Listen.HTTPPort, m.Sender.ID)
				rsp, _ := httplib.Get(url).Response()
				b.SendAlbum(m.Sender, tb.Album{&tb.Photo{File: tb.FromReader(rsp.Body)}})
			}
		})
		logs.Info("监听tgbot")
		b.Start()
	}()
}

func tgBotNotify(msg string) {
	if b == nil {
		return
	}
	if Config.TelegramUserID == 0 {
		logs.Warn("tgbot未绑定用id")
		return
	}
	b.Send(&tb.User{ID: Config.TelegramUserID}, msg)
}

func SendTgMsg(id int, msg string) {
	b.Send(&tb.User{ID: id}, msg)
}
