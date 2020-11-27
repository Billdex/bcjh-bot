package gamedata

type GameData struct {
	Chefs       []ChefData     `json:"chefs"`       // 厨师
	Equips      []EquipData    `json:"equips"`      // 厨具
	Recipes     []RecipeData   `json:"recipes"`     // 菜谱
	Combos      []ComboData    `json:"combos"`      // 合成菜
	Guests      []GuestData    `json:"guests"`      // 贵客
	Materials   []MaterialData `json:"materials"`   // 材料
	Skills      []SkillData    `json:"skills"`      // 技能
	Decorations []Decoration   `json:"decorations"` // 装修家具
}
