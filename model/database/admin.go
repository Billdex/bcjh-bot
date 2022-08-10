package database

type Admin struct {
	Id int   `xorm:"id autoincr"`
	QQ int64 `xorm:"qq"`
}

func (Admin) TableName() string {
	return "admin"
}
