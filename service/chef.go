package service

import (
	"bcjh-bot/logger"
	"bcjh-bot/models"
)

func ChefQuery(msg *models.OneBotMsg, args []string) {
	logger.Info("厨师查询，参数:", args)
}
