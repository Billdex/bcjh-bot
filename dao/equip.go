package dao

import (
	"bcjh-bot/model/database"
)

const CacheKeyEquipList = "equip_list"

// ClearEquipsCache 清除厨具数据缓存
func ClearEquipsCache() {
	Cache.Delete(CacheKeyEquipList)
}

// FindAllEquips 查询全部厨具信息
func FindAllEquips() ([]database.Equip, error) {
	equips := make([]database.Equip, 0)
	err := SimpleFindDataWithCache(CacheKeyEquipList, &equips, func(dest interface{}) error {
		return DB.OrderBy("equip_id").Find(dest)
	})
	return equips, err
}
