package database

// 装修
type Decoration struct {
	Name     string  `xorm:"name"`      // 家具名称
	Position string  `xorm:"position"`  // 家具位置
	Suit     string  `xorm:"suit"`      // 套装名称
	Origin   string  `xorm:"origin"`    // 家具来源
	Id       int     `xorm:"id"`        // 家具ID
	TipMin   int     `xorm:"tipMin"`    // 单次最小玉璧产值
	TipMax   int     `xorm:"tipMax"`    // 单次最大玉璧产值
	TipTime  int     `xorm:"tipTime"`   // 玉璧产出冷却时间（单位：s）
	Gold     float32 `xorm:"gold"`      // 营业加成（以1为单位）
	SuitGold float32 `xorm:"suit_gold"` // 整套家具营业加成（以1为单位）
}

func (Decoration) TableName() string {
	return "decoration"
}
