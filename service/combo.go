package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/model/database"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"fmt"
	"strconv"
)

func ComboQuery(c *onebot.Context, args []string) {
	logger.Info("后厨合成菜前置菜谱查询，参数:", args)

	if len(args) == 0 {
		//err := bot.SendMessage(c, comboHelp())
		//if err != nil {
		//	logger.Error("发送信息失败!", err)
		//}
		return
	}

	comboRecipeName := args[0]

	// 判断菜名是否唯一
	recipes := make([]database.Recipe, 0)
	recipeId, err := strconv.Atoi(comboRecipeName)
	if err == nil {
		err = database.DB.Where("gallery_id = ?", fmt.Sprintf("%03d", recipeId)).Find(&recipes)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			_ = bot.SendMessage(c, util.SystemErrorNote)
			return
		}
	} else {
		err = database.DB.Where("name like ?", "%"+comboRecipeName+"%").Find(&recipes)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			_ = bot.SendMessage(c, util.SystemErrorNote)
			return
		}
	}
	if len(recipes) == 0 {
		_ = bot.SendMessage(c, "没有查询到相关餐谱呢")
		return
	}
	if len(recipes) > 1 {
		msg := "你想查询哪个菜谱呢:"
		for _, recipe := range recipes {
			msg += fmt.Sprintf("\n%s", recipe.Name)
		}
		_ = bot.SendMessage(c, msg)
		return
	}

	comboRecipeName = recipes[0].Name

	preRecipes := make([]database.Recipe, 0)
	err = database.DB.Where("combo = ?", comboRecipeName).Find(&preRecipes)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_ = bot.SendMessage(c, util.SystemErrorNote)
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

	err = bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}
