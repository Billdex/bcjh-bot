package onebot

import (
	"errors"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type WsConnection struct {
	Conn     *websocket.Conn
	SendChan chan []byte
	RecvChan chan []byte
	isClose  bool
}

const (
	MAX_MESSAGE_NUM = 1000
)

type onRecvHandler func(data []byte)

type onCloseHandler func(code int, message string)

func NewWsConnection(conn *websocket.Conn, rh onRecvHandler, ch onCloseHandler) *WsConnection {
	wsc := &WsConnection{
		Conn:     conn,
		SendChan: make(chan []byte, MAX_MESSAGE_NUM),
		RecvChan: make(chan []byte, MAX_MESSAGE_NUM),
		isClose:  false,
	}

	wsc.Conn.SetCloseHandler(func(code int, text string) error {
		ch(code, text)
		return nil
	})

	// 接收消息
	go func() {
		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				log.Errorf("failed to read message, err: %+v", err)
				_ = wsc.Close()
				return
			}
			if messageType == websocket.PingMessage {
				_ = wsc.Send([]byte("pong"))
				continue
			}
			go rh(data)
		}
	}()

	// 发送消息
	go func() {
		for message := range wsc.SendChan {
			if wsc.Conn == nil {
				return
			}
			err := wsc.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				_ = wsc.Close()
				return
			}
		}
	}()
	return wsc
}

func (wsc *WsConnection) Send(data []byte) error {
	if wsc.isClose {
		return errors.New("connection is closed")
	}
	wsc.SendChan <- data
	return nil
}

func (wsc *WsConnection) Close() error {
	err := wsc.Conn.Close()
	if err != nil {
		return err
	}
	wsc.isClose = true
	return nil
}
