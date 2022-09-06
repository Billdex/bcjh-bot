package database

// RecipeMaterial 菜谱食材数据
type RecipeMaterial struct {
	RecipeGalleryId string `xorm:"recipe_id"`   // 菜谱图鉴ID
	MaterialId      int    `xorm:"material_id"` // 食材ID
	Quantity        int    `xorm:"quantity"`    // 数量
	Efficiency      int    `xorm:"efficiency"`  // 食材效率

	MaterialName string `xorm:"-"` // 食材名称，查询的时候用
}

func (RecipeMaterial) TableName() string {
	return "recipe_material"
}
