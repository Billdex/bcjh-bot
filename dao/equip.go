package dao

import (
	"bcjh-bot/model/database"
	"fmt"
	"regexp"
	"strings"
)

const CacheKeyEquipList = "equip_list"

// ClearEquipsCache 清除厨具数据缓存
func ClearEquipsCache() {
	Cache.Delete(CacheKeyEquipList)
}

// FindAllEquips 查询全部厨具信息
func FindAllEquips() ([]database.Equip, error) {
	equips := make([]database.Equip, 0)
	err := SimpleFindDataWithCache(CacheKeyEquipList, &equips, func(dest interface{}) error {
		results := make([]database.Equip, 0)
		err := DB.OrderBy("equip_id").Find(&results)
		if err != nil {
			return err
		}
		// 载入技能数据
		mSkills, err := GetSkillsMap()
		if err != nil {
			return fmt.Errorf("载入技能数据出错 %v", err)
		}
		for i := range results {
			for _, id := range results[i].Skills {
				results[i].SkillDescs = append(results[i].SkillDescs, mSkills[id].Description)
			}
		}

		*dest.(*[]database.Equip) = results
		return nil
	})
	return equips, err
}

// SearchEquipsWithName 根据名称查询厨具数据
func SearchEquipsWithName(name string) ([]database.Equip, error) {
	pattern := strings.ReplaceAll(name, "%", ".*")
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("描述格式有误 %v", err)
	}
	equips, err := FindAllEquips()
	if err != nil {
		return nil, fmt.Errorf("查询厨具数据失败 %v", err)
	}
	result := make([]database.Equip, 0)
	for _, equip := range equips {
		if re.MatchString(equip.Name) {
			result = append(result, equip)
		}
	}
	return result, nil
}
