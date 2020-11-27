package gamedata

// 任务
type QuestData struct {
	QuestId     int    `json:"questId"`
	QuestIdDisp     int    `json:"questIdDisp"`
	Type string `json:"type"`
	Goal string `json:"goal"`
	Rewards []struct {
		Name string
		Quantity string
	} `json:"rewards"`
	Conditions []struct{
		asdf struct{

		}
	}



	Description string `json:"desc"`
	Effects     []struct {
		Calculation string  `json:"cal"`
		Type        string  `json:"type"`
		Condition   string  `json:"condition"`
		Tag         int     `json:"tag"`
		Value       float64 `json:"value"`
	} `json:"effect"`
}
