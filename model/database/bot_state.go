package database

type BotState struct {
	BotId   int64 `xorm:"bot_id pk"`
	GroupId int64 `xorm:"group_id pk"`
	State   bool  `xorm:"state"`
}

func (BotState) TableName() string {
	return "bot_state"
}
