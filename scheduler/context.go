package pbbot_scheduler

import (
	"errors"
	"github.com/ProtobufBot/go-pbbot"
	"github.com/ProtobufBot/go-pbbot/proto_gen/onebot"
)

type Context struct {
	scheduler *Scheduler
	handlers  []HandleFunc
	index     int

	bot                 *pbbot.Bot
	privateMessageEvent *onebot.PrivateMessageEvent
	groupMessageEvent   *onebot.GroupMessageEvent
	rawMessage          string
	PretreatedMessage   string

	replyMessage *pbbot.Msg
}

func (c *Context) GetBot() *pbbot.Bot {
	return c.bot
}

func (c *Context) GetPrivateMessageEvent() (*onebot.PrivateMessageEvent, bool) {
	if c.privateMessageEvent == nil {
		return nil, false
	} else {
		return c.privateMessageEvent, true
	}
}

func (c *Context) GetGroupMessageEvent() (*onebot.GroupMessageEvent, bool) {
	if c.groupMessageEvent == nil {
		return nil, false
	} else {
		return c.groupMessageEvent, true
	}
}

func (c *Context) GetRawMessage() string {
	return c.rawMessage
}

func (c *Context) Reply(msg *pbbot.Msg, autoEscape bool) (int32, error) {
	if c.bot == nil {
		return -1, errors.New("context未记录Bot信息")
	}
	if c.privateMessageEvent != nil {
		resp, err := c.bot.SendPrivateMessage(c.privateMessageEvent.UserId, msg, autoEscape)
		if err != nil {
			return -1, err
		} else {
			return resp.MessageId, nil
		}
	} else if c.groupMessageEvent != nil {
		resp, err := c.bot.SendGroupMessage(c.groupMessageEvent.GroupId, msg, autoEscape)
		if err != nil {
			return -1, err
		} else {
			return resp.MessageId, nil
		}
	} else {
		return -1, errors.New("messageEvent信息有误")
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

func MustGroupAdmin(c *Context) {
	if event, ok := c.GetGroupMessageEvent(); ok {
		if event.Sender.Role != "owner" && event.Sender.Role != "admin" {
			c.Abort()
		}
	} else {
		c.Abort()
	}
}
