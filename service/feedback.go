package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/model/database"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"strings"
)

const (
	close = iota
	open
	finish
)

func Feedback(c *onebot.Context, args []string) {
	logger.Info("有人提交了反馈:", args)
	if len(args) == 0 {
		return
	}
	message := strings.Join(args, util.ArgsSplitCharacter)
	feedback := new(database.Feedback)
	feedback.Sender = c.Sender.UserId
	feedback.Nickname = c.Sender.Nickname
	feedback.Message = message
	feedback.Status = open
	affected, err := database.DB.Insert(feedback)
	if err != nil {
		logger.Error("插入数据失败!", err)
		_ = bot.SendMessage(c, "不小心反馈失败了呢!")
		return
	}
	if affected == 0 {
		_ = bot.SendMessage(c, "不小心反馈失败了呢!")
		return
	}
	err = bot.SendMessage(c, "收到! 反馈信息已经记在小本本上啦~")
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}
