package database

type Equip struct {
	EquipId   int    `xorm:"equip_id"`   // 厨具ID
	Name      string `xorm:"name"`       // 厨具名称
	GalleryId string `xorm:"gallery_id"` // 图鉴ID
	Origin    string `xorm:"origin"`     // 来源
	Rarity    int    `xorm:"rarity"`     // 稀有度
	Skills    []int  `xorm:"skills"`     // 技能效果组
}

func (Equip) TableName() string {
	return "equip"
}
