package database

// 任务
type Quest struct {
	QuestId     int            `xorm:"questId comment('任务 ID')"`
	QuestIdDisp string         `xorm:"questIdDisp comment('系列任务编号')"`
	Type        string         `xorm:"type comment('系列类型')"`
	Goal        string         `xorm:"goal comment('任务目标')"`
	Rewards     []QuestRewards `xorm:"rewards comment('任务奖励')"`
}

type QuestRewards struct {
	Name     string `json:"name"`     // 奖励名称
	Quantity string `json:"quantity"` // 奖励数量
}

func (Quest) TableName() string {
	return "quest"
}
