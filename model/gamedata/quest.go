package gamedata

// 任务
type QuestData struct {
	QuestId     int     `json:"questId"`
	QuestIdDisp float32 `json:"questIdDisp"`
	Type        string  `json:"type"`
	Goal        string  `json:"goal"`
	Rewards     []struct {
		Name     string `json:"name"`     // 奖励名称
		Quantity string `json:"quantity"` // 奖励数量
	} `json:"rewards"`
}

// type QuestRewards struct {
// 	Name     string `json:"name"`     // 奖励名称
// 	Quantity string `json:"quantity"` // 奖励数量
// }
