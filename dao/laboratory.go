package dao

import (
	"bcjh-bot/model/database"
)

const CacheKeyLaboratoryList = "laboratory_list"

// FindAllLaboratory 查询所有实验室菜谱数据
func FindAllLaboratory() ([]database.Laboratory, error) {
	var laboratories []database.Laboratory
	err := SimpleFindDataWithCache(CacheKeyLaboratoryList, &laboratories, func(dest interface{}) error {
		return DB.Find(dest)
	})
	return laboratories, err
}
