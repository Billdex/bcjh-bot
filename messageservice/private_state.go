package messageservice

import (
	"bcjh-bot/global"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util/logger"
)

func AllowPrivate(c *scheduler.Context) {
	if c.GetMessageType() == onebot.MessageTypePrivate {
		if err := global.SetBotState(c.GetBot().BotId, 0, true); err != nil {
			logger.Error("启用机器人私聊出错:", err)
			return
		}
		_, _ = c.Reply("已允许使用私聊")
	} else if c.GetMessageType() == onebot.MessageTypeGroup {
		atList := c.GetAtList()
		for _, at := range atList {
			if at == c.GetBot().BotId {
				if err := global.SetBotState(c.GetBot().BotId, 0, true); err != nil {
					logger.Error("启用机器人私聊出错:", err)
					return
				}
				_, _ = c.Reply("已允许使用私聊")
				break
			}
		}
	}
}

func DisablePrivate(c *scheduler.Context) {
	if c.GetMessageType() == onebot.MessageTypePrivate {
		if err := global.SetBotState(c.GetBot().BotId, 0, false); err != nil {
			logger.Error("禁用机器人私聊出错:", err)
			_, _ = c.Reply("禁用失败")
			return
		}
		_, _ = c.Reply("已禁用私聊查询")
	} else if c.GetMessageType() == onebot.MessageTypeGroup {
		atList := c.GetAtList()
		for _, at := range atList {
			if at == c.GetBot().BotId {
				if err := global.SetBotState(c.GetBot().BotId, 0, false); err != nil {
					logger.Error("禁用机器人私聊出错:", err)
					_, _ = c.Reply("禁用失败")
					return
				}
				_, _ = c.Reply("已禁用私聊查询")
				break
			}
		}
	}
}
