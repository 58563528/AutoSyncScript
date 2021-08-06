package models

import "github.com/robfig/cron/v3"

var c *cron.Cron

func initCron() {
	c = cron.New()
	c.Start()
}
