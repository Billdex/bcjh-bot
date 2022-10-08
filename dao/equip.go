package dao

import (
	"bcjh-bot/model/database"
)

const CacheKeyDataEquips = "data_equips"

// FindAllEquips 查询全部厨具信息
func FindAllEquips() ([]database.Equip, error) {
	equips := make([]database.Equip, 0)
	err := SimpleFindDataWithCache(CacheKeyDataEquips, &equips)
	return equips, err
}
