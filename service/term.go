package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util/logger"
)

func TermInfo(c *onebot.Context, args []string) {
	logger.Info("术语信息查询")
	msg := termHelp()
	err := bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}
