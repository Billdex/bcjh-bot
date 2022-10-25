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

// ComboQuery 后厨合成菜查询
func ComboQuery(c *scheduler.Context) {
	comboRecipeName := c.PretreatedMessage

	// 查询全部菜谱
	allRecipes, err := dao.FindAllRecipes()
	if err != nil {
		logger.Error("获取菜谱数据出错", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}

	// 匹配符合的后厨菜名和对应需要的前置菜谱
	mMatchRecipes := make(map[string][]database.Recipe)
	var matchComboName string
	for _, recipe := range allRecipes {
		for _, combo := range recipe.Combo {
			if strings.Contains(combo, comboRecipeName) {
				mMatchRecipes[combo] = append(mMatchRecipes[combo], recipe)
				matchComboName = combo
				break
			}
		}
	}

	if len(mMatchRecipes) == 0 {
		_, _ = c.Reply(fmt.Sprintf("%s 不是后厨合成菜哦", comboRecipeName))
		return
	}

	if len(mMatchRecipes) > 1 {
		msg := "你想查询哪个菜谱呢"
		for combo := range mMatchRecipes {
			msg += fmt.Sprintf("\n%s", combo)
		}
		_, _ = c.Reply(msg)
		return
	}

	recipes := mMatchRecipes[matchComboName]
	msg := fmt.Sprintf("合成%s需要以下菜谱", matchComboName)
	for _, recipe := range recipes {
		msg += fmt.Sprintf("\n「%s-%s」", recipe.Name, recipe.Origin)
	}

	_, _ = c.Reply(msg)
}
