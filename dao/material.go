package dao

import (
	"bcjh-bot/model/database"
	"fmt"
)

// GetMaterialById 根据食材id查询食材的名称与来源
func GetMaterialById(id int) (database.Material, error) {
	var material database.Material
	has, err := DB.Where("material_id = ?", id).Get(&material)
	if err != nil {
		return material, err
	}
	if !has {
		return material, fmt.Errorf("菜谱%d数据缺失", id)
	}
	return material, nil
}

// FindRecipeMaterialByRecipeGalleryId 根据菜谱的图鉴ID查询对应的食材数据
func FindRecipeMaterialByRecipeGalleryId(id string, withName bool) ([]database.RecipeMaterial, error) {
	var recipeMaterials []database.RecipeMaterial
	err := DB.Where("recipe_id = ?", id).Find(&recipeMaterials)
	if err != nil {
		return nil, err
	}
	if withName {
		for i := range recipeMaterials {
			material, err := GetMaterialById(recipeMaterials[i].MaterialId)
			if err != nil {
				return nil, err
			}
			recipeMaterials[i].MaterialName = material.Name
		}
	}
	return recipeMaterials, err
}
