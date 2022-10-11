package dao

import (
	"bcjh-bot/model/database"
)

const CacheKeyChefList = "chef_list"

// ClearChefsCache 清除厨师数据缓存
func ClearChefsCache() {
	Cache.Delete(CacheKeyChefList)
}

// FindAllChefs 查询全部厨师信息
func FindAllChefs() ([]database.Chef, error) {
	chefs := make([]database.Chef, 0)
	err := SimpleFindDataWithCache(CacheKeyChefList, &chefs, func(dest interface{}) error {
		return DB.OrderBy("chef_id").Find(dest)
	})
	return chefs, err
}
