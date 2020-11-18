package server

import (
	"bcjh-bot/logger"
	"bcjh-bot/service"
	"bcjh-bot/util"
	"encoding/json"
	"net/http"
)

func Run(port string) error {
	if "" == port {
		port = ":5800"
	}

	service.RegisterInstructions()
	logger.Info("指令处理函数注册完毕")

	http.HandleFunc("/", MsgHandler)
	return http.ListenAndServe(port, nil)
}

func MsgHandler(w http.ResponseWriter, r *http.Request) {
	var msg CQHttpMsg
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		logger.Error("数据格式有误", err)
		return
	}

	//判断前缀
	text, hasPrefix := service.PrefixFilter(msg.RawMessage, util.PrefixCharacter)
	if !hasPrefix {
		return
	}
	logger.Debugf("收到一条消息Msg:%+v\n正文内容:%v\n", msg, text)

	//分发指令
	instruction, args := service.InstructionFilter(text, service.Ins.GetInstructions())
	logger.Debugf("instruction:%v, args:%v", instruction, args)
	if instruction != nil {
		instruction(args)
	}

	w.Header().Set("Content-Type", "application/json")
	return
}
