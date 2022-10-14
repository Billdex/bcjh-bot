package middleware

import (
	"bcjh-bot/dao"
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
	if senderRole == onebot.GroupSenderRoleOwner || senderRole == onebot.GroupSenderRoleAdmin || dao.IsSuperAdmin(c.GetSenderId()) {
		c.Next()
	} else {
		if ok, _ := dao.GetBotState(c.GetBot().BotId, event.GroupId); ok {
			_, _ = c.Reply(e.PermissionDeniedNote)
		}
		c.SetWarnMessage("用户没有管理员权限")
		c.Abort()
		return
	}
}

func MustSuperAdmin(c *scheduler.Context) {
	if dao.IsSuperAdmin(c.GetSenderId()) {
		c.Next()
	} else {
		if c.GetMessageType() == onebot.MessageTypeGroup {
			if ok, _ := dao.GetBotState(c.GetBot().BotId, c.GetGroupEvent().GroupId); ok {
				_, _ = c.Reply(e.PermissionDeniedNote)
			}
		}
		c.SetWarnMessage("用户没有超管权限")
		c.Abort()
	}
}
