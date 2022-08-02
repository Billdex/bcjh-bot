package database

type Admin struct {
	Id int   `xorm:"id pk autoincr"`
	QQ int64 `xorm:"qq bigint unique"`
}

func (Admin) TableName() string {
	return "admin"
}
