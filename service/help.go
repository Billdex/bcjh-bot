package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/config"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"fmt"
	"strings"
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
		case "游戏术语", "术语", "黑话":
			msg = termHelp()
		case "厨师", "厨子":
			msg = chefHelp()
		case "厨具":
			msg = equipmentHelp()
		case "菜谱":
			msg = recipeHelp()
		case "贵客":
			msg = guestHelp()
		case "符文":
			msg = antiqueHelp()
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
	preChar := util.PrefixCharacters[0]
	sb := strings.Builder{}
	sb.WriteString("【爆炒机器人查询使用说明】\n")
	sb.WriteString(fmt.Sprintf("使用『%s功能名 参数』查询信息\n", preChar))
	sb.WriteString(fmt.Sprintf("示例「%s厨师 羽十六」\n", preChar))
	sb.WriteString("目前提供以下功能:\n")
	sb.WriteString(fmt.Sprintf("帮助, 反馈, 图鉴网, 术语, 厨师, 厨具, 菜谱, 贵客, 符文\n"))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("使用 %s帮助 功能名 查询用法\n", preChar))
	sb.WriteString(fmt.Sprintf("示例 %s帮助 厨师\n", preChar))
	sb.WriteString("\n")
	sb.WriteString("数据来源: L图鉴网\n")
	sb.WriteString("https://foodgame.gitee.io\n")
	sb.WriteString("项目地址（给个star呀）:\n")
	sb.WriteString("https://github.com/Billdex/bcjh-bot")
	return sb.String()
}

// 反馈功能指引
func feedbackHelp() string {
	var msg string
	msg += fmt.Sprintf("[问题反馈与建议]\n")
	msg += fmt.Sprintf("在使用过程中如果遇到了什么bug或者有什么好的建议，可以通过该功能反馈给我\n")
	msg += fmt.Sprintf("反馈方式:\n")
	msg += fmt.Sprintf("%s反馈 问题描述或建议\n", util.PrefixCharacters[0])
	msg += fmt.Sprintf("如果比较紧急可以私聊我们:\n")
	msg += fmt.Sprintf("QQ:591404144(罗觉觉)或646792290(汪汪泥)")
	return msg
}

// 图鉴网功能指引
func galleryWebsiteHelp() string {
	var msg string
	msg += fmt.Sprintf("[图鉴网-网址查询]\n")
	msg += fmt.Sprintf("给出L图鉴网与手机版图鉴网地址，方便记不住网址的小可爱快速访问。")
	return msg
}

// 游戏术语
func termHelp() string {
	var msg string
	//msg += fmt.Sprintf("[爆炒江湖游戏术语]\n")
	//msg += fmt.Sprintf("[技法]: 游戏内有炒、烤、煮、蒸、炸、切共6种技法，每道菜有1或2种技法属性，厨师技法达到菜谱所有要求即可做这道菜。\n")
	//msg += fmt.Sprintf("[菜谱品阶]: 菜谱有可、优、特、神、传5个品阶, 当做菜的厨师达到菜谱所有技法要求的1、2、3、4、5倍时可达成对应品阶。" +
	//	"不同品阶的金币收益不同，分别为100%%(可),110%%(优),130%%(特),150%%(神),200%%(传)\n")
	//msg += fmt.Sprintf("[熟练/专精]: 营业制作的菜会增加菜谱熟练度, 品阶越高增加越多, 熟练度满后会提高一定售价。\n")
	//msg += fmt.Sprintf("[碰瓷/升阶贵客]: 当一道菜首次提升至优、特、神时必来贵客(不需要做完整一组, 只需做一份即可), " +
	//	"注意, 碰瓷过高品阶贵客后无法再碰瓷低品阶贵客, 如直接做到神级后, 便无法再碰瓷优和特品阶的贵客。")
	termImagePath := config.AppConfig.Resource.Image + "/游戏术语.jpg"
	CQImage := bot.GetCQImage(termImagePath, "file")
	msg += fmt.Sprintf("%s", CQImage)
	return msg
}

// 厨师功能指引
func chefHelp() string {
	preChar := util.PrefixCharacters[0]
	var msg string
	msg += fmt.Sprintf("[厨师信息查询]\n")
	msg += fmt.Sprintf("[基础信息查询]: %s厨师 厨师名\n", preChar)
	msg += fmt.Sprintf("示例: %s厨师 羽十六", preChar)
	return msg
}

// 厨具功能指引
func equipmentHelp() string {
	preChar := util.PrefixCharacters[0]
	var msg string
	msg += fmt.Sprintf("[厨具信息查询]\n")
	msg += fmt.Sprintf("[基础信息查询]: %s厨具 厨具名\n", preChar)
	msg += fmt.Sprintf("示例: %s厨具 金烤叉", preChar)
	return msg
}

// 菜谱功能指引
func recipeHelp() string {
	preChar := util.PrefixCharacters[0]
	var msg string
	msg += fmt.Sprintf("[菜谱信息查询]\n")
	msg += fmt.Sprintf("[基础信息查询]: %s菜谱 菜谱名\n", preChar)
	msg += fmt.Sprintf("示例: %s菜谱 荷包蛋\n", preChar)
	msg += fmt.Sprintf("[复合信息查询]:\n")
	msg += fmt.Sprintf("%s菜谱 查询条件-参数-筛选条件\n", preChar)
	msg += fmt.Sprintf("示例: %s菜谱 食材-茄子-单时间-$100-p2\n", preChar)
	msg += fmt.Sprintf("%s菜谱 任意-耗材效率-4火\n", preChar)
	msg += fmt.Sprintf("目前提供以下查询条件:\n")
	msg += fmt.Sprintf("任意(不填参数), 食材, 技法, 贵客, 符文, 来源\n")
	msg += fmt.Sprintf("目前提供以下筛选条件(可叠加):\n")
	msg += fmt.Sprintf("单价下限($100), 稀有度下限(3火/星), 页数(p5)\n")
	msg += fmt.Sprintf("目前提供以下排序方式\n")
	msg += fmt.Sprintf("单时间, 总时间, 单价, 金币效率, ")
	msg += fmt.Sprintf("耗材效率, 食材效率")
	return msg
}

// 贵客功能指引
func guestHelp() string {
	preChar := util.PrefixCharacters[0]
	var msg string
	msg += fmt.Sprintf("[贵客信息查询]\n")
	msg += fmt.Sprintf("[基础信息查询]: %s贵客 贵客名\n", preChar)
	msg += fmt.Sprintf("示例: %s贵客 如来", preChar)
	return msg
}

// 符文功能指引
func antiqueHelp() string {
	preChar := util.PrefixCharacters[0]
	var msg string
	msg += fmt.Sprintf("[符文信息查询]\n")
	msg += fmt.Sprintf("提供根据符文名查询对应菜谱的功能, 并按照一组时间升序排序\n")
	msg += fmt.Sprintf("[基础信息查询]: %s符文 符文名\n", preChar)
	msg += fmt.Sprintf("当结果过多时可以使用「p」参数分页\n")
	msg += fmt.Sprintf("示例: %s符文 五香果 %s符文 一昧真火-p2", preChar, preChar)
	return msg
}
