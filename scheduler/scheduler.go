package pbbot_scheduler

import (
	"errors"
	"github.com/ProtobufBot/go-pbbot"
	"github.com/ProtobufBot/go-pbbot/proto_gen/onebot"
	"net/http"
	"strings"
)

type Scheduler struct {
	*CmdGroup
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

func (s *Scheduler) Process(bot *pbbot.Bot, event interface{}) error {
	c := s.createContext()
	var rawMessage string
	if privateEvent, ok := event.(*onebot.PrivateMessageEvent); ok {
		c.privateMessageEvent = privateEvent
		rawMessage = privateEvent.RawMessage
	} else if groupEvent, ok := event.(*onebot.GroupMessageEvent); ok {
		c.groupMessageEvent = groupEvent
		rawMessage = groupEvent.RawMessage
	} else {
		return errors.New("event类型错误!必须为*onebot.PrivateMessageEvent或*onebot.GroupMessageEvent")
	}
	c.bot = bot
	c.rawMessage = rawMessage
	handlerChain, content, found := s.findHandler(rawMessage)
	if found {
		c.handlers = handlerChain
		c.PretreatedMessage = content
		for c.index < len(c.handlers) {
			c.handlers[c.index](c)
			c.index++
		}
	}
	return nil
}

func (s *Scheduler) findHandler(message string) ([]HandleFunc, string, bool) {
	return s.CmdGroup.SearchHandlerChain(strings.TrimSpace(message))
}

func (s *Scheduler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := s.createContext()
	s.Process(c)
}

func (s *Scheduler) Run(port string) error {
	if port == "" {
		port = ":5800"
	}
	http.HandleFunc("/", s.ServeHTTP)
	return http.ListenAndServe(port, nil)
}
