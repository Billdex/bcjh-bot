package middleware

import (
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
)

// MustAtSelf 判定该消息是否是指向自己的
func MustAtSelf(c *scheduler.Context) {
	// 如果消息类型为私聊，则始终视为指向自己
	if c.GetMessageType() == onebot.MessageTypePrivate {
		c.Next()
		return
	}
	// 如果消息类型为群聊，则需要判断 at 列表里是否有自己
	if c.GetMessageType() == onebot.MessageTypeGroup {
		atList := c.GetAtList()
		for _, at := range atList {
			if at == c.GetBot().BotId {
				c.Next()
				return
			}
		}
	}
	// 获取不到消息类型则直接忽略该消息
	c.Abort()
	return
}
