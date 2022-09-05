package dao

import (
	"bcjh-bot/model/database"
)

// FindAllEquips 查询全部厨具信息
func FindAllEquips() ([]database.Equip, error) {
	equips := make([]database.Equip, 0)
	err := DB.Find(&equips)
	return equips, err
}
