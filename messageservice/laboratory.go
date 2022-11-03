package messageservice

import (
	"bcjh-bot/dao"
	"bcjh-bot/model/database"
	"bcjh-bot/scheduler"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"fmt"
	"strings"
)

func LaboratoryQuery(c *scheduler.Context) {
	arg := c.PretreatedMessage
	laboratories, err := dao.FindAllLaboratory()
	if err != nil {
		logger.Error("查询实验室数据出错!", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}

	targets := make([]database.Laboratory, 0)
	for _, item := range laboratories {
		if strings.Contains(item.Target, arg) {
			targets = append(targets, item)
		}
	}

	var msg string
	if len(targets) == 0 {
		msg = fmt.Sprintf("%s似乎不是实验室菜谱", arg)
	} else if len(targets) == 1 {
		msg += fmt.Sprintf("「%s」%s", targets[0].Skill, targets[0].Target)
		msg += fmt.Sprintf("\n%s", strings.Repeat("🔥", targets[0].Rarity))
		msg += fmt.Sprintf("\n消耗符文:\n「%s」*%d", targets[0].Antique, targets[0].AntiqueNumber)
		msg += fmt.Sprintf("\n消耗厨具: ")
		if len(targets[0].Equips) == 0 {
			msg += "无"
		} else {
			for _, equip := range targets[0].Equips {
				msg += fmt.Sprintf("\n「%s」", equip)
			}
		}
		msg += fmt.Sprintf("\n前置菜谱: ")
		if len(targets[0].Recipes) == 0 {
			msg += "无"
		} else {
			for _, recipe := range targets[0].Recipes {
				msg += fmt.Sprintf("\n「%s」", recipe)
			}
		}
	} else {
		msg += "找到以下多个实验室菜谱\n"
		for _, target := range targets {
			msg += fmt.Sprintf("\n%s", target.Target)
		}
	}

	_, _ = c.Reply(msg)
}
