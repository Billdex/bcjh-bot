package service

import (
	"bcjh-bot/model/onebot"
	"bcjh-bot/util/logger"
)

func EquipmentQuery(c *onebot.Context, args []string) {
	logger.Info("厨具查询，参数:", args)
}
