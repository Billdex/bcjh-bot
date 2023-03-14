package gamedata

// 任务
type QuestData struct {
	QuestId     int     `json:"questId"`
	QuestIdDisp float32 `json:"questIdDisp"`
	Type        string  `json:"type"`
	Goal        string  `json:"goal"`
	PreId       string  `json:"preId"`
	Rewards     []struct {
		Name     string `json:"name"`     // 奖励名称
		Quantity string `json:"quantity"` // 奖励数量
	} `json:"rewards"`
	Conditions []struct {
		RecipeId     int    `json:"recipeId,omitempty"`
		Rank         int    `json:"rank,omitempty"`
		Num          int    `json:"num,omitempty"`
		GoldEff      bool   `json:"goldEff,omitempty"`
		MaterialId   int    `json:"materialId,omitempty"`
		Guest        string `json:"guest,omitempty"`
		AnyGuest     bool   `json:"anyGuest,omitempty"`
		Skill        string `json:"skill,omitempty"`
		MaterialEff  bool   `json:"materialEff,omitempty"`
		NewGuest     bool   `json:"newGuest,omitempty"`
		Rarity       int    `json:"rarity,omitempty"`
		Price        int    `json:"price,omitempty"`
		Category     string `json:"category,omitempty"`
		Condiment    string `json:"condiment,omitempty"`
		CondimentEff bool   `json:"condimentEff,omitempty"`
	} `json:"conditions"`
}

// type QuestRewards struct {
// 	Name     string `json:"name"`     // 奖励名称
// 	Quantity string `json:"quantity"` // 奖励数量
// }
