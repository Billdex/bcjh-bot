package database

import "time"

// Strategy 攻略数据
type Strategy struct {
	Id         int       `xorm:"id autoincr pk"`
	Keyword    string    `xorm:"keyword"`
	Value      string    `xorm:"value longtext"`
	CreateTime time.Time `xorm:"'create_time' created"`
	UpdateTime time.Time `xorm:"'update_time' updated"`
}

func (Strategy) TableName() string {
	return "strategy"
}
