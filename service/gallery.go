package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util/logger"
	"fmt"
)

func GalleryWebsite(c *onebot.Context, args []string) {
	logger.Info("查询图鉴网, 参数:", args)

	var msg string
	foodgame := "https://foodgame.gitee.io/"
	bcjh := "https://bcjh.gitee.io/"

	msg += fmt.Sprintf("L图鉴网: %s\n", foodgame)
	msg += fmt.Sprintf("白菜菊花手机图鉴网: %s", bcjh)

	err := bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}
