package dao

import (
	"bcjh-bot/model/database"
)

// FindAllChefs 查询全部厨师信息
func FindAllChefs() ([]database.Chef, error) {
	chefs := make([]database.Chef, 0)
	err := DB.Find(&chefs)
	return chefs, err
}
