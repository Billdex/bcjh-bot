package dao

import (
	"bcjh-bot/model/database"
)

const CacheKeyTarotList = "tarot_list"

// FindAllTarots 查询全部签文信息
func FindAllTarots() ([]database.Tarot, error) {
	tarots := make([]database.Tarot, 0)
	err := SimpleFindDataWithCache(CacheKeyTarotList, &tarots, func(dest interface{}) error {
		return DB.OrderBy("id").Find(dest)
	})
	return tarots, err
}

// FindTarotsWithScore 根据分值查询签文列表
func FindTarotsWithScore(score int) ([]database.Tarot, error) {
	tarots, err := FindAllTarots()
	if err != nil {
		return nil, err
	}
	result := make([]database.Tarot, 0)
	for _, tarot := range tarots {
		if tarot.Score == score {
			result = append(result, tarot)
		}
	}
	return result, nil
}
