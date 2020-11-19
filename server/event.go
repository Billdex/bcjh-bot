package server

import (
	"bcjh-bot/logger"
	"bcjh-bot/models"
	"bcjh-bot/service"
	"bcjh-bot/util"
)

//处理消息事件
func MessageEventHandler(msg *models.OneBotMsg) {
	//判断前缀
	text, hasPrefix := service.PrefixFilter(msg.RawMessage, util.PrefixCharacter)
	if !hasPrefix {
		return
	}
	logger.Debugf("收到一条消息事件信息Msg:%+v\n正文内容:%v\n", msg, text)

	//分发指令
	instruction, args := service.InstructionFilter(text, service.Ins.GetInstructions())
	logger.Debugf("instruction:%v, args:%v", instruction, args)
	if instruction != nil {
		instruction(msg, args)
	}
}

//处理通知事件
func NoticeEventHandler(msg *models.OneBotMsg) {
	logger.Info("收到一条通知事件信息:", msg.NoticeType)
}

//处理请求事件
func RequestEventHandler(msg *models.OneBotMsg) {
	logger.Info("收到一条请求事件信息:", msg.RequestType)
}

//处理元事件
func MetaEventHandler(msg *models.OneBotMsg) {
	logger.Info("收到一条元事件信息:", msg.MetaEventType)
}
