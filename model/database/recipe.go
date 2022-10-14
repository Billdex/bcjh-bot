package database

import (
	"fmt"
	"image"
	"regexp"
	"strings"
)

type Recipe struct {
	RecipeId           int      `xorm:"pk recipe_id"`        // 菜谱ID
	Name               string   `xorm:"name"`                // 菜名
	GalleryId          string   `xorm:"index gallery_id"`    // 图鉴ID
	Rarity             int      `xorm:"rarity"`              // 稀有度
	Origin             string   `xorm:"origin"`              // 来源
	Stirfry            int      `xorm:"stirfry"`             // 炒技法
	Bake               int      `xorm:"bake"`                // 烤技法
	Boil               int      `xorm:"boil"`                // 煮技法
	Steam              int      `xorm:"steam"`               // 蒸技法
	Fry                int      `xorm:"fry"`                 // 炸技法
	Cut                int      `xorm:"cut"`                 // 切技法knife
	Condiment          string   `xorm:"condiment"`           // 调料
	Price              int      `xorm:"price"`               // 价格
	ExPrice            int      `xorm:"exPrice"`             // 熟练加价
	GoldEfficiency     int      `xorm:"gold_efficiency"`     // 金币效率
	MaterialEfficiency int      `xorm:"material_efficiency"` // 耗材效率
	Gift               string   `xorm:"gift"`                // 第一次做到神级送的符文
	Guests             []string `xorm:"guests"`              // 升阶贵客
	Time               int      `xorm:"'time'"`              // 每份时间(秒)
	Limit              int      `xorm:"limit"`               // 每组数量
	TotalTime          int      `xorm:"total_time"`          // 每组时间(秒)
	Unlock             string   `xorm:"unlock"`              // 可解锁
	Combo              []string `xorm:"combo"`               // 可合成

	Materials  []RecipeMaterial `xorm:"-"` // 所需食材数据
	GuestGifts []GuestGift      `xorm:"-"` // 贵客礼物数据
}

func (Recipe) TableName() string {
	return "recipe"
}

func (recipe Recipe) GetSkillValueMap() map[string]int {
	m := map[string]int{
		"stirfry": recipe.Stirfry,
		"bake":    recipe.Bake,
		"boil":    recipe.Boil,
		"steam":   recipe.Steam,
		"fry":     recipe.Fry,
		"cut":     recipe.Cut,
	}
	return m
}

// NeedSkill 判断菜谱是否需要某个技法
func (recipe Recipe) NeedSkill(skill string) (bool, error) {
	switch skill {
	case "炒":
		return recipe.Stirfry > 0, nil
	case "烤":
		return recipe.Bake > 0, nil
	case "煮":
		return recipe.Boil > 0, nil
	case "蒸":
		return recipe.Steam > 0, nil
	case "炸":
		return recipe.Fry > 0, nil
	case "切":
		return recipe.Cut > 0, nil
	default:
		return false, fmt.Errorf("%s是什么技法呀", skill)
	}
}

// UsedMaterial 判断菜谱是否使用了某个食材
func (recipe Recipe) UsedMaterial(material string) bool {
	for i := range recipe.Materials {
		if recipe.Materials[i].Material.Name == material {
			return true
		}
	}
	return false
}

// UsedMaterials 判断菜谱是否使用了某些食材
func (recipe Recipe) UsedMaterials(materials []string) bool {
	for _, material := range materials {
		if recipe.UsedMaterial(material) {
			return true
		}
	}
	return false
}

// HasMaterialFrom 判断菜谱是否有食材来自某个采集点
func (recipe Recipe) HasMaterialFrom(origin string) bool {
	for i := range recipe.Materials {
		if recipe.Materials[i].Material.Origin == origin {
			return true
		}
	}
	return false
}

// HasMaterialOrigins 判断菜谱是否有食材来自某些采集点
func (recipe Recipe) HasMaterialOrigins(origins []string) bool {
	for _, origin := range origins {
		if recipe.HasMaterialFrom(origin) {
			return true
		}
	}
	return false
}

// HasGuest 判断菜谱是否会来某个贵客
func (recipe Recipe) HasGuest(guest string) bool {
	pattern := strings.ReplaceAll(guest, "%", ".*")
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	for i := range recipe.GuestGifts {
		if re.MatchString(recipe.GuestGifts[i].GuestName) {
			return true
		}
	}
	return false
}

// HasAntique 判断菜谱是否会给某个菜谱
func (recipe Recipe) HasAntique(antique string) bool {
	pattern := strings.ReplaceAll(antique, "%", ".*")
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	for i := range recipe.GuestGifts {
		if re.MatchString(recipe.GuestGifts[i].Antique) {
			return true
		}
	}
	return false
}

// RecipeData 用于绘制厨师图片数据信息的模型
type RecipeData struct {
	Recipe
	Avatar image.Image
	Skills []RecipeSkillData
}

type RecipeSkillData struct {
	Type  string
	Value int
	Image image.Image
}
