package models

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
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
			if m.Text == "status" {
				tgBotNotify(Count())
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
