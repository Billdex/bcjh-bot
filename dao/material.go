package dao

import (
	"bcjh-bot/model/database"
)

const (
	CacheKeyMaterialList       = "material_list"
	CacheKeyRecipeMaterialList = "recipe_material_list"
)

// ClearMaterialsCache 清除食材数据缓存
func ClearMaterialsCache() {
	Cache.Delete(CacheKeyMaterialList)
	Cache.Delete(CacheKeyRecipeMaterialList)
}

// FindAllMaterials 查询全部食材信息
func FindAllMaterials() ([]database.Material, error) {
	var materials []database.Material
	err := SimpleFindDataWithCache(CacheKeyMaterialList, &materials, func(dest interface{}) error {
		return DB.OrderBy("material_id").Find(dest)
	})
	return materials, err
}

// FindAllRecipeMaterials 获取所有菜谱食材数据
func FindAllRecipeMaterials() ([]database.RecipeMaterial, error) {
	var recipeMaterials []database.RecipeMaterial
	err := SimpleFindDataWithCache(CacheKeyRecipeMaterialList, &recipeMaterials, func(dest interface{}) error {
		return DB.OrderBy("recipe_id").Find(dest)
	})
	return recipeMaterials, err
}

// GetRecipeMaterialsMap 获取 map 格式的菜谱食材关联数据，key 为菜谱图鉴 id
func GetRecipeMaterialsMap() (map[string][]database.RecipeMaterial, error) {
	recipeMaterials, err := FindAllRecipeMaterials()
	if err != nil {
		return nil, err
	}
	materials, err := FindAllMaterials()
	if err != nil {
		return nil, err
	}
	mMaterialNames := make(map[int]string)
	for _, material := range materials {
		mMaterialNames[material.MaterialId] = material.Name
	}
	mResult := make(map[string][]database.RecipeMaterial)
	for _, recipeMaterial := range recipeMaterials {
		id := recipeMaterial.RecipeGalleryId
		recipeMaterial.MaterialName = mMaterialNames[recipeMaterial.MaterialId]
		mResult[id] = append(mResult[id], recipeMaterial)
	}
	return mResult, nil
}
