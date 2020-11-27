package database

type SkillEffect struct {
	Calculation string  `json:"calculation"`
	Type        string  `json:"type"`
	Condition   string  `json:"condition"`
	Tag         int     `json:"tag"` // 对厨师生效的性别 1:男 2:女
	Value       float64 `json:"value"`
}

type Skill struct {
	SkillId     int           `xorm:"pk skill_id comment('技能ID')"`
	Description string        `xorm:"description comment('技能描述')"`
	Effects     []SkillEffect `xorm:"effect comment('效果')"`
}

func (Skill) TableName() string {
	return "skill"
}
