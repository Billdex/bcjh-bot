package scheduler

import (
	"bcjh-bot/scheduler/onebot"
	"errors"
)

type Context struct {
	scheduler *Scheduler
	handlers  []HandleFunc
	index     int

	bot               *onebot.Bot
	event             interface{}
	privateEvent      *onebot.MessageEventPrivateReq
	groupEvent        *onebot.MessageEventGroupReq
	messageType       string
	rawMessage        string
	PretreatedMessage string
}

func (c *Context) GetBot() *onebot.Bot {
	return c.bot
}

func (c *Context) GetEvent() interface{} {
	return c.event
}

func (c *Context) GetMessageType() string {
	return c.messageType
}

func (c *Context) GetRawMessage() string {
	return c.rawMessage
}

func (c *Context) GetEventTime() int64 {
	switch c.messageType {
	case onebot.MessageTypePrivate:
		return c.privateEvent.Time
	case onebot.MessageTypeGroup:
		return c.groupEvent.Time
	default:
		return 0
	}
}

func (c *Context) GetSenderId() int64 {
	switch c.messageType {
	case onebot.MessageTypePrivate:
		return c.privateEvent.Sender.UserId
	case onebot.MessageTypeGroup:
		return c.groupEvent.Sender.UserId
	default:
		return 0
	}
}

func (c *Context) GetSenderNickname() string {
	switch c.messageType {
	case onebot.MessageTypePrivate:
		return c.privateEvent.Sender.Nickname
	case onebot.MessageTypeGroup:
		return c.groupEvent.Sender.Nickname
	default:
		return ""
	}
}

func (c *Context) Reply(msg string) (int32, error) {
	if c.bot == nil {
		return -1, errors.New("bot信息未记录或连接已断开")
	}
	switch c.messageType {
	case onebot.MessageTypePrivate:
		return c.bot.SendPrivateMessage(c.privateEvent.Sender.UserId, msg)
	case onebot.MessageTypeGroup:
		return c.bot.SendGroupMessage(c.groupEvent.GroupId, msg)
	default:
		return 0, nil
	}
}

func (c *Context) Next() {
	c.index++
	for c.index < (len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) IsAborted() bool {
	return c.index >= len(c.handlers)
}

func (c *Context) Abort() {
	c.index = len(c.handlers)
}
