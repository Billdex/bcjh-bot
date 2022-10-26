package messageservice

import (
	"bcjh-bot/dao"
	"bcjh-bot/scheduler"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"fmt"
	"strings"
)

var prefixCharacters = []string{"#", "＃"}

func HelpGuide(c *scheduler.Context) {
	_, _ = c.Reply(introHelp())
}

// 功能指引
func introHelp() string {
	preChar := prefixCharacters[0]
	sb := strings.Builder{}
	sb.WriteString("【爆炒江湖查询机器人】\n")
	sb.WriteString(fmt.Sprintf("使用方式『%s功能名 参数』\n", preChar))
	sb.WriteString(fmt.Sprintf("示例「%s厨师 羽十六」\n", preChar))
	sb.WriteString("\n")
	sb.WriteString("详情请看说明文档:\n")
	sb.WriteString("http://bcjhbot.billdex.cn\n")
	sb.WriteString("数据来源: L图鉴网\n")
	sb.WriteString("https://foodgame.github.io")
	return sb.String()
}

// 反馈功能指引
func feedbackHelp() string {
	var msg string
	msg += fmt.Sprintf("【问题反馈与建议】\n")
	msg += fmt.Sprintf("在使用过程中如果遇到了什么bug或者有什么好的建议，可以通过该功能反馈给我\n")
	msg += fmt.Sprintf("反馈方式:\n")
	msg += fmt.Sprintf("「%s反馈 问题描述或建议」\n", prefixCharacters[0])
	msg += fmt.Sprintf("如果比较紧急可以私聊我们:\n")
	msg += fmt.Sprintf("QQ:591404144(罗觉觉)或646792290(汪汪泥)")
	return msg
}

// 厨师功能指引
func chefHelp() string {
	preChar := prefixCharacters[0]
	var msg string
	msg += fmt.Sprintf("【厨师信息查询】\n")
	msg += fmt.Sprintf("基础信息查询:『%s厨师 厨师名』\n", preChar)
	msg += fmt.Sprintf("示例:「%s厨师 羽十六」", preChar)
	return msg
}

// 厨具功能指引
func equipmentHelp() string {
	preChar := prefixCharacters[0]
	var msg string
	msg += fmt.Sprintf("【厨具信息查询】\n")
	msg += fmt.Sprintf("基础信息查询:『%s厨具 厨具名』\n", preChar)
	msg += fmt.Sprintf("示例:「%s厨具 金烤叉」", preChar)
	return msg
}

// 菜谱功能指引
func recipeHelp() string {
	preChar := prefixCharacters[0]
	var msg string
	msg += fmt.Sprintf("【菜谱信息查询】\n")
	msg += fmt.Sprintf("基础信息查询:『%s菜谱 菜谱名』\n", preChar)
	msg += fmt.Sprintf("示例:「%s菜谱 荷包蛋」\n", preChar)
	msg += fmt.Sprintf("复合信息查询说明请看文档:\n")
	msg += "http://bcjhbot.billdex.cn/#/usage/recipe"
	return msg
}

// 食材及效率查询
func materialHelp() string {
	prefix := prefixCharacters[0]
	sb := strings.Builder{}
	sb.WriteString("【食材及食材效率查询】\n")
	sb.WriteString(fmt.Sprintf("查询方式:『%s食材 食材名』", prefix))
	return sb.String()
}

// 调料查询功能指引
func condimentHelp() string {
	prefix := prefixCharacters[0]
	sb := strings.Builder{}
	sb.WriteString("【调料信息查询】\n")
	sb.WriteString("1. 简单查询，接名称或ID:\n")
	sb.WriteString(fmt.Sprintf("『%s调料 香菜』『%s调料 1』\n", prefix, prefix))
	sb.WriteString("2. 限制稀有度，和菜名写在一起:\n")
	sb.WriteString(fmt.Sprintf("『%s调料 三火』『%s调料 三星香菜』『%s调料 3星香菜』\n", prefix, prefix, prefix))
	sb.WriteString("3. 限制来源，或对应阁楼的技法:\n")
	sb.WriteString(fmt.Sprintf("『%s调料 %%切』『%s调料 三火 切』『%s调料 1星 梵正』",
		prefix, prefix, prefix))
	// sb.WriteString("4. 限制技能:\n")
	// sb.WriteString(fmt.Sprintf("『%s调料 三火%s炒技法+15』『%s调料 三火%s采集』『%s调料 三火%s售价』\n",
	// 	prefix, split, prefix, split, prefix, split))
	return sb.String()
}

