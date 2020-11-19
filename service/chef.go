package service

import (
	"bcjh-bot/logger"
	"bcjh-bot/model"
)

func ChefQuery(msg *model.OneBotMsg, args []string) {
	logger.Info("厨师查询，参数:", args)
}
