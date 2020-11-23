package service

import (
	"bcjh-bot/model/onebot"
	"bcjh-bot/util/logger"
)

func RecipeQuery(c *onebot.Context, args []string) {
	logger.Info("菜谱查询, 参数:", args)
}
