package database

type Recipe struct {
	RecipeId           int      `xorm:"pk recipe_id comment('菜谱ID')"`
	Name               string   `xorm:"name comment('菜名')"`
	GalleryId          string   `xorm:"index gallery_id comment('图鉴ID')"`
	Rarity             int      `xorm:"rarity comment('稀有度')"`
	Origin             string   `xorm:"origin comment('来源')"`
	Stirfry            int      `xorm:"stirfry comment('炒技法')"`
	Bake               int      `xorm:"bake comment('烤技法')"`
	Boil               int      `xorm:"boil comment('煮技法')"`
	Steam              int      `xorm:"steam comment('蒸技法')"`
	Fry                int      `xorm:"fry comment('炸技法')"`
	Cut                int      `xorm:"cut comment('切技法knife')"`
	Price              int      `xorm:"price comment('价格')"`
	ExPrice            int      `xorm:"exPrice comment('熟练加价')"`
	GoldEfficiency     int      `xorm:"gold_efficiency comment('金币效率')"`
	MaterialEfficiency int      `xorm:"material_efficiency comment('耗材效率')"`
	Gift               string   `xorm:"gift comment('第一次做到神级送的符文')"`
	Guests             []string `xorm:"guests comment('升阶贵客')"`
	Time               int      `xorm:"'time' comment('每份时间(秒)')"`
	Limit              int      `xorm:"limit comment('每组数量')"`
	TotalTime          int      `xorm:"total_time comment('每组时间(秒)')"`
	Unlock             string   `xorm:"unlock comment('可解锁')"`
	Combo              string   `xorm:"combo comment('可合成')"`
}

func (Recipe) TableName() string {
	return "recipe"
}
