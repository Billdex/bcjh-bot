package bot

import "golang.org/x/net/websocket"

type WsConnection struct {
	Conn     *websocket.Conn
	SendChan chan []byte
	RecvChan chan []byte
	isClose  bool
}

const (
	MAX_MESSAGE_NUM = 1000
)

func NewWsConnection(conn *websocket.Conn) (*WsConnection, error) {
	c := &WsConnection{
		Conn:     conn,
		SendChan: make(chan []byte, MAX_MESSAGE_NUM),
		RecvChan: make(chan []byte, MAX_MESSAGE_NUM),
		isClose:  false,
	}
}

func (c WsConnection) Close() error {
	err := c.Conn.Close()
	if err != nil {
		return err
	}
	c.isClose = true
	return nil
}
