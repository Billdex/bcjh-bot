package database

type BlackList struct {
	QQ      int64 `xorm:"qq pk"`
	GroupId int64 `xorm:"group_id pk"`
	EndTime int64 `xorm:"end_time"`
}

func (BlackList) TableName() string {
	return "black_list"
}
