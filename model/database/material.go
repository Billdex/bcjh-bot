package database

type Material struct {
	MaterialId int    `xorm:"material_id"` // 食材ID
	Name       string `xorm:"name"`        // 食材名
	Origin     string `xorm:"origin"`      // 食材来源
}

func (Material) TableName() string {
	return "material"
}
