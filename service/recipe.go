package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/model/database"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"errors"
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

	var err error
	var msg string
	if len(args) > 1 {
		var order string
		//var orderField, orderType string
		if len(args) >= 3 {
			//orderField, orderType = getRecipeOrderType(args[2])
			order = args[2]
		} else {
			//orderField, orderType = getRecipeOrderType("")
			order = ""
		}
		switch args[0] {
		case "食材":
			if len(args) < 2 {
				_ = bot.SendMessage(c, "参数有误")
				return
			}
			msg, err = getRecipeMsgWithMaterial(args[1], order)
		default:
			msg = "参数有误!"
		}
		if err != nil {
			logger.Error("查询数据出错!", err)
			_ = bot.SendMessage(c, "查询数据失败!")
			return
		}
	} else {
		msg, err = getRecipeMsgWithName(args[0])
		if err != nil {
			logger.Error("查询数据出错!", err)
			_ = bot.SendMessage(c, "查询数据失败!")
			return
		}
	}

	err = bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}

func getRecipeOrderType(order string) (string, string) {
	switch order {
	case "单时间":
		return "time", "ASC"
	case "总时间":
		return "time*limit", "ASC"
	case "售价":
		return "price", "DESC"
	case "赚钱效率":
		return "price*3600/time", "DESC"
	default:
		return "gallery_id", "ASC"
	}
}

func getRecipeMsgWithName(arg string) (string, error) {
	recipes := make([]database.Recipe, 0)
	err := database.DB.Where("gallery_id = ?", arg).Asc("gallery_id").Find(&recipes)
	if err != nil {
		return "", err
	}
	if len(recipes) == 0 {
		err = database.DB.Where("name like ?", "%"+arg+"%").Asc("gallery_id").Find(&recipes)
		if err != nil {
			return "", err
		}
	}
	var msg string
	if len(recipes) == 0 {
		return "哎呀，好像找不到呢!", nil
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
				return "", err
			}
			if !has {
				return "", err
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
			return "", errors.New(fmt.Sprintf("%s升阶贵客数据有误!", recipe.Name))
		}

		msg += fmt.Sprintf("%s %s %s\n", recipe.GalleryId, recipe.Name, rarity)
		msg += fmt.Sprintf("售价: %d(%d)\n", recipe.Price, recipe.Price+recipe.ExPrice)
		msg += fmt.Sprintf("赚钱效率: %d/h\n", goldEfficiency)
		msg += fmt.Sprintf("来源: %s\n", recipe.Origin)
		msg += fmt.Sprintf("单份耗时: %s\n", time)
		msg += fmt.Sprintf("每组份数: %d\n", recipe.Limit)
		msg += fmt.Sprintf("一组耗时: %s\n", allTime)
		msg += fmt.Sprintf("炒:%d 烤:%d 煮:%d\n", recipe.Stirfry, recipe.Bake, recipe.Boil)
		msg += fmt.Sprintf("蒸:%d 炸:%d 切:%d\n", recipe.Steam, recipe.Fry, recipe.Cut)
		msg += fmt.Sprintf("食材: %s\n", materials)
		msg += fmt.Sprintf("耗材效率: %d/h\n", materialEfficiency)
		msg += fmt.Sprintf("神级符文: %s\n", recipe.Gift)
		msg += fmt.Sprintf("可解锁: %s\n", recipe.Unlock)
		msg += fmt.Sprintf("可合成: %s\n", recipe.Combo)
		msg += fmt.Sprintf("贵客-符文: %s\n", recipe.GuestAntiques)
		msg += fmt.Sprintf("升阶贵客: %s", guests)
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

	return msg, nil
}

func getRecipeMsgWithMaterial(arg string, order string) (string, error) {
	recipes := make([]database.Recipe, 0)
	material := new(database.Material)
	has, err := database.DB.Where("name = ?", arg).Get(material)
	if err != nil {
		return "", err
	}
	if !has {
		return "食材参数有误!", nil
	}

	queryArg := fmt.Sprintf("%%\"MaterialId\":%d%%", material.MaterialId)
	orderField, orderType := getRecipeOrderType(order)
	switch orderType {
	case "ASC":
		err = database.DB.Where("materials like ?", queryArg).Asc(orderField).Find(&recipes)
	case "DESC":
		err = database.DB.Where("materials like ?", queryArg).Desc(orderField).Find(&recipes)
	default:
		err = database.DB.Where("materials like ?", queryArg).Asc(orderField).Find(&recipes)
	}
	if err != nil {
		return "", err
	}

	var msg string
	if len(recipes) == 0 {
		return "哎呀，好像找不到呢!", nil
	} else {
		msg = "查询到以下菜谱:\n"
		for p, recipe := range recipes {
			var thirdInfo string
			switch order {
			case "单时间":
				thirdInfo = util.FormatSecondToString(recipe.Time)
			case "总时间":
				thirdInfo = util.FormatSecondToString(recipe.Time * recipe.Limit)
			case "售价":
				thirdInfo = fmt.Sprintf("$%d", recipe.Price)
			case "赚钱效率":
				thirdInfo = fmt.Sprintf("$%d/h", recipe.Price*3600/recipe.Time)
			default:
				thirdInfo = ""
			}
			msg += fmt.Sprintf("%s %s %s", recipe.GalleryId, recipe.Name, thirdInfo)
			if p != len(recipes)-1 {
				msg += "\n"
				if p == util.MaxSearchList {
					msg += "......"
					break
				}
			}
		}
	}

	return msg, nil
}

//func conditionQueryRecipe(condition string, arg string, order string) []database.Recipe{
//	query := "condition = ?"
//
//
//}
