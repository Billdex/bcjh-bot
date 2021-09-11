package middleware

import (
	"bcjh-bot/global"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
)

func MustAdmin(c *scheduler.Context) {
	if c.GetMessageType() != onebot.MessageTypeGroup || c.GetGroupEvent() == nil {
		c.Abort()
		return
	}
	event := c.GetGroupEvent()
	senderRole := event.Sender.Role
	if senderRole == onebot.GroupSenderRoleOwner || senderRole == onebot.GroupSenderRoleAdmin || global.IsSuperAdmin(c.GetSenderId()) {
		c.Next()
	} else {
		c.Abort()
		return
	}
}

func MustSuperAdmin(c *scheduler.Context) {
	if global.IsSuperAdmin(c.GetSenderId()) {
		c.Next()
	} else {
		c.Abort()
	}
}
