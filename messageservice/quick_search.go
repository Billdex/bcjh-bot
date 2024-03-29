package messageservice

import (
	"bcjh-bot/dao"
	"bcjh-bot/model/database"
	"bcjh-bot/scheduler"
	"fmt"
	"sync"
)

func QuickSearch(c *scheduler.Context) {
	param := c.PretreatedMessage
	if param == "" || param == "%" {
		return
	}

	var wg sync.WaitGroup
	var recipes []database.Recipe
	var chefs []database.Chef
	var equips []database.Equip
	var strategies []database.Strategy
	wg.Add(4)
	go func() {
		recipes, _ = dao.SearchRecipesWithName(param)
		wg.Done()
	}()
	go func() {
		chefs, _ = dao.SearchChefsWithName(param)
		wg.Done()
	}()
	go func() {
		equips, _ = dao.SearchEquipsWithName(param)
		wg.Done()
	}()
	go func() {
		strategies, _ = dao.SearchStrategiesWithKeyword(param)
		wg.Done()
	}()
	wg.Wait()

	// 查询到多条结果的时候按顺序输出，每种三个，其他情况输出单个结果
	// 但是如果出现了完全匹配的结果数据，则直接输出该条数据作为结果
	var msg string
	total := len(recipes) + len(chefs) + len(equips) + len(strategies)
	if total > 1 {
		msg += "查询到以下结果:"
		for i := range recipes {
			if recipes[i].Name == param {
				_, _ = c.Reply(echoRecipeMessage(recipes[i]))
				return
			}
			if i >= len(recipes)*6/total && i > 1 {
				msg += "\n......"
				break
			}
			msg += fmt.Sprintf("\n菜谱 %s", recipes[i].Name)
		}
		for i := range chefs {
			if chefs[i].Name == param {
				_, _ = c.Reply(echoChefMessage(chefs[i]))
				return
			}
			if i >= len(chefs)*6/total && i > 1 {
				msg += "\n......"
				break
			}
			msg += fmt.Sprintf("\n厨师 %s", chefs[i].Name)
		}
		for i := range equips {
			if equips[i].Name == param {
				_, _ = c.Reply(echoEquipMessage(equips[i]))
				return
			}
			if i >= len(equips)*6/total && i > 1 {
				msg += "\n......"
				break
			}
			msg += fmt.Sprintf("\n厨具 %s", equips[i].Name)
		}
		for i := range strategies {
			if strategies[i].Keyword == param {
				_, _ = c.Reply(strategies[i].Value)
				return
			}
			if i >= len(strategies)*6/total && i > 1 {
				msg += "\n......"
				break
			}
			msg += fmt.Sprintf("\n攻略 %s", strategies[i].Keyword)
		}
	} else {
		if len(recipes) == 1 {
			msg = echoRecipeMessage(recipes[0])
		} else if len(chefs) == 1 {
			msg = echoChefMessage(chefs[0])
		} else if len(equips) == 1 {
			msg = echoEquipMessage(equips[0])
		} else if len(strategies) == 1 {
			msg = strategies[0].Value
		} else {
			msg = "没有找到相关结果!"
		}
	}
	_, _ = c.Reply(msg)
}
