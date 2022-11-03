package middleware

import (
	"bcjh-bot/dao"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util/logger"
)

// CheckPluginState 使用该中间件需要在global.plugin_state中添加对应插件的名称与别名
func CheckPluginState(defaultState bool) scheduler.HandleFunc {
	return func(c *scheduler.Context) {
		if c.GetMessageType() == onebot.MessageTypePrivate || c.GetPrivateEvent() != nil {
			c.Next()
			return
		}
		event := c.GetGroupEvent()
		if pluginName, ok := dao.GetPluginName(c.GetKeyword()); ok {
			if pluginOn, err := dao.GetPluginState(event.GroupId, pluginName, defaultState); err != nil {
				logger.Errorf("获取插件状态失败 %v", err)
				c.Abort()
				return
			} else {
				if pluginOn {
					c.Next()
				} else {
					c.SetWarnMessage("插件功能未启用")
					c.Abort()
				}
			}
		} else {
			c.SetWarnMessage("插件功能未启用")
			c.Abort()
			return
		}
	}
}
