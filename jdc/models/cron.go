package models

import "github.com/robfig/cron/v3"

var c *cron.Cron

func initCron() {
	c = cron.New()
	c.Start()
	c.AddFunc("*/10 * * * *", func() {
		Save <- &JdCookie{}
	})
	c.AddFunc("1 * * * *", initCookie)
}
