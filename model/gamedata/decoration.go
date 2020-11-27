package gamedata

// 家具套装
type Decoration struct {
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Suit     string  `json:"suit"`
	Origin   string  `json:"origin"`
	Id       int     `json:"id"`
	TipMin   int     `json:"tipMin"`
	TipMax   int     `json:"tipMax"`
	TipTime  int     `json:"tipTime"`
	Gold     float32 `json:"gold"`
	SuitGold float32 `json:"suitGold"`

	// Icon     int     `json:"icon"`

}
