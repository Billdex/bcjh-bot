package dao

import (
	"bcjh-bot/model/database"
	"fmt"
	"regexp"
	"strings"
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

// SearchMaterialsWithName 根据名称筛选食材列表
func SearchMaterialsWithName(name string) ([]database.Material, error) {
	pattern := strings.ReplaceAll(name, "%", ".*")
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("食材描述格式有误 %v", err)
	}
	materials, err := FindAllMaterials()
	if err != nil {
		return nil, fmt.Errorf("查询食材数据失败 %v", err)
	}
	result := make([]database.Material, 0)
	for _, material := range materials {
		if re.MatchString(material.Name) {
			result = append(result, material)
		}
	}
	return result, nil
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
	mMaterials := make(map[int]database.Material)
	for _, material := range materials {
		mMaterials[material.MaterialId] = material
	}
	mResult := make(map[string][]database.RecipeMaterial)
	for _, recipeMaterial := range recipeMaterials {
		id := recipeMaterial.RecipeGalleryId
		recipeMaterial.Material = mMaterials[recipeMaterial.MaterialId]
		mResult[id] = append(mResult[id], recipeMaterial)
	}
	return mResult, nil
}
