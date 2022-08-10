package database

import "time"

type Feedback struct {
	Id         int64     `xorm:"id pk autoincr"`
	Sender     int64     `xorm:"sender BIGINT"`
	Nickname   string    `xorm:"nickname"`
	Message    string    `xorm:"message"`
	Status     int       `xorm:"status"`
	CreateTime time.Time `xorm:"create_time created"`
	UpdateTime time.Time `xorm:"update_time updated"`
}

func (Feedback) TableName() string {
	return "feedback"
}
