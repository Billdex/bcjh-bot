package database

// 装修
type Decoration struct {
	Name     string  `xorm:"name comment('家具名称')"`
	Position string  `xorm:"position comment('家具位置')"`
	Suit     string  `xorm:"suit comment('套装名称')"`
	Origin   string  `xorm:"origin comment('家具来源')"`
	Id       int     `xorm:"id comment('家具ID')"`
	TipMin   int     `xorm:"tipMin comment('单次最小玉璧产值')"`
	TipMax   int     `xorm:"tipMax comment('单次最大玉璧产值')"`
	TipTime  int     `xorm:"tipTime comment('玉璧产出冷却时间（单位：s）')"`
	Gold     float32 `xorm:"gold comment('营业加成（以1为单位）')"`
	SuitGold float32 `xorm:"suitGold comment('整套家具营业加成（以1为单位）')"`

	// Icon     int     `xorm:"icon comment('家具图片ID')"`

}

func (Decoration) TableName() string {
	return "decoration"
}
