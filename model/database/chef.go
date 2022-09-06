package database

import "image"

// Chef 厨师数据对应的数据库模型
type Chef struct {
	ChefId        int    `xorm:"pk chef_id"`       // 厨师 id
	Name          string `xorm:"name"`             // 厨师名字
	Gender        int    `xorm:"gender"`           // 性别
	Rarity        int    `xorm:"index rarity"`     // 稀有度
	Origin        string `xorm:"origin"`           // 来源
	GalleryId     string `xorm:"index gallery_id"` // 图鉴id
	Stirfry       int    `xorm:"stirfry"`          // 炒技法
	Bake          int    `xorm:"bake"`             // 烤技法
	Boil          int    `xorm:"boil"`             // 煮技法
	Steam         int    `xorm:"steam"`            // 蒸技法
	Fry           int    `xorm:"fry"`              // 炸技法
	Cut           int    `xorm:"cut"`              // 切技法knife
	Meat          int    `xorm:"meat"`             // 肉类采集
	Flour         int    `xorm:"flour"`            // 面类采集
	Fish          int    `xorm:"fish"`             // 水产采集
	Vegetable     int    `xorm:"vegetable"`        // 蔬菜采集
	Sweet         int    `xorm:"sweet"`            // 甜
	Sour          int    `xorm:"sour"`             // 酸
	Spicy         int    `xorm:"spicy"`            // 辣
	Salty         int    `xorm:"salty"`            // 咸
	Bitter        int    `xorm:"bitter"`           // 苦
	Tasty         int    `xorm:"tasty"`            // 鲜
	SkillId       int    `xorm:"skill_id"`         // 技能id
	UltimateGoals []int  `xorm:"ultimate_goals"`   // 修炼任务id数组
	UltimateSkill int    `xorm:"ultimate_skill"`   // 修炼效果id
}

func (Chef) TableName() string {
	return "chef"
}

// ChefData 厨师数据的另一种模型，专用于厨师图鉴的绘制
type ChefData struct {
	Chef
	Avatar        image.Image
	Skill         string
	UltimateGoals []string
	UltimateSkill string
}

// GetCondimentData 获取厨师的调味数据，返回值分别为数值和类型名称
func (c ChefData) GetCondimentData() (int, string) {
	if c.Sweet > 0 {
		return c.Sweet, "Sweet"
	} else if c.Sour > 0 {
		return c.Sour, "Sour"
	} else if c.Spicy > 0 {
		return c.Spicy, "Spicy"
	} else if c.Salty > 0 {
		return c.Salty, "Salty"
	} else if c.Bitter > 0 {
		return c.Bitter, "Bitter"
	} else if c.Tasty > 0 {
		return c.Tasty, "Tasty"
	} else {
		// 没有调料数值则视为甜味
		return 0, "Sweet"
	}
}

func (c ChefData) GetCondimentValue() int {
	value, _ := c.GetCondimentData()
	return value
}

func (c ChefData) GetCondimentType() string {
	_, CondimentType := c.GetCondimentData()
	return CondimentType
}

func (c ChefData) GetUltimateSkill() string {
	if c.UltimateSkill == "" {
		return "暂无"
	}
	return c.UltimateSkill
}

func (c ChefData) GetUltimateGoals() []string {
	goals := make([]string, 3)
	for i := 0; i < 3; i++ {
		if i >= len(c.UltimateGoals) || c.UltimateGoals[i] == "" {
			goals[i] = "暂无"
		} else {
			goals[i] = c.UltimateGoals[i]
		}
	}
	return goals
}
