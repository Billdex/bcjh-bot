package database

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
