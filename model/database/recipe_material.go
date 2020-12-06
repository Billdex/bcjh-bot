package database

type RecipeMaterial struct {
	RecipeGalleryId string `xorm:"recipe_id comment('菜谱图鉴ID')"`
	MaterialId      int    `xorm:"material_id comment('食材ID')"`
	Quantity        int    `xorm:"quantity comment('数量')"`
	Efficiency      int    `xorm:"efficiency comment('食材效率')"`
}

func (RecipeMaterial) TableName() string {
	return "recipe_material"
}
