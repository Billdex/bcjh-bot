package database

import "time"

type Feedback struct {
	Id         int64          `xorm:"id pk autoincr" json:"id"`
	Sender     int64          `xorm:"sender BIGINT" json:"sender"`
	Nickname   string         `xorm:"nickname" json:"nickname"`
	Message    string         `xorm:"message" json:"message"`
	Status     FeedbackStatus `xorm:"status" json:"status"`
	CreateTime time.Time      `xorm:"create_time created" json:"create_time"`
	UpdateTime time.Time      `xorm:"update_time updated" json:"-"`
}

func (Feedback) TableName() string {
	return "feedback"
}

type FeedbackStatus int

const (
	FeedbackStatusClosed FeedbackStatus = iota
	FeedbackStatusOpen
	FeedbackStatusFinished
)
