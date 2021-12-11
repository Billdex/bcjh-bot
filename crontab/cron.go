package crontab

import (
	"bcjh-bot/scheduler/onebot"
	"github.com/robfig/cron/v3"
)

func Register(server *onebot.Server) {
	crontab := cron.New(cron.WithSeconds())

	// 每周日晚上八点发厨神提醒
	crontab.AddFunc("0 0 20 * * 0", cookRaceRemind(server))

	crontab.Start()
}
