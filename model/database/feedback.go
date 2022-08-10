package database

import "time"

type Feedback struct {
	Sender     int64     `xorm:"sender BIGINT"`
	Nickname   string    `xorm:"nickname"`
	Message    string    `xorm:"message"`
	Status     int       `xorm:"status"`
	CreateTime time.Time `xorm:"create_time created"`
}

func (Feedback) TableName() string {
	return "feedback"
}
