package dao

import "bcjh-bot/model/database"

const CacheKeyQuestList = "quest_list"

// ClearQuestsCache 清除任务数据缓存
func ClearQuestsCache() {
	Cache.Delete(CacheKeyQuestList)
}

// FindAllQuests 查询全部任务信息
func FindAllQuests() ([]database.Quest, error) {
	var quests []database.Quest
	err := SimpleFindDataWithCache(CacheKeyQuestList, &quests, func(dest interface{}) error {
		return DB.OrderBy("quest_id").Find(dest)
	})
	return quests, err
}

// GetQuestsMap 获取 map 格式的任务数据，key 为任务 id
func GetQuestsMap() (map[int]database.Quest, error) {
	quests, err := FindAllQuests()
	if err != nil {
		return nil, err
	}
	mResult := make(map[int]database.Quest)
	for _, quest := range quests {
		mResult[quest.QuestId] = quest
	}
	return mResult, nil
}
