package database

import (
	"image"
	"regexp"
	"strings"
)

type Equip struct {
	EquipId   int    `xorm:"equip_id"`   // 厨具ID
	Name      string `xorm:"name"`       // 厨具名称
	GalleryId string `xorm:"gallery_id"` // 图鉴ID
	Origin    string `xorm:"origin"`     // 来源
	Rarity    int    `xorm:"rarity"`     // 稀有度
	Skills    []int  `xorm:"skills"`     // 技能效果组

	SkillDescs []string `xorm:"-"` // 技能效果描述组
}

func (Equip) TableName() string {
	return "equip"
}

// FormatRarity 格式化稀有度输出
func (equip Equip) FormatRarity() string {
	return strings.Repeat("🔥", equip.Rarity)
}

// HasSkill 判断厨具是否具有某个技能效果
func (equip Equip) HasSkill(skill string) bool {
	pattern := strings.ReplaceAll(skill, "%", ".*")
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	for i := range equip.SkillDescs {
		if re.MatchString(equip.SkillDescs[i]) {
			return true
		}
	}
	return false
}

type EquipData struct {
	Equip
	Avatar image.Image
	Skills []Skill
}
