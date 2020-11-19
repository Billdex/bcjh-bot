package service

import (
	"bcjh-bot/logger"
	"bcjh-bot/models"
)

func RecipeQuery(msg *models.OneBotMsg, args []string) {
	logger.Info("菜谱查询, 参数:", args)
}
