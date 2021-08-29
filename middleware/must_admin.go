package middleware

import (
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
	if senderRole == onebot.GroupSenderRoleOwner || senderRole == onebot.GroupSenderRoleAdmin {
		c.Next()
	} else {
		c.Abort()
		return
	}
}
