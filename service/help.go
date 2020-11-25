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
	if len(args) >= 1 {
		switch args[0] {
		case "菜谱":
			msg = recipeHelp()
		default:
			msg = introHelp()
		}
	} else {
		msg = introHelp()
	}

	err := bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}

func introHelp() string {
	var msg string
	msg += fmt.Sprintf("爆炒江湖信息查询机器人\n")
	msg += fmt.Sprintf("使用 %s指令名 参数 查询信息\n", util.PrefixCharacter)
	msg += fmt.Sprintf("示例 %s厨师 羽十六\n", util.PrefixCharacter)

	msg += fmt.Sprintf("目前提供以下查询指令:\n")
	msg += fmt.Sprintf("帮助  图鉴网\n")
	msg += fmt.Sprintf("厨师  厨具  菜谱  贵客\n")
	msg += fmt.Sprintf("\n")

	msg += fmt.Sprintf("使用 %s帮助 指令名 查询用法\n", util.PrefixCharacter)
	msg += fmt.Sprintf("示例 %s帮助 厨师\n\n", util.PrefixCharacter)

	msg += fmt.Sprintf("数据来源: L图鉴网\n")
	msg += fmt.Sprintf("https://foodgame.gitee.io)\n")
	msg += fmt.Sprintf("项目地址(来提issue呀):\n")
	msg += fmt.Sprintf("https://github.com/Billdex/bcjh-bot")
	return msg
}

func recipeHelp() string {
	var msg string
	msg += fmt.Sprintf("菜谱信息查询功能\n")
	msg += fmt.Sprintf("基础信息查询:\n")
	msg += fmt.Sprintf("%s菜谱 菜谱名\n", util.PrefixCharacter)
	msg += fmt.Sprintf("示例 %s菜谱 荷包蛋\n", util.PrefixCharacter)
	msg += fmt.Sprintf("复合信息查询:\n")
	msg += fmt.Sprintf("%s菜谱 筛选方式-参数-排序方式\n", util.PrefixCharacter)
	msg += fmt.Sprintf("示例 %s菜谱 食材-牛排-单价\n", util.PrefixCharacter)
	msg += fmt.Sprintf("目前提供以下筛选方式:\n")
	msg += fmt.Sprintf("食材\n")
	msg += fmt.Sprintf("目前提供以下排序方式\n")
	msg += fmt.Sprintf("单时间 总时间 单价 金币效率\n")
	msg += fmt.Sprintf("耗材效率 食材效率")
	return msg
}
