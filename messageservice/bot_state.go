package messageservice

import (
	"bcjh-bot/dao"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util/logger"
)

// EnableBotInGroup 在群聊中启用机器人
func EnableBotInGroup(c *scheduler.Context) {
	if c.GetMessageType() != onebot.MessageTypeGroup || c.GetGroupEvent() == nil {
		return
	}
	if err := dao.SetBotState(c.GetBotId(), c.GetGroupId(), true); err != nil {
		logger.Error("设置群内机器人启动出错:", err)
		return
	}
	_, _ = c.Reply("已开机")
}

// DisableBotInGroup 在群聊中停用机器人
func DisableBotInGroup(c *scheduler.Context) {
	if c.GetMessageType() != onebot.MessageTypeGroup || c.GetGroupEvent() == nil {
		return
	}
	if err := dao.SetBotState(c.GetBotId(), c.GetGroupId(), false); err != nil {
		logger.Error("设置群内机器人关闭出错:", err)
		_, _ = c.Reply("关机失败")
		return
	}
	_, _ = c.Reply("已关机")
}
