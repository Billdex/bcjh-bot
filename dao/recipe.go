package dao

import (
	"bcjh-bot/model/database"
)

// FindAllRecipes 查询全部菜谱信息
func FindAllRecipes() ([]database.Recipe, error) {
	recipes := make([]database.Recipe, 0)
	err := DB.Find(&recipes)
	return recipes, err
}
