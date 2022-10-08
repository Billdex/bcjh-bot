package dao

import (
	"bcjh-bot/model/database"
)

const CacheKeyDataChefs = "data_chefs"

// FindAllChefs 查询全部厨师信息
func FindAllChefs() ([]database.Chef, error) {
	chefs := make([]database.Chef, 0)
	err := SimpleFindDataWithCache(CacheKeyDataChefs, &chefs)
	return chefs, err
}
