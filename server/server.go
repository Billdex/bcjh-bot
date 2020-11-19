package server

import (
	"bcjh-bot/logger"
	"bcjh-bot/models"
	"bcjh-bot/service"
	"bcjh-bot/util"
	"encoding/json"
	"net/http"
)

//启动服务
func Run(port string) error {
	if "" == port {
		port = ":5800"
	}

	service.RegisterInstructions()
	logger.Info("指令处理函数注册完毕")

	http.HandleFunc("/", OneBotMsgHandler)
	return http.ListenAndServe(port, nil)
}

func OneBotMsgHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debug("收到一条OneBot消息：", r.Body)
	var msg models.OneBotMsg
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		logger.Error("数据格式有误", err)
		return
	}

	switch msg.MessageType {
	case util.OneBotMessageEvent:
		MessageEventHandler(&msg)
	case util.OneBotNoticeEvent:
		NoticeEventHandler(&msg)
	case util.OneBotRequestEvent:
		RequestEventHandler(&msg)
	case util.OneBotMetaEvent:
		MetaEventHandler(&msg)
	default:
		logger.Info("未知OneBot消息类型:", msg.MessageType)
	}

	w.Header().Set("Content-Type", "application/json")
	return
}
