package service

import (
	"bcjh-bot/logger"
	"bcjh-bot/model/onebot"
)

func RecipeQuery(c *onebot.Context, args []string) {
	logger.Info("菜谱查询, 参数:", args)
}
