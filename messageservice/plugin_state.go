package messageservice

import (
	"bcjh-bot/dao"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util/logger"
	"fmt"
)

func EnablePluginInGroup(c *scheduler.Context) {
	if c.GetMessageType() != onebot.MessageTypeGroup || c.GetGroupEvent() == nil {
		return
	}
	if pluginName, ok := dao.GetPluginName(c.PretreatedMessage); ok {
		event := c.GetGroupEvent()
		if err := dao.SetPluginState(event.GroupId, pluginName, true); err != nil {
			logger.Errorf("设置群内功能%s启动出错:%v", pluginName, err)
			return
		}
		_, _ = c.Reply(fmt.Sprintf("%s功能已启用", pluginName))
	} else {
		_, _ = c.Reply(fmt.Sprintf("%s功能不存在或无法设置启用状态", c.PretreatedMessage))
	}
}

func DisablePluginInGroup(c *scheduler.Context) {
	if c.GetMessageType() != onebot.MessageTypeGroup || c.GetGroupEvent() == nil {
		return
	}
	if pluginName, ok := dao.GetPluginName(c.PretreatedMessage); ok {
		event := c.GetGroupEvent()
		if err := dao.SetPluginState(event.GroupId, pluginName, false); err != nil {
			logger.Errorf("设置群内功能%s关闭出错:%v", pluginName, err)
			_, _ = c.Reply(fmt.Sprintf("%s功能关闭出错", pluginName))
			return
		}
		_, _ = c.Reply(fmt.Sprintf("%s功能已关闭", pluginName))
	} else {
		_, _ = c.Reply(fmt.Sprintf("%s功能不存在或禁止关闭", c.PretreatedMessage))
	}
}
