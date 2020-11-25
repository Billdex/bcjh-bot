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
		case "帮助":
			msg = introHelp()
		case "反馈":
			msg = feedbackHelp()
		case "图鉴网":
			msg = galleryWebsiteHelp()
		case "菜谱":
			msg = recipeHelp()
		default:
			msg = "似乎还没有开发这个功能呢~"
		}
	} else {
		msg = introHelp()
	}

	err := bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}

// 功能指引
func introHelp() string {
	var msg string
	msg += fmt.Sprintf("[爆炒江湖信息查询机器人]\n")
	msg += fmt.Sprintf("使用 %s功能名 参数 查询信息\n", util.PrefixCharacter)
	msg += fmt.Sprintf("示例 %s厨师 羽十六\n", util.PrefixCharacter)

	msg += fmt.Sprintf("目前提供以下功能:\n")
	msg += fmt.Sprintf("帮助, 反馈, 图鉴网, ")
	msg += fmt.Sprintf("厨师, 厨具, 菜谱, 贵客\n")
	msg += fmt.Sprintf("\n")

	msg += fmt.Sprintf("使用 %s帮助 功能名 查询用法\n", util.PrefixCharacter)
	msg += fmt.Sprintf("示例 %s帮助 厨师\n\n", util.PrefixCharacter)

	msg += fmt.Sprintf("数据来源: L图鉴网\n")
	msg += fmt.Sprintf("https://foodgame.gitee.io\n")
	msg += fmt.Sprintf("项目地址(给个star呀):\n")
	msg += fmt.Sprintf("https://github.com/Billdex/bcjh-bot")
	return msg
}

// 反馈功能指引
func feedbackHelp() string {
	var msg string
	msg += fmt.Sprintf("[问题反馈与建议]\n")
	msg += fmt.Sprintf("在使用过程中如果遇到了什么bug或者有什么好的建议，可以通过该功能反馈给我\n")
	msg += fmt.Sprintf("反馈方式:\n")
	msg += fmt.Sprintf("%s反馈 问题描述或建议\n", util.PrefixCharacter)
	msg += fmt.Sprintf("如果比较紧急可以私聊我:\n")
	msg += fmt.Sprintf("QQ:591404144")
	return msg
}

// 图鉴网功能指引
func galleryWebsiteHelp() string {
	var msg string
	msg += fmt.Sprintf("[图鉴网 网址查询]\n")
	msg += fmt.Sprintf("给出L图鉴网与手机版图鉴网地址，方便记不住网址的小可爱快速访问。")
	return msg
}

// 菜谱功能指引
func recipeHelp() string {
	var msg string
	msg += fmt.Sprintf("[菜谱信息查询]\n")
	msg += fmt.Sprintf("基础信息查询: %s菜谱 菜谱名\n", util.PrefixCharacter)
	msg += fmt.Sprintf("示例: %s菜谱 荷包蛋\n", util.PrefixCharacter)
	msg += fmt.Sprintf("复合信息查询:\n")
	msg += fmt.Sprintf("%s菜谱 筛选条件-参数-排序方式-单价下限或稀有度下限\n", util.PrefixCharacter)
	msg += fmt.Sprintf("示例: %s菜谱 食材-茄子-单时间-$100\n", util.PrefixCharacter)
	msg += fmt.Sprintf("%s菜谱 任意-耗材效率-4火\n", util.PrefixCharacter)
	msg += fmt.Sprintf("目前提供以下筛选条件:\n")
	msg += fmt.Sprintf("任意(不填参数), 食材, 技法\n")
	msg += fmt.Sprintf("目前提供以下排序方式\n")
	msg += fmt.Sprintf("单时间, 总时间, 单价, 金币效率, ")
	msg += fmt.Sprintf("耗材效率, 食材效率")
	return msg
}
