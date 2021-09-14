package middleware

import (
	"bcjh-bot/global"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util/e"
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
		if ok, _ := global.GetBotState(c.GetBot().BotId, event.GroupId); ok {
			_, _ = c.Reply(e.PermissionDeniedNote)
		}
		c.Abort()
		return
	}
}

func MustSuperAdmin(c *scheduler.Context) {
	if global.IsSuperAdmin(c.GetSenderId()) {
		c.Next()
	} else {
		if c.GetMessageType() == onebot.MessageTypeGroup {
			if ok, _ := global.GetBotState(c.GetBot().BotId, c.GetGroupEvent().GroupId); ok {
				_, _ = c.Reply(e.PermissionDeniedNote)
			}
		}
		c.Abort()
	}
}
