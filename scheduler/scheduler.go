package scheduler

import (
	"bcjh-bot/scheduler/onebot"
	"strings"
)

type Scheduler struct {
	*CmdGroup
	Engine *onebot.Server
}

type HandleFunc func(*Context)

func New() *Scheduler {
	scheduler := &Scheduler{
		CmdGroup: &CmdGroup{
			isHandleNode: false,
			ignoreCase:   false,
			Keywords:     []string{""},
			BaseHandlers: []HandleFunc{},
			subCmdGroups: make([]*CmdGroup, 0),
		},
	}
	scheduler.CmdGroup.scheduler = scheduler

	return scheduler
}

func (s *Scheduler) createContext() *Context {
	return &Context{
		scheduler: s,
		handlers:  make([]HandleFunc, 0),
		index:     0,
	}
}

func (s *Scheduler) Process(bot *onebot.Bot, event interface{}) {
	c := s.createContext()
	if privateEvent, ok := event.(*onebot.MessageEventPrivateReq); ok {
		c.privateEvent = privateEvent
		c.messageType = onebot.MessageTypePrivate
		c.rawMessage = privateEvent.RawMessage
	} else if groupEvent, ok := event.(*onebot.MessageEventGroupReq); ok {
		c.groupEvent = groupEvent
		c.messageType = onebot.MessageTypeGroup
		c.rawMessage = groupEvent.RawMessage
	} else {
		return
	}
	c.bot = bot
	c.event = event
	keyword, handlerChain, content, found := s.findHandler(c.rawMessage)
	if found {
		c.keyword = keyword
		c.handlers = handlerChain
		c.PretreatedMessage = content
		for c.index < len(c.handlers) {
			c.handlers[c.index](c)
			c.index++
		}
	}
}

func (s *Scheduler) findHandler(message string) (string, []HandleFunc, string, bool) {
	return s.CmdGroup.SearchHandlerChain(strings.TrimSpace(message))
}

func (s *Scheduler) Serve(port string, path string, handler *onebot.Handler) error {
	s.Engine = onebot.New(port, path)
	handler.HandlePrivateMessage = func(bot *onebot.Bot, req *onebot.MessageEventPrivateReq) {
		s.Process(bot, req)
	}
	handler.HandleGroupMessage = func(bot *onebot.Bot, req *onebot.MessageEventGroupReq) {
		s.Process(bot, req)
	}
	return s.Engine.Serve(*handler)
}
