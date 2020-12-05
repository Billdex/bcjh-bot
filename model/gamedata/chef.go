package gamedata

type ChefData struct {
	ChefId        int    `json:"chefId"`
	Name          string `json:"name"`
	Tags          []int  `json:"tags"`
	Rarity        int    `json:"rarity"`
	Origin        string `json:"origin"`
	GalleryId     string `json:"galleryId"`
	Stirfry       int    `json:"stirfry"`
	Bake          int    `json:"bake"`
	Boil          int    `json:"boil"`
	Steam         int    `json:"steam"`
	Fry           int    `json:"fry"`
	Cut           int    `json:"knife"`
	Meat          int    `json:"meat"`
	Flour         int    `json:"creation"`
	Fish          int    `json:"fish"`
	Vegetable     int    `json:"veg"`
	Sweet         int    `json:"sweet"`
	Sour          int    `json:"sour"`
	Spicy         int    `json:"spicy"`
	Salty         int    `json:"salty"`
	Bitter        int    `json:"bitter"`
	Tasty         int    `json:"tasty"`
	SkillId       int    `json:"skill"`
	UltimateGoals []int  `json:"ultimateGoal"`
	UltimateSkill int    `json:"ultimateSkill"`
}
