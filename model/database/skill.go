package database

// Skill 技能数据
type Skill struct {
	SkillId     int           `xorm:"pk skill_id"` // 技能ID
	Description string        `xorm:"description"` // 技能描述
	Effects     []SkillEffect `xorm:"effect"`      // 效果
}

func (Skill) TableName() string {
	return "skill"
}

// SkillEffect 技能效果详情
type SkillEffect struct {
	Calculation string  `json:"calculation"`
	Type        string  `json:"type"`
	Condition   string  `json:"condition"`
	Tag         int     `json:"tag"` // 对厨师生效的性别 1:男 2:女
	Value       float64 `json:"value"`
}
