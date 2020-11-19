package service

import (
	"bcjh-bot/logger"
	"bcjh-bot/models"
)

func EquipmentQuery(msg *models.OneBotMsg, args []string) {
	logger.Info("厨具查询，参数:", args)
}
