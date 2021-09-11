package messageservice

import (
	"bcjh-bot/model/database"
	"bcjh-bot/scheduler"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"fmt"
	"strconv"
	"strings"
)

func ComboQuery(c *scheduler.Context) {
	arg := strings.TrimSpace(c.PretreatedMessage)

	if arg == "" {
		_, _ = c.Reply(comboHelp())
		return
	}

	comboRecipeName := arg

	// 判断菜名是否唯一
	recipes := make([]database.Recipe, 0)
	recipeId, err := strconv.Atoi(comboRecipeName)
	if err == nil {
		err = database.DB.Where("gallery_id = ?", fmt.Sprintf("%03d", recipeId)).Find(&recipes)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			_, _ = c.Reply(e.SystemErrorNote)
			return
		}
	} else {
		err = database.DB.Where("name like ?", "%"+comboRecipeName+"%").Find(&recipes)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			_, _ = c.Reply(e.SystemErrorNote)
			return
		}
	}
	if len(recipes) == 0 {
		_, _ = c.Reply("没有查询到相关餐谱呢")
		return
	}
	if len(recipes) > 1 {
		msg := "你想查询哪个菜谱呢:"
		for _, recipe := range recipes {
			msg += fmt.Sprintf("\n%s", recipe.Name)
		}
		_, _ = c.Reply(msg)
		return
	}

	comboRecipeName = recipes[0].Name

	preRecipes := make([]database.Recipe, 0)
	err = database.DB.Where("combo = ?", comboRecipeName).Find(&preRecipes)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}
	var msg string
	if len(preRecipes) == 0 {
		msg = fmt.Sprintf("%s不是后厨合成菜哦!", comboRecipeName)
	} else {
		msg += fmt.Sprintf("合成%s需要以下前置菜谱:", comboRecipeName)
		for _, recipe := range preRecipes {
			msg += fmt.Sprintf("\n「%s-%s」", recipe.Name, recipe.Origin)
		}
	}

	_, _ = c.Reply(msg)
}
