package noticeservice

import (
	"bcjh-bot/global"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util/logger"
	"fmt"
)

func ReplyRecallMessage(bot *onebot.Bot, event *onebot.NoticeEventGroupRecall) {
	if botOn, err := global.GetBotState(bot.BotId, event.GroupId); err != nil || !botOn {
		return
	}
	if pluginOn, err := global.GetPluginState(event.GroupId, "防撤回", false); err != nil || !pluginOn {
		return
	}

	recallMsg, err := bot.GetMsgInfo(int32(event.MessageId))
	if err != nil {
		logger.Error("获取群消息失败", err)
		return
	}
	msg := fmt.Sprintf("群成员撤回了一条消息\n")
	msg += fmt.Sprintf("消息内容:\n%s", recallMsg.RawMessage)

	_, _ = bot.SendGroupMessage(event.GroupId, msg)
}