// 贵客功能指引
func guestHelp() string {
	preChar := prefixCharacters[0]
	var msg string
	msg += fmt.Sprintf("【贵客信息查询】\n")
	msg += fmt.Sprintf("基础信息查询:『%s贵客 贵客名』\n", preChar)
	msg += fmt.Sprintf("示例:「%s贵客 如来」", preChar)
	return msg
}

// 符文功能指引
func antiqueHelp() string {
	preChar := prefixCharacters[0]
	var msg string
	msg += fmt.Sprintf("【符文信息查询】\n")
	msg += fmt.Sprintf("提供根据符文名查询对应菜谱的功能, 并按照一组时间升序排序\n")
	msg += fmt.Sprintf("结果过多可使用「p」参数分页\n")
	msg += fmt.Sprintf("示例:『%s符文 五香果』『%s符文 一昧真火 p2』", preChar, preChar)
	return msg
}

// 任务功能指引
func questHelp() string {
	prefix := prefixCharacters[0]
	maxLen := 5
	sb := strings.Builder{}
	sb.WriteString("【任务信息查询】\n")
	sb.WriteString(fmt.Sprintf("1. 主线，接ID（可以指定长度，最多%d条）:\n", maxLen))
	sb.WriteString(fmt.Sprintf("『%s主线 1』『%s主线 1 5』\n", prefix, prefix))
	sb.WriteString("2. 支线，接ID:\n")
	sb.WriteString(fmt.Sprintf("『%s支线 9.1』\n", prefix))
	sb.WriteString("3. 限时，使用#攻略 限时 查看限时任务攻略")
	return sb.String()
}

// 碰瓷功能指引
func upgradeGuestHelp() string {
	prefix := prefixCharacters[0]
	sb := strings.Builder{}
	sb.WriteString("【升阶贵客查询】\n")
	sb.WriteString(fmt.Sprintf("查询碰瓷贵客可用的菜:\n"))
	sb.WriteString("结果过多可使用「p」参数分页\n")
	sb.WriteString(fmt.Sprintf("示例:『%s碰瓷 如来』『%s碰瓷 唐伯虎 p2』", prefix, prefix))
	return sb.String()

}

// 后厨合成菜功能指引
func comboHelp() string {
	prefix := prefixCharacters[0]
	sb := strings.Builder{}
	sb.WriteString("【后厨合成菜谱查询】\n")
	sb.WriteString(fmt.Sprintf("查询后厨合成菜的前置菜谱:\n"))
	sb.WriteString(fmt.Sprintf("示例:『%s后厨 BBQ烧烤』", prefix))
	return sb.String()
}

// LaboratoryHelp 实验室前置功能指引
func LaboratoryHelp() string {
	prefix := prefixCharacters[0]
	sb := strings.Builder{}
	sb.WriteString("【实验室菜谱查询】\n")
	sb.WriteString(fmt.Sprintf("查询实验室菜谱的前置材料, 暂未收录厨具数据\n"))
	sb.WriteString(fmt.Sprintf("示例:『%s实验室 猪肉脯』", prefix))
	return sb.String()
}

// 攻略功能指引
func strategyHelp() string {
	keywords, err := dao.LoadStrategyKeywords()
	if err != nil {
		logger.Errorf("查询攻略关键词列表失败 %v", err)
		return e.SystemErrorNote
	}
	sb := strings.Builder{}
	sb.WriteString("【游戏攻略快捷查询】\n")
	sb.WriteString("收录了一些简要的游戏攻略，查询方式:『#攻略 关键词』\n")
	sb.WriteString("目前收录了以下内容:\n")
	for _, keyword := range keywords {
		sb.WriteString(fmt.Sprintf("「%s」 ", keyword))
	}
	return sb.String()
}
