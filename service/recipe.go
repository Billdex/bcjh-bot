package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/model/database"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"fmt"
)

func RecipeQuery(c *onebot.Context, args []string) {
	logger.Info("菜谱查询, 参数:", args)

	if len(args) == 0 {
		err := bot.SendMessage(c,
			fmt.Sprintf("指令示例:\n"+
				"%s菜谱 荷包蛋", util.PrefixCharacter))
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
		return
	}
	if args[0] == "%" {
		err := bot.SendMessage(c, "参数有误!")
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
		return
	}

	recipes := make([]database.Recipe, 0)
	err := database.DB.Where("gallery_id = ?", args[0]).Asc("gallery_id").Find(&recipes)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_ = bot.SendMessage(c, "查询数据失败!")
		return
	}
	if len(recipes) == 0 {
		err = database.DB.Where("name like ?", "%"+args[0]+"%").Asc("gallery_id").Find(&recipes)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			_ = bot.SendMessage(c, "查询数据失败!")
			return
		}
	}

	var msg string
	if len(recipes) == 0 {
		msg = "未查询到数据!"
	} else if len(recipes) == 1 {
		recipe := recipes[0]
		rarity := ""
		for i := 0; i < recipe.Rarity; i++ {
			rarity += "🔥"
		}
		goldEfficiency := (int)(float64(recipe.Price) * (3600.0 / float64(recipe.Time)))
		time := util.FormatSecondToString(recipe.Time)
		allTime := util.FormatSecondToString(recipe.Time * recipe.Limit)

		materials := ""
		materialQuantities := 0
		for _, m := range recipe.Materials {
			materialQuantities += m.Quantity
			material := new(database.Material)
			has, err := database.DB.Where("material_id = ?", m.MaterialId).Get(material)
			if err != nil {
				logger.Error("查询数据库出错!", err)
				_ = bot.SendMessage(c, "查询数据失败!")
				return
			}
			if !has {
				_ = bot.SendMessage(c, "查询数据失败!")
				return
			}
			materials += fmt.Sprintf("%s*%d ", material.Name, m.Quantity)
		}
		materialEfficiency := (int)(float64(materialQuantities) * (3600.0 / float64(recipe.Time)))

		guests := ""
		if len(recipe.Guests) == 3 {
			if recipe.Guests[0] != "" {
				guests += fmt.Sprintf("优-%s, ", recipe.Guests[0])
			} else {
				guests += fmt.Sprintf("优-未知,")
			}
			if recipe.Guests[1] != "" {
				guests += fmt.Sprintf("特-%s, ", recipe.Guests[1])
			} else {
				guests += fmt.Sprintf("特-未知,")
			}
			if recipe.Guests[2] != "" {
				guests += fmt.Sprintf("神-%s", recipe.Guests[2])
			} else {
				guests += fmt.Sprintf("神-未知")
			}
		} else {
			logger.Errorf("%s贵客数据有误!", recipe.Name)
			_ = bot.SendMessage(c, "查询数据失败!")
			return
		}

		msg += fmt.Sprintf("%s %s %s\n", recipe.GalleryId, recipe.Name, rarity)
		msg += fmt.Sprintf("售价: %d(%d)  效率: %d/h\n", recipe.Price, recipe.Price+recipe.ExPrice, goldEfficiency)
		msg += fmt.Sprintf("来源: %s\n", recipe.Origin)
		msg += fmt.Sprintf("单份耗时: %s\n", time)
		msg += fmt.Sprintf("每组份数: %d\n", recipe.Limit)
		msg += fmt.Sprintf("一组耗时: %s\n", allTime)
		msg += fmt.Sprintf("炒:%d 烤:%d 煮:%d\n", recipe.Stirfry, recipe.Bake, recipe.Boil)
		msg += fmt.Sprintf("蒸:%d 炸:%d 切:%d\n", recipe.Steam, recipe.Fry, recipe.Cut)
		msg += fmt.Sprintf("材料: %s\n", materials)
		msg += fmt.Sprintf("耗材效率: %d/h\n", materialEfficiency)
		msg += fmt.Sprintf("神级符文: %s\n", recipe.Gift)
		msg += fmt.Sprintf("可解锁: %s\n", recipe.Unlock)
		msg += fmt.Sprintf("可合成: %s\n", recipe.Combo)
		msg += fmt.Sprintf("贵客: %s", guests)
	} else {
		msg = "查询到以下菜谱:\n"
		for p, recipe := range recipes {
			msg += fmt.Sprintf("%s %s", recipe.GalleryId, recipe.Name)
			if p != len(recipes)-1 {
				msg += "\n"
				if p == util.MaxSearchList {
					msg += "......"
					break
				}
			}
		}
	}

	err = bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}
