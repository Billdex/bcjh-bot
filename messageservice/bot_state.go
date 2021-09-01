package messageservice

import (
	"bcjh-bot/global"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util/logger"
)

func EnableBotInGroup(c *scheduler.Context) {
	if c.GetMessageType() != onebot.MessageTypeGroup || c.GetGroupEvent() == nil {
		return
	}
	atList := c.GetAtList()
	for _, at := range atList {
		if at == c.GetBot().BotId {
			event := c.GetGroupEvent()
			if err := global.SetBotState(event.SelfId, event.GroupId, true); err != nil {
				logger.Error("设置群内机器人启动出错:", err)
				return
			}
			_, _ = c.Reply("已开机")
			break
		}
	}
}

func DisableBotInGroup(c *scheduler.Context) {
	if c.GetMessageType() != onebot.MessageTypeGroup || c.GetGroupEvent() == nil {
		return
	}
	atList := c.GetAtList()
	for _, at := range atList {
		if at == c.GetBot().BotId {
			event := c.GetGroupEvent()
			if err := global.SetBotState(event.SelfId, event.GroupId, false); err != nil {
				logger.Error("设置群内机器人关闭出错:", err)
				_, _ = c.Reply("关机出错")
				return
			}
			_, _ = c.Reply("已关机")
			break
		}
	}
}
