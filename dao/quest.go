package dao

import "bcjh-bot/model/database"

func GetQuestById(id int) (database.Quest, error) {
	var quest database.Quest
	_, err := DB.Where("quest_id = ?", id).Get(&quest)
	return quest, err
}

func GetQuestsByIds(ids []int) ([]database.Quest, error) {
	quests := make([]database.Quest, 0, len(ids))
	err := DB.In("quest_id", ids).Find(&quests)
	return quests, err
}
