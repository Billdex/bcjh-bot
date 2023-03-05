package userdata

type RepCol struct {
	Id            bool `json:"id"`
	Img           bool `json:"img"`
	Rarity        bool `json:"rarity"`
	Skills        bool `json:"skills"`
	SkillsSim     bool `json:"skills_sim"`
	Condiment     bool `json:"condiment"`
	Materials     bool `json:"materials"`
	Price         bool `json:"price"`
	ExPrice       bool `json:"exPrice"`
	Time          bool `json:"time"`
	Limit         bool `json:"limit"`
	TotalPrice    bool `json:"total_price"`
	TotalTimeShow bool `json:"total_time_show"`
	GoldEff       bool `json:"gold_eff"`
	MaterialEff   bool `json:"material_eff"`
	CondiEff      bool `json:"condi_eff"`
	Origin        bool `json:"origin"`
	Unlock        bool `json:"unlock"`
	Combo         bool `json:"combo"`
	Guests        bool `json:"guests"`
	DegreeGuests  bool `json:"degree_guests"`
	Gift          bool `json:"gift"`
	Got           bool `json:"got"`
}

type CalRepCol struct {
	Id            bool `json:"id"`
	Rarity        bool `json:"rarity"`
	SkillsSim     bool `json:"skills_sim"`
	Skills        bool `json:"skills"`
	Materials     bool `json:"materials"`
	Origin        bool `json:"origin"`
	Limit         bool `json:"limit"`
	Price         bool `json:"price"`
	BuffRule      bool `json:"buff_rule"`
	PriceRule     bool `json:"price_rule"`
	PriceTotal    bool `json:"price_total"`
	TotalTimeShow bool `json:"total_time_show"`
	GoldEff       bool `json:"gold_eff"`
}

type ChefCol struct {
	Id            bool `json:"id"`
	Img           bool `json:"img"`
	Rarity        bool `json:"rarity"`
	Skills        bool `json:"skills"`
	Skill         bool `json:"skill"`
	Gather        bool `json:"gather"`
	Condiment     bool `json:"condiment"`
	Sex           bool `json:"sex"`
	Origin        bool `json:"origin"`
	UltimateGoal  bool `json:"ultimateGoal"`
	UltimateSkill bool `json:"ultimateSkill"`
	Got           bool `json:"got"`
}

type EquipCol struct {
	Id     bool `json:"id"`
	Img    bool `json:"img"`
	Rarity bool `json:"rarity"`
	Skill  bool `json:"skill"`
	Origin bool `json:"origin"`
}

type CondimentCol struct {
	Id     bool `json:"id"`
	Img    bool `json:"img"`
	Rarity bool `json:"rarity"`
	Skill  bool `json:"skill"`
	Origin bool `json:"origin"`
}

type DecorationCol struct {
	Checkbox bool `json:"checkbox"`
	Id       bool `json:"id"`
	Img      bool `json:"img"`
	Gold     bool `json:"gold"`
	TipMin   bool `json:"tipMin"`
	TipMax   bool `json:"tipMax"`
	TipTime  bool `json:"tipTime"`
	EffMin   bool `json:"effMin"`
	EffMax   bool `json:"effMax"`
	EffAvg   bool `json:"effAvg"`
	Position bool `json:"position"`
	Suit     bool `json:"suit"`
	SuitGold bool `json:"suitGold"`
	Origin   bool `json:"origin"`
}

type MapCol struct {
	Field1 bool `json:"0"`
	Field2 bool `json:"1"`
	Field3 bool `json:"2"`
	Field4 bool `json:"3"`
	Field5 bool `json:"4"`
}
