package database

type WhiteList struct {
	Id      int    `xorm:"id pk autoincr"`
	Plugin  string `xorm:"plugin"`
	GroupId int    `xorm:"group_id"`
}

func (WhiteList) TableName() string {
	return "white_list"
}
