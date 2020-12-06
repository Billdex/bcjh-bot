package gamedata

// 调料
type Condiment struct {
	CondimentId int `json:"condimentId"`
	Name string `json:"name"`
	Rarity int `json:"rarity"`
	Skill []int `json:"skill"`
	Origin string `json:"origin"`
}