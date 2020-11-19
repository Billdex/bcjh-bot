package service

import (
	"bcjh-bot/logger"
	"bcjh-bot/models"
)

func UpdateData(msg *models.OneBotMsg, args []string) {
	logger.Info("更新数据, 参数:", args)
}
