package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"fmt"
)

func HelpGuide(c *onebot.Context, args []string) {
	logger.Info("帮助查询, 参数:", args)

	var msg string
	msg += fmt.Sprintf("按照格式发送信息进行查询:\n")
	msg += fmt.Sprintf("%s指令 参数\n", util.PrefixCharacter)
	msg += fmt.Sprintf("示例 %s厨师 羽十六\n", util.PrefixCharacter)

	msg += fmt.Sprintf("目前提供以下查询指令:\n")
	msg += fmt.Sprintf("帮助  图鉴网\n")
	msg += fmt.Sprintf("厨师  厨具  菜谱\n")

	msg += fmt.Sprintf("数据来源: L图鉴网(https://foodgame.gitee.io)")

	err := bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}
