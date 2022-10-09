package middleware

import (
	"bcjh-bot/dao"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util/logger"
)

func CheckBlackList(c *scheduler.Context) {
	if c.GetMessageType() == onebot.MessageTypePrivate || c.GetPrivateEvent() != nil {
		c.Next()
		return
	}
	allow, err := dao.GetUserAllowState(c.GetSenderId(), c.GetGroupId())
	if err != nil {
		logger.Errorf("获取用户 %d group %d allow state 出错 %v", c.GetSenderId(), c.GetGroupId(), err)
		allow = true
	}
	if allow {
		c.Next()
		return
	} else {
		c.Abort()
		c.SetWarnMessage("该用户当前处于禁用状态")
		return
	}
}
