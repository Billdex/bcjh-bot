package middleware

import (
	"bcjh-bot/global"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
)

func CheckBlackList(c *scheduler.Context) {
	if c.GetMessageType() == onebot.MessageTypePrivate || c.GetPrivateEvent() != nil {
		c.Next()
		return
	}
	if allow := global.GetUserAllowState(c.GetSenderId(), c.GetGroupEvent().GroupId); allow {
		c.Next()
		return
	} else {
		c.Abort()
		return
	}
}
