package database

import "time"

type Exchange struct {
	Id         int       `xorm:"id autoincr pk"`
	Content    string    `xorm:"content longtext"`
	CreateTime time.Time `xorm:"'create_time' created"`
	UpdateTime time.Time `xorm:"'update_time' updated"`
}

func (Exchange) TableName() string {
	return "exchange"
}
