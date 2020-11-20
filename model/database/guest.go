package database

type GuestGift struct {
	Antique string `json:"antique"`
	Recipe  string `json:"recipe"`
}

type Guest struct {
	GuestId   int         `xorm:"guest_id comment('贵客ID')"`
	Name      string      `xorm:"name comment('贵客名称')"`
	GalleryId string      `xorm:"gallery_id comment('图鉴ID')"`
	Gifts     []GuestGift `xorm:"gifts comment('礼物符文')"`
}

func (Guest) TableName() string {
	return "guest"
}
