package database

import (
	"time"
)

type BlackList struct {
	QQ          int64     `xorm:"qq pk"`
	GroupId     int64     `xorm:"group_id pk"`
	EndTime     int64     `xorm:"end_time"`
	EndDatetime time.Time `xorm:"end_datetime"`
}

func (BlackList) TableName() string {
	return "black_list"
}
