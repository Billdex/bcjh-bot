package database

type Admin struct {
	Id int `xorm:"id autoincr"`
	QQ int `xorm:"qq"`
}

func (Admin) TableName() string {
	return "admin"
}
