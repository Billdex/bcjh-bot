package database

type GuestGift struct {
	GuestId   string `xorm:"guest_id comment('贵客图鉴ID')"`
	GuestName string `xorm:"guest_name comment('贵客名称')"`
	Antique   string `xorm:"antique comment('礼物名')"`
	Recipe    string `xorm:"recipe comment('菜谱名')"`
}

func (GuestGift) TableName() string {
	return "guest_gift"
}
