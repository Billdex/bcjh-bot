package database

import "time"

type UserData struct {
	// 自增 ID
	Id int `xorm:"id autoincr pk"`

	// 用户绑定的 qq 号
	QQ int64 `xorm:"qq bigint unique"`

	// 用户上传白菜菊花数据时使用的名称
	User string `xorm:"user"`

	// 用户上传的白菜菊花个人数据 ID
	BcjhID int `xorm:"bcjh_id"`

	// 白菜菊花用户个人数据内容
	Data string `xorm:"data text"`

	// 用户上传白菜菊花数据的时间
	CreateTime time.Time `xorm:"create_time"`
}

func (UserData) TableName() string {
	return "user_data"
}
