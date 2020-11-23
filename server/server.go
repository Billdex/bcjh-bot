package server

import (
	"bcjh-bot/model/onebot"
	"bcjh-bot/service"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
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
	context := new(onebot.Context)
	err := json.NewDecoder(r.Body).Decode(context)
	if err != nil {
		logger.Error("数据格式有误", err)
		return
	}

	switch context.PostType {
	case util.OneBotMessageEvent:
		MessageEventHandler(context)
	case util.OneBotNoticeEvent:
		NoticeEventHandler(context)
	case util.OneBotRequestEvent:
		RequestEventHandler(context)
	case util.OneBotMetaEvent:
		MetaEventHandler(context)
	default:
		logger.Info("未知OneBot事件类型:", context.MessageType)
	}

	w.Header().Set("Content-Type", "application/json")
	return
}
