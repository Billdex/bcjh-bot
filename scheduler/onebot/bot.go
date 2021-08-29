package onebot

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Bot struct {
	BotId         int64
	Session       *WsConnection
	AutoEscape    bool
	OnebotHandler Handler
}

type Handler struct {
	HandlePrivateMessage func(bot *Bot, req *MessageEventPrivateReq)
	HandleGroupMessage   func(bot *Bot, req *MessageEventGroupReq)
}

type MessageHandler func(bot *Bot, data []byte)

func NewBot(botId int64, conn *websocket.Conn, handler Handler, ch onCloseHandler, autoEscape bool) *Bot {
	bot := &Bot{
		BotId:         botId,
		AutoEscape:    autoEscape,
		OnebotHandler: handler,
	}
	onRecvHandler := bot.handleRecv
	bot.Session = NewWsConnection(conn, onRecvHandler, ch)
	return bot
}

func (bot *Bot) handleRecv(data []byte) {
	switch gjson.Get(string(data), "post_type").String() {
	case PostTypeMessageEvent:
		bot.handleMessageEvent(data)
	case PostTypeNoticeEvent:
		bot.handleNoticeEvent(data)
	case PostTypeRequestEvent:
		bot.handleRequestEvent(data)
	case PostTypeMetaEvent:
		bot.handleMetaEvent(data)
	default:
		return
	}
}

func (bot *Bot) handleMessageEvent(data []byte) {
	msgType := gjson.Get(string(data), "message_type").String()
	if msgType == MessageTypePrivate {
		req := &MessageEventPrivateReq{}
		if err := json.Unmarshal(data, req); err != nil {
			log.Errorf("解析私聊消息出错: %v, 原始json数据: %s\n", err, string(data))
			return
		}
		if bot.OnebotHandler.HandlePrivateMessage != nil {
			bot.OnebotHandler.HandlePrivateMessage(bot, req)
		}
	} else if msgType == MessageTypeGroup {
		req := &MessageEventGroupReq{}
		if err := json.Unmarshal(data, req); err != nil {
			log.Errorf("解析群消息出错: %v, 原始json数据: %s\n", err, string(data))
			return
		}
		if bot.OnebotHandler.HandleGroupMessage != nil {
			bot.OnebotHandler.HandleGroupMessage(bot, req)
		}
	}
}

func (bot *Bot) handleNoticeEvent(data []byte) {

}

func (bot *Bot) handleRequestEvent(data []byte) {

}

func (bot *Bot) handleMetaEvent(data []byte) {

}

func (bot *Bot) SendPrivateMessage(userId int64, message string) (int32, error) {
	req := actionApiReq{
		Action: "send_private_msg",
	}
	params := sendPrivateMsgParams{
		UserId:     userId,
		Message:    message,
		AutoEscape: bot.AutoEscape,
	}
	req.Params = params
	data, err := json.Marshal(&req)
	if err != nil {
		return 0, err
	}
	return 0, bot.Session.Send(data)
}

func (bot *Bot) SendGroupMessage(groupId int64, message string) (int32, error) {
	req := actionApiReq{
		Action: "send_group_msg",
	}
	params := sendGroupMsgParams{
		GroupId:    groupId,
		Message:    message,
		AutoEscape: bot.AutoEscape,
	}
	req.Params = params
	data, err := json.Marshal(&req)
	if err != nil {
		return 0, err
	}
	return 0, bot.Session.Send(data)
}
