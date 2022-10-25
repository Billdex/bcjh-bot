package dao

import (
	"bcjh-bot/model/database"
	"fmt"
	"regexp"
	"strings"
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

// SearchChefsWithName 根据名字搜索厨师
func SearchChefsWithName(name string) ([]database.Chef, error) {
	pattern := strings.ReplaceAll(name, "%", ".*")
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("描述格式有误 %v", err)
	}
	chefs, err := FindAllChefs()
	if err != nil {
		return nil, fmt.Errorf("查询厨师数据失败 %v", err)
	}
	result := make([]database.Chef, 0)
	for _, chef := range chefs {
		if re.MatchString(chef.Name) {
			result = append(result, chef)
		}
	}
	return result, nil
}
