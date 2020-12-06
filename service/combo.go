package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/model/database"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"fmt"
)

func ComboQuery(c *onebot.Context, args []string) {
	logger.Info("后厨合成菜前置菜谱查询，参数:", args)

	if len(args) == 0 {
		err := bot.SendMessage(c, comboHelp())
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
		return
	}

	comboRecipeName := args[0]

	preRecipes := make([]database.Recipe, 0)
	err := database.DB.Where("combo = ?", comboRecipeName).Find(&preRecipes)
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
			msg += fmt.Sprintf("\n「%s-%s-%s」", recipe.GalleryId, recipe.Name, recipe.Origin)
		}
	}

	err = bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}
