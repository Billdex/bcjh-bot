package userdata

import (
	"bcjh-bot/util"
)

// UserUltimateData 用户修炼数据
type UserUltimateData struct {
	DecoBuff util.IntString `json:"decoBuff"`
	Stirfry  util.IntString `json:"Stirfry"`
	Boil     util.IntString `json:"Boil"`
	Knife    util.IntString `json:"Knife"`
	Fry      util.IntString `json:"Fry"`
	Bake     util.IntString `json:"Bake"`
	Steam    util.IntString `json:"Steam"`
	Male     util.IntString `json:"Male"`
	Female   util.IntString `json:"Female"`
	All      util.IntString `json:"All"`

	// buff
	Partial ultimateChef `json:"Partial"`
	Self    ultimateChef `json:"Self"`

	MaxLimit1  util.IntString `json:"MaxLimit_1"`
	MaxLimit2  util.IntString `json:"MaxLimit_2"`
	MaxLimit3  util.IntString `json:"MaxLimit_3"`
	MaxLimit4  util.IntString `json:"MaxLimit_4"`
	MaxLimit5  util.IntString `json:"MaxLimit_5"`
	PriceBuff1 util.IntString `json:"PriceBuff_1"`
	PriceBuff2 util.IntString `json:"PriceBuff_2"`
	PriceBuff3 util.IntString `json:"PriceBuff_3"`
	PriceBuff4 util.IntString `json:"PriceBuff_4"`
	PriceBuff5 util.IntString `json:"PriceBuff_5"`
}

type ultimateChef struct {
	Id  []string `json:"id"`
	Row []struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		SubName string `json:"subName"`
		Effect  []struct {
			Type      string  `json:"type"`
			Value     float64 `json:"value"`
			Condition string  `json:"condition"`
			Cal       string  `json:"cal"`
		} `json:"effect"`
	} `json:"row"`
}
