package database

// 任务
type Quest struct {
	QuestId     int            `xorm:"quest_id"`      // 任务 ID
	QuestIdDisp string         `xorm:"quest_id_disp"` // 系列任务编号
	Type        string         `xorm:"type"`          // 系列类型
	Goal        string         `xorm:"goal"`          // 任务目标
	Rewards     []QuestRewards `xorm:"rewards"`       // 任务奖励
}

type QuestRewards struct {
	Name     string `json:"name"`     // 奖励名称
	Quantity string `json:"quantity"` // 奖励数量
}

func (Quest) TableName() string {
	return "quest"
}
