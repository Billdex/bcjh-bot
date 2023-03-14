package dao

import (
	"bcjh-bot/model/database"
	"sort"
)

const CacheKeyQuestList = "quest_list"

var mainQuests []database.Quest

// ClearQuestsCache 清除任务数据缓存
func ClearQuestsCache() {
	Cache.Delete(CacheKeyQuestList)
	mainQuests = []database.Quest{}
}

// FindAllQuests 查询全部任务信息
func FindAllQuests() ([]database.Quest, error) {
	var quests []database.Quest
	err := SimpleFindDataWithCache(CacheKeyQuestList, &quests, func(dest interface{}) error {
		return DB.OrderBy("quest_id").Find(dest)
	})
	return quests, err
}

// FindAllMainQuests 查询全部主线任务
func FindAllMainQuests() ([]database.Quest, error) {
	if len(mainQuests) != 0 {
		return mainQuests, nil
	}
	quests, err := FindAllQuests()
	if err != nil {
		return nil, err
	}
	results := make([]database.Quest, 0)
	for i := range quests {
		if quests[i].Type == "主线任务" {
			results = append(results, quests[i])
		}
	}
	// 任务列表顺序排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].QuestId < results[j].QuestId
	})
	mainQuests = results
	return mainQuests, nil
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

// FindQuestsWithIds 根据 id 列表查询任务列表
func FindQuestsWithIds(ids []int) ([]database.Quest, error) {
	quests := make([]database.Quest, 0, len(ids))
	questMap, err := GetQuestsMap()
	if err != nil {
		return quests, err
	}
	for _, qid := range ids {
		if quest, ok := questMap[qid]; ok {
			quests = append(quests, quest)
		}
	}
	return quests, nil
}

// GetMaxMainQuestId 查询最大主线任务id
func GetMaxMainQuestId() (int, error) {
	if len(mainQuests) != 0 {
		return mainQuests[len(mainQuests)-1].QuestId, nil
	}
	quests, err := FindAllMainQuests()
	if err != nil {
		return 0, err
	}
	return quests[len(mainQuests)-1].QuestId, nil
}

// FindMainQuestsWithLimit 查询从某个任务 id 起始的任务数据列表，最长不得超过 5 条
func FindMainQuestsWithLimit(startId int, limit int) ([]database.Quest, error) {
	if startId < 0 {
		startId = 0
	}
	if limit > 5 {
		limit = 5
	}
	if limit < 1 {
		limit = 1
	}
	quests, err := FindAllMainQuests()
	if err != nil {
		return nil, err
	}
	results := make([]database.Quest, 0, limit)
	for i := range quests {
		if quests[i].QuestId >= startId && quests[i].QuestId < startId+limit {
			results = append(results, quests[i])
		}
	}
	return results, nil
}
