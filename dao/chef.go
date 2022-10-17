package dao

import (
	"bcjh-bot/model/database"
	"fmt"
)

const CacheKeyChefList = "chef_list"

// ClearChefsCache 清除厨师数据缓存
func ClearChefsCache() {
	Cache.Delete(CacheKeyChefList)
}

// FindAllChefs 查询全部厨师信息
func FindAllChefs() ([]database.Chef, error) {
	chefs := make([]database.Chef, 0)
	err := SimpleFindDataWithCache(CacheKeyChefList, &chefs, func(dest interface{}) error {
		results := make([]database.Chef, 0)
		err := DB.OrderBy("chef_id").Find(&results)
		if err != nil {
			return err
		}
		// 载入技能数据
		mSkills, err := GetSkillsMap()
		if err != nil {
			return fmt.Errorf("载入技能数据出错 %v", err)
		}
		for i := range results {
			results[i].SkillDesc = mSkills[results[i].SkillId].Description
			results[i].UltimateSkillDesc = mSkills[results[i].UltimateSkill].Description
		}

		*dest.(*[]database.Chef) = results
		return nil
	})
	return chefs, err
}
