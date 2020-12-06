package server

import (
	"bcjh-bot/model/onebot"
	"bcjh-bot/service"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
)

// 处理消息事件
func MessageEventHandler(c *onebot.Context) {
	// 判断前缀
	var text string
	var hasPrefix bool
	for _, prefix := range util.PrefixCharacters {
		text, hasPrefix = PrefixFilter(c.RawMessage, prefix)
		if hasPrefix {
			break
		}
	}
	if !hasPrefix {
		return
	}
	logger.Debugf("收到一条消息事件信息Msg:%+v\n正文内容:%v\n", c, text)

	// 分发指令
	instruction, args := InstructionFilter(text, service.Ins.GetInstructions())
	logger.Debugf("instruction:%v, args:%v\n", instruction, args)
	if instruction != nil {
		instruction(c, args)
	}
}

// 处理通知事件
func NoticeEventHandler(c *onebot.Context) {
	logger.Info("收到一条通知事件信息:", c.NoticeType)
}

// 处理请求事件
func RequestEventHandler(c *onebot.Context) {
	logger.Info("收到一条请求事件信息:", c.RequestType)
}

// 处理元事件
func MetaEventHandler(c *onebot.Context) {
	logger.Info("收到一条元事件信息:", c.MetaEventType)
}
