package database

type Guest struct {
	GuestId   string `xorm:"guest_id"`
	GuestName string `xorm:"guest_name"`
}

func (Guest) TableName() string {
	return "guest"
}
