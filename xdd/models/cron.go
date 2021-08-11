package models

import "github.com/robfig/cron/v3"

var c *cron.Cron

func initCron() {
	c = cron.New()
	if Config.DailyAssetPushCron != "" {
		c.AddFunc(Config.DailyAssetPushCron, DailyAssetsPush)
	}
	c.Start()
}
