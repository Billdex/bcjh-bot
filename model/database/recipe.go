package database

type RecipeMaterial struct {
	MaterialId int `xorm:"material_id"`
	Quantity   int `xorm:"quantity"`
}

type Recipe struct {
	RecipeId  int              `xorm:"recipe_id comment('菜谱ID')"`
	Name      string           `xorm:"name comment('菜名')"`
	GalleryId string           `xorm:"gallery_id comment('图鉴ID')"`
	Rarity    int              `xorm:"rarity comment('稀有度')"`
	Origin    string           `xorm:"origin comment('来源')"`
	Stirfry   int              `xorm:"stirfry comment('炒技法')"`
	Bake      int              `xorm:"bake comment('烤技法')"`
	Boil      int              `xorm:"boil comment('煮技法')"`
	Steam     int              `xorm:"steam comment('蒸技法')"`
	Fry       int              `xorm:"fry comment('炸技法')"`
	Cut       int              `xorm:"knife comment('切技法')"`
	Price     int              `xorm:"price comment('价格')"`
	ExPrice   int              `xorm:"exPrice comment('熟练加价')"`
	Gift      string           `xorm:"gift comment('第一次做到神级送的符文')"`
	Guests    []string         `xorm:"guests comment('升阶贵客')"`
	Limit     int              `xorm:"limit comment('每组数量')"`
	Time      int              `xorm:"'time' comment('每份时间(秒)')"`
	Unlock    string           `xorm:"unlock comment('可解锁')"`
	Combo     string           `xorm:"combo comment('可合成')"`
	Materials []RecipeMaterial `xorm:"materials comment(材料)"`
}

func (Recipe) TableName() string {
	return "recipe"
}
