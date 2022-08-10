package crontab

import (
	"bcjh-bot/global"
	"bcjh-bot/scheduler"
	"bcjh-bot/util/logger"
	"time"
)

func cookRaceRemind(s *scheduler.Scheduler) func() {
	return func() {
		bots := s.Engine.GetBots()
		msg := "本周厨神大赛就快要结束啦，请还未参赛的小伙伴尽快参赛哦！对啦，周常也别忘记做咯~"
		for _, bot := range bots {
			groups, err := bot.GetGroupList()
			if err != nil {
				logger.Errorf("未获取[bot %d]到group列表, err:%v", bot.BotId, err)
				continue
			}
			for _, group := range groups {
				if botOk, _ := global.GetBotState(bot.BotId, group.GroupId); botOk {
					if pluginOk, _ := global.GetPluginState(group.GroupId, "厨神提醒", true); pluginOk {
						_, _ = bot.SendGroupMessage(group.GroupId, msg)
						// 等待一点时间，免得发得太快被风控
						time.Sleep(200 * time.Millisecond)
					}
				}
			}
		}
	}
}
