package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct{}

func NewServer() *Server {
	return new(Server)
}

func (Server) Run(port string) error {
	http.HandleFunc("/", MsgHandler)
	if "" == port {
		port = ":5800"
	}
	return http.ListenAndServe(port, nil)
}

func MsgHandler(w http.ResponseWriter, r *http.Request) {
	var msg CQHttpMsg
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("数据有误")
		return
	}

	fmt.Printf("Msg:%+v\n", msg)
}
