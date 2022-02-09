package middleware

import (
	"bcjh-bot/global"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util"
)

func CheckBotState(c *scheduler.Context) {
	if c.GetMessageType() == onebot.MessageTypePrivate || c.GetPrivateEvent() != nil {
		if util.InKeywordList(c.GetKeyword(), "允许私聊", "开启私聊", "禁止私聊", "禁用私聊", "关闭私聊") {
			c.Next()
			return
		}
		if global.IsSuperAdmin(c.GetSenderId()) {
			c.Next()
			return
		}
		if privateOn, _ := global.GetBotState(c.GetBot().BotId, 0); privateOn {
			c.Next()
			return
		} else {
			c.Abort()
			c.SetWarnMessage("该 bot 未开启私聊查询")
			return
		}
	}
	if c.GetMessageType() == onebot.MessageTypeGroup || c.GetGroupEvent() != nil {
		if util.InKeywordList(c.GetKeyword(), "开机", "关机") {
			c.Next()
			return
		}
		event := c.GetGroupEvent()
		if botOn, _ := global.GetBotState(c.GetBot().BotId, event.GroupId); botOn {
			c.Next()
			return
		} else {
			c.Abort()
			c.SetWarnMessage("该群未启用此bot")
			return
		}
	}
}
