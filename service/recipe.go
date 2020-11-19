package service

import (
	"bcjh-bot/logger"
	"bcjh-bot/model"
)

func RecipeQuery(msg *model.OneBotMsg, args []string) {
	logger.Info("菜谱查询, 参数:", args)
}
