package messageservice

import (
	"bcjh-bot/model/database"
	"bcjh-bot/scheduler"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"fmt"
	"strings"
)

func LaboratoryQuery(c *scheduler.Context) {
	arg := strings.TrimSpace(c.PretreatedMessage)

	if arg == "" {
		_, _ = c.Reply(LaboratoryHelp())
		return
	}

	targets := make([]database.Laboratory, 0)
	err := database.DB.Where("target_name like ?", "%"+arg+"%").Find(&targets)
	if err != nil {
		logger.Error("数据库查询出错!")
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}

	var msg string
	if len(targets) == 0 {
		msg = fmt.Sprintf("%s似乎不是实验室菜谱", arg)
	} else if len(targets) == 1 {
		rarity := ""
		for i := 0; i < targets[0].Rarity; i++ {
			rarity += "🔥"
		}
		msg += fmt.Sprintf("「%s」%s", targets[0].Skill, targets[0].Target)
		msg += fmt.Sprintf("\n%s", rarity)
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
		var msg string
		msg += "找到以下多个实验室菜谱\n"
		for _, target := range targets {
			msg += fmt.Sprintf("\n%s", target.Target)
		}
	}

	_, _ = c.Reply(msg)
}
