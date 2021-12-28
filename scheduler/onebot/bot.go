package onebot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"math/rand"
	"sync"
	"time"
)

const (
	apiResponseMapKey = "%s_%d_%d%d"
)

type Bot struct {
	BotId         int64
	Session       *WsConnection
	AutoEscape    bool
	OnebotHandler Handler
	apiResponse   map[string]chan []byte
	apiResMux     sync.Mutex
}

type Handler struct {
	HandlePrivateMessage func(bot *Bot, req *MessageEventPrivateReq)
	HandleGroupMessage   func(bot *Bot, req *MessageEventGroupReq)

	HandleGroupRecallNotice func(bot *Bot, req *NoticeEventGroupRecall)
}

type MessageHandler func(bot *Bot, data []byte)

func NewBot(botId int64, conn *websocket.Conn, handler Handler, ch onCloseHandler, autoEscape bool) *Bot {
	bot := &Bot{
		BotId:         botId,
		AutoEscape:    autoEscape,
		OnebotHandler: handler,
		apiResponse:   make(map[string]chan []byte),
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
		if echo := gjson.Get(string(data), "echo").String(); echo != "" {
			bot.apiResMux.Lock()
			if ch, ok := bot.apiResponse[echo]; ok {
				ch <- data
			}
			bot.apiResMux.Unlock()
		}
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
	noticeType := gjson.Get(string(data), "notice_type").String()
	switch noticeType {
	case NoticeTypeGroupRecall:
		if bot.OnebotHandler.HandleGroupRecallNotice != nil {
			req := &NoticeEventGroupRecall{}
			if err := json.Unmarshal(data, req); err != nil {
				log.Errorf("解析群撤回消息出错: %v, 原始json数据: %s\n", err, string(data))
				return
			}
			bot.OnebotHandler.HandleGroupRecallNotice(bot, req)
		}

	}
}

func (bot *Bot) handleRequestEvent(data []byte) {

}

func (bot *Bot) handleMetaEvent(data []byte) {

}

func (bot *Bot) ActionRequestAPI(action string, params interface{}) ([]byte, error) {
	req := actionApiReq{
		Action: action,
		Params: params,
	}
	key := fmt.Sprintf(apiResponseMapKey, action, bot.BotId, time.Now().UnixNano(), rand.Intn(100))
	req.Echo = key
	recvChan := make(chan []byte, 1)
	bot.apiResMux.Lock()
	bot.apiResponse[key] = recvChan
	bot.apiResMux.Unlock()
	defer func() {
		bot.apiResMux.Lock()
		delete(bot.apiResponse, key)
		bot.apiResMux.Unlock()
		close(recvChan)
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	sendMsg, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}
	err = bot.Session.Send(sendMsg)
	if err != nil {
		return nil, err
	}
	select {
	case data := <-recvChan:
		return []byte(gjson.Get(string(data), "data").String()), nil
	case <-ctx.Done():
		return nil, errors.New("超时未收到返回数据")
	}
}

func (bot *Bot) GetMsgInfo(messageId int32) (MsgInfo, error) {
	params := getMsgInfoParams{
		MessageId: messageId,
	}
	data, err := bot.ActionRequestAPI("get_msg", params)
	if err != nil {
		return MsgInfo{}, err
	}
	msg := MsgInfo{}
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return MsgInfo{}, err
	}
	return msg, nil
}

func (bot *Bot) GetGroupInfo(groupId int64) (GroupInfo, error) {
	params := getGroupInfoParams{
		GroupId: groupId,
		NoCache: false,
	}
	data, err := bot.ActionRequestAPI("get_group_info", params)
	if err != nil {
		return GroupInfo{}, err
	}
	var group GroupInfo
	err = json.Unmarshal(data, &group)
	if err != nil {
		return GroupInfo{}, err
	}
	return group, nil
}

func (bot *Bot) GetGroupList() ([]GroupInfo, error) {
	var params interface{}
	data, err := bot.ActionRequestAPI("get_group_list", params)
	if err != nil {
		return nil, err
	}
	groups := make([]GroupInfo, 0)
	err = json.Unmarshal(data, &groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (bot *Bot) GetImageInfo(file string) (ImageInfo, error) {
	params := getImageParams{
		File: file,
	}
	data, err := bot.ActionRequestAPI("get_image", params)
	if err != nil {
		return ImageInfo{}, err
	}
	var image ImageInfo
	err = json.Unmarshal(data, &image)
	if err != nil {
		return ImageInfo{}, err
	}
	return image, nil
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
