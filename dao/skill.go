package dao

import "bcjh-bot/model/database"

const CacheKeySkillList = "skill_list"

// ClearSkillsCache 清除技能数据缓存
func ClearSkillsCache() {
	Cache.Delete(CacheKeySkillList)
}

// FindAllSkills 查询全部技能数据
func FindAllSkills() ([]database.Skill, error) {
	skills := make([]database.Skill, 0)
	err := SimpleFindDataWithCache(CacheKeySkillList, &skills, func(dest interface{}) error {
		return DB.OrderBy("skill_id").Find(&skills)
	})
	return skills, err
}

// GetSkillsMap 获取 map 格式的技能数据，key 为技能 id
func GetSkillsMap() (map[int]database.Skill, error) {
	skills, err := FindAllSkills()
	if err != nil {
		return nil, err
	}
	mResult := make(map[int]database.Skill)
	for _, skill := range skills {
		mResult[skill.SkillId] = skill
	}
	return mResult, nil
}
