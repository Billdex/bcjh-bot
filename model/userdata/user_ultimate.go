package userdata

// UserUltimateData 用户修炼数据
type UserUltimateData struct {
	DecoBuff string `json:"decoBuff"`
	Stirfry  string `json:"Stirfry"`
	Boil     string `json:"Boil"`
	Knife    string `json:"Knife"`
	Fry      string `json:"Fry"`
	Bake     string `json:"Bake"`
	Steam    string `json:"Steam"`
	Male     string `json:"Male"`
	Female   string `json:"Female"`
	All      string `json:"All"`

	// buff
	Partial ultimateChef `json:"Partial"`
	Self    ultimateChef `json:"Self"`

	MaxLimit1  string `json:"MaxLimit_1"`
	MaxLimit2  string `json:"MaxLimit_2"`
	MaxLimit3  string `json:"MaxLimit_3"`
	MaxLimit4  string `json:"MaxLimit_4"`
	MaxLimit5  string `json:"MaxLimit_5"`
	PriceBuff1 string `json:"PriceBuff_1"`
	PriceBuff2 string `json:"PriceBuff_2"`
	PriceBuff3 string `json:"PriceBuff_3"`
	PriceBuff4 string `json:"PriceBuff_4"`
	PriceBuff5 string `json:"PriceBuff_5"`
}

type ultimateChef struct {
	Id  []string `json:"id"`
	Row []struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		SubName string `json:"subName"`
		Effect  []struct {
			Type      string `json:"type"`
			Value     int    `json:"value"`
			Condition string `json:"condition"`
			Cal       string `json:"cal"`
		} `json:"effect"`
	} `json:"row"`
}
