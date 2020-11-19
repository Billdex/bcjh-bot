package service

import (
	"bcjh-bot/logger"
	"bcjh-bot/model/onebot"
)

func EquipmentQuery(c *onebot.Context, args []string) {
	logger.Info("厨具查询，参数:", args)
}
