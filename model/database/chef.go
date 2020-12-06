package database

type Chef struct {
	ChefId        int    `xorm:"pk chef_id comment('厨师id')"`
	Name          string `xorm:"name comment('厨师名字')"`
	Gender        int    `xorm:"gender comment('性别')"`
	Rarity        int    `xorm:"index rarity comment('稀有度')"`
	Origin        string `xorm:"comment('来源')"`
	GalleryId     string `xorm:"index gallery_id comment('图鉴id')"`
	Stirfry       int    `xorm:"stirfry comment('炒技法'')"`
	Bake          int    `xorm:"bake comment('烤技法')"`
	Boil          int    `xorm:"boil comment('煮技法')"`
	Steam         int    `xorm:"steam comment('蒸技法')"`
	Fry           int    `xorm:"fry comment('炸技法')"`
	Cut           int    `xorm:"cut comment('切技法knife')"`
	Meat          int    `xorm:"meat comment('肉类采集')"`
	Flour         int    `xorm:"flour comment('面类采集')"`
	Fish          int    `xorm:"fish comment('水产采集')"`
	Vegetable     int    `xorm:"vegetable comment('蔬菜采集')"`
	Sweet         int    `xorm:"sweet comment('甜')"`
	Sour          int    `xorm:"sour comment('酸')"`
	Spicy         int    `xorm:"spicy comment('辣')"`
	Salty         int    `xorm:"salty comment('咸')"`
	Bitter        int    `xorm:"bitter comment(苦)"`
	Tasty         int    `xorm:"tasty comment(鲜)"`
	SkillId       int    `xorm:"skill_id comment('技能id')"`
	UltimateGoals []int  `xorm:"ultimate_goals comment('修炼任务id数组')"`
	UltimateSkill int    `xorm:"ultimate_skill comment('修炼效果id')"`
}

func (Chef) TableName() string {
	return "chef"
}
