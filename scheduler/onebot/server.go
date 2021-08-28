package onebot

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Server struct {
	Bots map[int64]*Bot
	Port string
	Path string
}

func New(port string, path string) *Server {
	s := &Server{
		Bots: make(map[int64]*Bot),
		Port: port,
		Path: path,
	}
	return s
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
			log.Infof("Bot %d 断开连接\n", botId)
			delete(s.Bots, botId)
		}
		s.Bots[botId] = NewBot(botId, c, handler, onCloseHandler, false)
		log.Infof("Bot %d 已建立连接\n", botId)
		return
	})
	return http.ListenAndServe(s.Port, nil)
}
