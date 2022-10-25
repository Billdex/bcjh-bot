package dao

import (
	"bcjh-bot/model/database"
	"fmt"
	"regexp"
	"strings"
)

const CacheKeyRecipeList = "recipe_list"

// ClearRecipesCache 清除菜谱数据缓存
func ClearRecipesCache() {
	Cache.Delete(CacheKeyRecipeList)
}

// FindAllRecipes 查询全部菜谱信息
func FindAllRecipes() ([]database.Recipe, error) {
	recipes := make([]database.Recipe, 0)
	err := SimpleFindDataWithCache(CacheKeyRecipeList, &recipes, func(dest interface{}) error {
		var results []database.Recipe
		err := DB.OrderBy("recipe_id").Find(&results)
		if err != nil {
			return err
		}
		// 载入食材数据
		mMaterials, err := GetRecipeMaterialsMap()
		if err != nil {
			return fmt.Errorf("载入菜谱食材数据出错 %v", err)
		}
		// 载入贵客礼物
		mGuestGifts, err := GetRecipeGuestGiftsMap()
		if err != nil {
			return fmt.Errorf("载入菜谱贵客礼物数据出错 %v", err)
		}

		for i := range results {
			results[i].Materials = mMaterials[results[i].GalleryId]
			results[i].GuestGifts = mGuestGifts[results[i].Name]
		}

		*dest.(*[]database.Recipe) = results
		return nil
	})
	return recipes, err
}

// SearchRecipesWithName 根据菜谱名查询菜谱
func SearchRecipesWithName(name string) ([]database.Recipe, error) {
	pattern := strings.ReplaceAll(name, "%", ".*")
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("描述格式有误 %v", err)
	}
	recipes, err := FindAllRecipes()
	if err != nil {
		return nil, fmt.Errorf("查询菜谱数据失败 %v", err)
	}
	result := make([]database.Recipe, 0)
	for _, recipe := range recipes {
		if re.MatchString(recipe.Name) {
			result = append(result, recipe)
		}
	}
	return result, nil
}
