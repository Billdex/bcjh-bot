package database

type GuestGift struct {
	GuestId   string `xorm:"guest_id"`   // 贵客图鉴ID
	GuestName string `xorm:"guest_name"` // 贵客名称
	Antique   string `xorm:"antique"`    // 礼物名
	Recipe    string `xorm:"recipe"`     // 菜谱名
	TotalTime int    `xorm:"total_time"` // 总时间
}

func (GuestGift) TableName() string {
	return "guest_gift"
}
