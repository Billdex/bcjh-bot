package database

type Tarot struct {
	Id          int    `xorm:"id autoincr pk"`
	Score       int    `xorm:"score"`
	Description string `xorm:"description"`
}

func (Tarot) TableName() string {
	return "tarot"
}

func (t Tarot) Level() string {
	switch {
	case t.Score == 0:
		return "不知道吉不吉"
	case 0 < t.Score && t.Score < 15:
		return "小小吉"
	case 15 <= t.Score && t.Score < 40:
		return "小吉"
	case 40 <= t.Score && t.Score < 60:
		return "中吉"
	case 60 <= t.Score && t.Score < 85:
		return "大吉"
	case 85 <= t.Score && t.Score < 100:
		return "大大吉"
	case t.Score == 100:
		return "超吉"
	default:
		return "?"
	}
}
