package dao

import (
	"bcjh-bot/model/database"
)

const CacheKeyEquipList = "equip_list"

// FindAllEquips 查询全部厨具信息
func FindAllEquips() ([]database.Equip, error) {
	equips := make([]database.Equip, 0)
	err := SimpleFindDataWithCache(CacheKeyEquipList, &equips, func(dest interface{}) error {
		return DB.Find(dest)
	})
	return equips, err
}
