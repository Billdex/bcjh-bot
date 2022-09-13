package dao

import "bcjh-bot/model/database"

func FindAllSkills() ([]database.Skill, error) {
	skills := make([]database.Skill, 0)
	err := DB.Find(&skills)
	return skills, err
}

func FindSkillsByIds(ids []int) ([]database.Skill, error) {
	skills := make([]database.Skill, 0)
	err := DB.In("skill_id", ids).Find(&skills)
	return skills, err
}

func GetSkillById(id int) (database.Skill, error) {
	skill := database.Skill{}
	_, err := DB.Where("skill_id = ?", id).Get(&skill)
	return skill, err
}
