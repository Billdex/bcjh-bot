package middleware

import (
	"bcjh-bot/global"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
)

func CheckPluginState(defaultState bool) scheduler.HandleFunc {
	return func(c *scheduler.Context) {
		if c.GetMessageType() == onebot.MessageTypePrivate || c.GetPrivateEvent() != nil {
			c.Next()
			return
		}
		event := c.GetGroupEvent()
		if pluginName, ok := global.GetPluginName(c.GetKeyword()); ok {
			if pluginOn, err := global.GetPluginState(event.GroupId, pluginName, defaultState); err != nil {
				c.Abort()
				return
			} else {
				if pluginOn {
					c.Next()
				} else {
					c.Abort()
				}
			}
		} else {
			c.Abort()
			return
		}
	}
}
