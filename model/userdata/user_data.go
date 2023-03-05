package userdata

import (
	"encoding/json"
)

type UserData struct {
	RepCol        RepCol        `json:"repCol"`
	CalRepCol     CalRepCol     `json:"calRepCol"`
	ChefCol       ChefCol       `json:"chefCol"`
	EquipCol      EquipCol      `json:"equipCol"`
	CondimentCol  CondimentCol  `json:"condimentCol"`
	DecorationCol DecorationCol `json:"decorationCol"`
	MapCol        MapCol        `json:"mapCol"`
	// 用户修炼数据
	UserUltimate  UserUltimateData `json:"userUltimate"`
	UserNav       int              `json:"userNav"`
	ShowDetail    bool             `json:"showDetail"`
	DefaultEx     bool             `json:"defaultEx"`
	CalShowGot    bool             `json:"calShowGot"`
	HideSuspend   bool             `json:"hideSuspend"`
	HiddenMessage bool             `json:"hiddenMessage"`
	RepGot        json.RawMessage  `json:"repGot"`
	ChefGot       json.RawMessage  `json:"chefGot"`
	PlanList      []Plan           `json:"planList"`
	AllEx         bool             `json:"allEx"`
	CustomRules   struct{}         `json:"customRules"`
}

type Plan struct {
	Name string `json:"name"`
	Data struct {
		Chef struct {
			Field1 int `json:"1"`
			Field2 int `json:"2"`
		} `json:"Chef"`
		Equip struct {
			Field1 int `json:"1"`
			Field2 int `json:"2"`
		} `json:"Equip"`
		Condiment struct {
			Field1 int `json:"1"`
			Field2 int `json:"2"`
		} `json:"Condiment"`
		Rep struct {
			Field1 int `json:"1-1"`
			Field2 int `json:"2-1"`
			Field3 int `json:"2-2"`
		} `json:"Rep"`
		Cnt struct {
			Field1 int `json:"1-1"`
			Field2 int `json:"2-1"`
			Field3 int `json:"2-2"`
		} `json:"Cnt"`
		Ex struct {
			Field1 bool `json:"1-1"`
			Field2 bool `json:"2-1"`
			Field3 bool `json:"2-2"`
		} `json:"Ex"`
		Condi struct {
			Field1 bool `json:"1-1"`
			Field2 bool `json:"2-1"`
			Field3 bool `json:"2-2"`
		} `json:"Condi"`
	} `json:"data"`
}
