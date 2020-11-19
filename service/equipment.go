package service

import (
	"bcjh-bot/logger"
	"bcjh-bot/model"
)

func EquipmentQuery(msg *model.OneBotMsg, args []string) {
	logger.Info("厨具查询，参数:", args)
}
