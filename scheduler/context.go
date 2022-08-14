package scheduler

import (
	"bcjh-bot/scheduler/onebot"
	"errors"
	"regexp"
	"strconv"
)

type Context struct {
	scheduler *Scheduler
	keyword   string
	handlers  []HandleFunc
	index     int

	bot               *onebot.Bot
	event             interface{}
	privateEvent      *onebot.MessageEventPrivateReq
	groupEvent        *onebot.MessageEventGroupReq
	messageType       string
	rawMessage        string
	PretreatedMessage string
	WarnMessage       string
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

func (c *Context) Reply(msg string) (int32, error) {
	if c.bot == nil {
		return -1, errors.New("bot信息未记录或连接已断开")
	}
	switch c.messageType {
	case onebot.MessageTypePrivate:
		if c.privateEvent.SubType == "friend" {
			return c.bot.SendPrivateMessage(c.privateEvent.Sender.UserId, msg)
		}
	case onebot.MessageTypeGroup:
		return c.bot.SendGroupMessage(c.groupEvent.GroupId, msg)
	default:
		return 0, nil
	}
	return 0, nil
}

func (c *Context) GetBot() *onebot.Bot {
	return c.bot
}

func (c *Context) GetBotId() int64 {
	if c.bot != nil {
		return c.bot.BotId
	}
	return 0
}

func (c *Context) GetLinkBotList() []*onebot.Bot {
	return c.scheduler.Engine.GetBots()
}

func (c *Context) GetKeyword() string {
	return c.keyword
}

func (c *Context) GetEvent() interface{} {
	return c.event
}

func (c *Context) GetPrivateEvent() *onebot.MessageEventPrivateReq {
	return c.privateEvent
}

func (c *Context) GetGroupEvent() *onebot.MessageEventGroupReq {
	return c.groupEvent
}

func (c *Context) GetMessageType() string {
	return c.messageType
}

func (c *Context) GetRawMessage() string {
	return c.rawMessage
}

func (c *Context) SetWarnMessage(msg string) {
	c.WarnMessage = msg
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

func (c *Context) GetGroupId() int64 {
	if c.messageType == onebot.MessageTypeGroup && c.GetGroupEvent() != nil {
		return c.GetGroupEvent().GroupId
	}
	return 0
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

func (c *Context) GetAtList() []int64 {
	atList := make([]int64, 0)
	reg := `\[CQ:at,qq=(\d+)\]`
	pattern := regexp.MustCompile(reg)
	allIndexes := pattern.FindAllSubmatch([]byte(c.PretreatedMessage), -1)
	for _, loc := range allIndexes {
		qq, _ := strconv.ParseInt(string(loc[1]), 10, 64)
		atList = append(atList, qq)
	}
	return atList
}

func (c *Context) GetImageList() []string {
	imgList := make([]string, 0)
	reg := `\[CQ:image,file=(.*?.image).*?\]`
	pattern := regexp.MustCompile(reg)
	allIndexes := pattern.FindAllSubmatch([]byte(c.PretreatedMessage), -1)
	for _, loc := range allIndexes {
		img := string(loc[1])
		imgList = append(imgList, img)
	}
	return imgList
}
