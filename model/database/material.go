package database

type Material struct {
	MaterialId int    `xorm:"material_id comment('食材ID')"`
	Name       string `xorm:"name comment('食材名')"`
	Origin     string `xorm:"origin comment('食材来源')"`
}

func (Material) TableName() string {
	return "material"
}
