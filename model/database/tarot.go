package database

type Tarot struct {
	Id          int    `xorm:"id autoincr pk"`
	Score       int    `xorm:"score"`
	Description string `xorm:"description"`
}

func (Tarot) TableName() string {
	return "tarot"
}
