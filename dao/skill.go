package dao

import "bcjh-bot/model/database"

func GetSkillById(id int) (database.Skill, error) {
	skill := database.Skill{}
	_, err := DB.Where("skill_id = ?", id).Get(&skill)
	return skill, err
}
