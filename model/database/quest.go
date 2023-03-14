package database

// 任务
type Quest struct {
	QuestId     int              `xorm:"quest_id"`      // 任务 ID
	QuestIdDisp string           `xorm:"quest_id_disp"` // 系列任务编号
	Type        string           `xorm:"type"`          // 系列类型
	Goal        string           `xorm:"goal"`          // 任务目标
	Rewards     []QuestRewards   `xorm:"rewards"`       // 任务奖励
	Conditions  []QuestCondition `json:"conditions"`    // 任务条件
}

type QuestRewards struct {
	Name     string `json:"name"`     // 奖励名称
	Quantity string `json:"quantity"` // 奖励数量
}

type QuestCondition struct {
	RecipeId     int    `json:"recipeId"`
	Rank         int    `json:"rank"`
	Num          int    `json:"num"`
	GoldEff      bool   `json:"goldEff"`
	MaterialId   int    `json:"materialId"`
	Guest        string `json:"guest"`
	AnyGuest     bool   `json:"anyGuest"`
	Skill        string `json:"skill"`
	MaterialEff  bool   `json:"materialEff"`
	NewGuest     bool   `json:"newGuest"`
	Rarity       int    `json:"rarity"`
	Price        int    `json:"price"`
	Category     string `json:"category"`
	Condiment    string `json:"condiment"`
	CondimentEff bool   `json:"condimentEff"`
}

func (Quest) TableName() string {
	return "quest"
}
