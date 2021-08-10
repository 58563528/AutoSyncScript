package models

func (ck *JdCookie) Push(msg string) {
	if Config.QywxKey != "" {
		go qywxNotify(&QywxConfig{Content: msg})
	}
	if Config.TelegramBotToken != "" {
		go tgBotNotify(msg)
	}
	if Config.QQID != 0 {
		go NotifyQQ(Config.QQID, msg)
	}
}
