package dao

import (
	"bcjh-bot/model/database"
	"fmt"
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
