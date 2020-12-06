package gamedata

type SkillData struct {
	SkillId     int    `json:"skillId"`
	Description string `json:"desc"`
	Effects     []struct {
		Calculation string  `json:"cal"`
		Type        string  `json:"type"`
		Condition   string  `json:"condition"`
		Tag         int     `json:"tag"`
		Value       float64 `json:"value"`
	} `json:"effect"`
}
