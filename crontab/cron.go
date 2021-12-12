package crontab

import (
	"bcjh-bot/scheduler"
	"github.com/robfig/cron/v3"
)

var crontab *cron.Cron

func Register(server *scheduler.Scheduler) {
	crontab = cron.New(cron.WithSeconds())

	// 每周日晚上八点发厨神提醒
	crontab.AddFunc("0 0 20 * * 0", cookRaceRemind(server))

	crontab.Start()
}
