package server

import "time"

type CQHttpMsg struct {
	Time        time.Time `json:"time"`
	SelfId      int       `json:"self_id"`
	PostType    string    `json:"post_type"`
	MessageType string    `json:"message_type"`
	SubType     string    `json:"sub_type"`
	MessageId   int       `json:"message_id"`
	UserId      int       `json:"user_id"`
	Message     string    `json:"message"`
	RawMessage  string    `json:"raw_message"`
	Font        int       `json:"font"`
	Sender      struct {
		Nickname string `json:"nickname"`
		Sex      string `json:"sex"`
		Age      int    `json:"age"`
	} `json:"sender"`
}
