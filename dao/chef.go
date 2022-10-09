package dao

import (
	"bcjh-bot/model/database"
)

const CacheKeyChefList = "chef_list"

// FindAllChefs 查询全部厨师信息
func FindAllChefs() ([]database.Chef, error) {
	chefs := make([]database.Chef, 0)
	err := SimpleFindDataWithCache(CacheKeyChefList, &chefs, func(dest interface{}) error {
		return DB.Find(dest)
	})
	return chefs, err
}
