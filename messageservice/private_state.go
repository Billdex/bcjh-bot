package messageservice

import (
	"bcjh-bot/global"
	"bcjh-bot/scheduler"
	"bcjh-bot/util/logger"
)

// AllowPrivate 启用机器人私聊查询
func AllowPrivate(c *scheduler.Context) {
	if err := global.SetBotState(c.GetBotId(), 0, true); err != nil {
		logger.Error("启用机器人私聊出错:", err)
		return
	}
	_, _ = c.Reply("已允许使用私聊")
}

// DisablePrivate 禁用机器人私聊查询
func DisablePrivate(c *scheduler.Context) {
	if err := global.SetBotState(c.GetBotId(), 0, false); err != nil {
		logger.Error("禁用机器人私聊出错:", err)
		_, _ = c.Reply("禁用失败")
		return
	}
	_, _ = c.Reply("已禁用私聊查询")
}
