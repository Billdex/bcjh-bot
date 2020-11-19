package service

import (
	"bcjh-bot/logger"
	"bcjh-bot/model/onebot"
)

func ChefQuery(c *onebot.Context, args []string) {
	logger.Info("厨师查询，参数:", args)
}
