package messageservice

import (
	"bcjh-bot/global"
	"bcjh-bot/scheduler"
	"fmt"
)

func PublicNotice(c *scheduler.Context) {
	msg := c.PretreatedMessage
	bots := c.GetLinkBotList()
	for _, bot := range bots {
		groups, err := c.GetBot().GetGroupList()
		if err != nil {
			_, _ = c.Reply(fmt.Sprintf("未获取[bot %d]到group列表, err:%v", bot.BotId, err))
		}
		for _, group := range groups {
			if ok, _ := global.GetBotState(bot.BotId, group.GroupId); ok {
				_, _ = bot.SendGroupMessage(group.GroupId, msg)
			}
		}
	}
}
