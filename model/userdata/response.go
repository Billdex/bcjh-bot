package userdata

import (
	"time"
)

type Response struct {
	Result     bool      `json:"result"`
	Id         int       `json:"id"`
	User       string    `json:"user"`
	Data       string    `json:"data"`
	CreateTime time.Time `json:"create_time"`
	Msg        string    `json:"msg"`
}
