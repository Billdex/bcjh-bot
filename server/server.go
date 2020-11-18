package server

import (
	"bcjh-bot/logger"
	"bcjh-bot/util"
	"encoding/json"
	"net/http"
)

func Run(port string) error {
	http.HandleFunc("/", MsgHandler)
	if "" == port {
		port = ":5800"
	}
	return http.ListenAndServe(port, nil)
}

func MsgHandler(w http.ResponseWriter, r *http.Request) {
	var msg CQHttpMsg
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		logger.Info("数据格式有误", err)
		return
	}

	text, hasPrefix := util.PrefixFilter(msg.RawMessage, util.PrefixCharacter)
	if !hasPrefix {
		return
	}
	logger.Infof("收到一条消息Msg:%+v\n正文内容:%v\n", msg, text)

	w.Header().Set("Content-Type", "application/json")
	return
}
