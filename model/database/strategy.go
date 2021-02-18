package database

import "time"

type Strategy struct {
	Id         int       `xorm:"id autoincr pk comment('ID')"`
	Keyword    string    `xorm:"keyword"`
	Value      string    `xorm:"value"`
	CreateTime time.Time `xorm:"'create_time' created"`
	UpdateTime time.Time `xorm:"'update_time' updated"`
}

func (Strategy) TableName() string {
	return "strategy"
}
