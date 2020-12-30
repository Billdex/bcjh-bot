package database

type Laboratory struct {
	Target        string   `xorm:"target_name"`
	TargetType    string   `xorm:"target_type"`
	Rarity        int      `xorm:"rarity"`
	Skill         string   `xorm:"skill"`
	Antique       string   `xorm:"antique"`
	AntiqueNumber int      `xorm:"antique_number"`
	Equips        []string `xorm:"equips"`
	Recipes       []string `xorm:"recipes"`
}

func (Laboratory) TableName() string {
	return "laboratory"
}
