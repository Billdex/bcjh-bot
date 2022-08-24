package database

type Recipe struct {
	RecipeId           int      `xorm:"pk recipe_id"`        // 菜谱ID
	Name               string   `xorm:"name"`                // 菜名
	GalleryId          string   `xorm:"index gallery_id"`    // 图鉴ID
	Rarity             int      `xorm:"rarity"`              // 稀有度
	Origin             string   `xorm:"origin"`              // 来源
	Stirfry            int      `xorm:"stirfry"`             // 炒技法
	Bake               int      `xorm:"bake"`                // 烤技法
	Boil               int      `xorm:"boil"`                // 煮技法
	Steam              int      `xorm:"steam"`               // 蒸技法
	Fry                int      `xorm:"fry"`                 // 炸技法
	Cut                int      `xorm:"cut"`                 // 切技法knife
	Condiment          string   `xorm:"condiment"`           // 调料
	Price              int      `xorm:"price"`               // 价格
	ExPrice            int      `xorm:"exPrice"`             // 熟练加价
	GoldEfficiency     int      `xorm:"gold_efficiency"`     // 金币效率
	MaterialEfficiency int      `xorm:"material_efficiency"` // 耗材效率
	Gift               string   `xorm:"gift"`                // 第一次做到神级送的符文
	Guests             []string `xorm:"guests"`              // 升阶贵客
	Time               int      `xorm:"'time'"`              // 每份时间(秒)
	Limit              int      `xorm:"limit"`               // 每组数量
	TotalTime          int      `xorm:"total_time"`          // 每组时间(秒)
	Unlock             string   `xorm:"unlock"`              // 可解锁
	Combo              []string `xorm:"combo"`               // 可合成
}

func (Recipe) TableName() string {
	return "recipe"
}
