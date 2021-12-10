package messageservice

import (
	"bcjh-bot/global"
	"bcjh-bot/scheduler"
	"fmt"
	"time"
)

func PublicNotice(c *scheduler.Context) {
	msg := fmt.Sprintf("来自「%s」的公告:\n%s", c.GetSenderNickname(), c.PretreatedMessage)
	bots := c.GetLinkBotList()
	for _, bot := range bots {
		groups, err := bot.GetGroupList()
		if err != nil {
			_, _ = c.Reply(fmt.Sprintf("未获取[bot %d]到group列表, err:%v", bot.BotId, err))
		}
		for _, group := range groups {
			if botOk, _ := global.GetBotState(bot.BotId, group.GroupId); botOk {
				if pluginOk, _ := global.GetPluginState(group.GroupId, "公告", true); pluginOk {
					_, _ = bot.SendGroupMessage(group.GroupId, msg)
					// 等待一点时间，免得发得太快被风控
					time.Sleep(200 * time.Millisecond)
				}
			}
		}
	}
}
