package gamedata

type Decoration struct {
	Id int `json:"id"`
	//Icon     int     `json:"icon"`
	Name     string  `json:"name"`
	TipMin   int     `json:"tipMin"`
	TipMax   int     `json:"tipMax"`
	TipTime  int     `json:"tipTime"`
	Gold     float32 `json:"gold"`
	Position string  `json:"position"`
	Suit     string  `json:"suit"`
	SuitGold float32 `json:"suitGold"`
	Origin   string  `json:"origin"`
}
