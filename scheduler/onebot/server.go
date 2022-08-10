package onebot

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Server struct {
	bots map[int64]*Bot
	Port string
	Path string

	mux sync.Mutex
}

func New(port string, path string) *Server {
	s := &Server{
		bots: make(map[int64]*Bot),
		Port: port,
		Path: path,
	}
	return s
}

func (s *Server) GetBots() []*Bot {
	s.mux.Lock()
	defer s.mux.Unlock()
	bots := make([]*Bot, 0)
	for _, bot := range s.bots {
		bots = append(bots, bot)
	}
	return bots
}

func (s *Server) AddBot(id int64, bot *Bot) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.bots[id] = bot
}

func (s *Server) RemoveBot(id int64) {
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.bots, id)
}

func (s *Server) Serve(handler Handler) error {
	http.HandleFunc(s.Path, func(w http.ResponseWriter, req *http.Request) {
		xSelfId := req.Header.Get("x-self-id")
		botId, err := strconv.ParseInt(xSelfId, 10, 64)
		if err != nil {
			return
		}
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		c, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			return
		}
		onCloseHandler := func(code int, message string) {
			log.Printf("Bot %d 断开连接\n", botId)
			s.RemoveBot(botId)
		}
		s.AddBot(botId, NewBot(botId, c, handler, onCloseHandler, false))
		log.Printf("Bot %d 已建立连接\n", botId)
		return
	})
	return http.ListenAndServe(s.Port, nil)
}
