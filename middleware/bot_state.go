package middleware

import (
	"bcjh-bot/global"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util"
)

func CheckBotState(c *scheduler.Context) {
	if c.GetMessageType() == onebot.MessageTypePrivate || c.GetPrivateEvent() != nil {
		c.Next()
		return
	}
	if util.InKeywordList(c.GetKeyword(), "开机", "关机") {
		c.Next()
	} else {
		event := c.GetGroupEvent()
		if botOn, err := global.GetBotState(event.SelfId, event.GroupId); err != nil {
			c.Abort()
			return
		} else {
			if botOn {
				c.Next()
			} else {
				c.Abort()
			}
		}
	}
}
