package dao

import (
	"bcjh-bot/model/database"
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
		return DB.OrderBy("recipe_id").Find(dest)
	})
	return recipes, err
}
