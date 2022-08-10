package database

type Equip struct {
	EquipId   int    `xorm:"equip_id comment('厨具ID')"`
	Name      string `xorm:"name comment('厨具名称')"`
	GalleryId string `xorm:"gallery_id comment('图鉴ID')"`
	Origin    string `xorm:"origin comment('来源')"`
	Rarity    int    `xorm:"rarity comment('稀有度')"`
	Skills    []int  `xorm:"skills comment('技能效果组')"`
}

func (Equip) TableName() string {
	return "equip"
}
