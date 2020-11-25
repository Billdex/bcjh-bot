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
		if len(args) >= 3 {
			order = args[2]
		} else {
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

	logger.Debug("发送一条消息:", msg)
	err = bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}

func getRecipeOrderString(order string) (string, bool) {
	switch order {
	case "单时间":
		return "`time` ASC", true
	case "总时间":
		return "`total_time` ASC", true
	case "单价":
		return "`price` DESC", true
	case "金币效率":
		return "`gold_efficiency` DESC", true
	case "耗材效率":
		return "`material_efficiency` DESC", true
	case "":
		return "`gallery_id` ASC", true
	default:
		return "", false
	}
}

func getRecipeOrderInfo(recipe database.Recipe, order string) string {
	switch order {
	case "单时间":
		return util.FormatSecondToString(recipe.Time)
	case "总时间":
		return util.FormatSecondToString(recipe.Time * recipe.Limit)
	case "单价":
		return fmt.Sprintf("💰%d", recipe.Price)
	case "金币效率":
		return fmt.Sprintf("💰%d/h", recipe.GoldEfficiency)
	case "耗材效率":
		return fmt.Sprintf("🥗%d/h", recipe.MaterialEfficiency)
	case "食材效率":
		return fmt.Sprintf("🥗%d/h", recipe.MaterialEfficiency)
	case "":
		return ""
	default:
		return ""
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
		logger.Info("未查询到菜谱")
		return "哎呀，好像找不到呢!", nil
	} else if len(recipes) == 1 {
		logger.Info("查询到一个菜谱")
		recipe := recipes[0]
		rarity := ""
		for i := 0; i < recipe.Rarity; i++ {
			rarity += "🔥"
		}
		time := util.FormatSecondToString(recipe.Time)
		allTime := util.FormatSecondToString(recipe.Time * recipe.Limit)

		materials := ""
		recipeMaterials := make([]database.RecipeMaterial, 0)
		err = database.DB.Where("recipe_id = ?", recipe.GalleryId).Find(&recipeMaterials)
		if err != nil {
			return "", err
		}
		for _, recipeMaterial := range recipeMaterials {
			material := new(database.Material)
			has, err := database.DB.Where("material_id = ?", recipeMaterial.MaterialId).Get(material)
			if err != nil {
				return "", err
			}
			if !has {
				return "", nil
			}
			materials += fmt.Sprintf("%s*%d ", material.Name, recipeMaterial.Quantity)
		}

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

		msg += fmt.Sprintf("[%s]%s %s\n", recipe.GalleryId, recipe.Name, rarity)
		msg += fmt.Sprintf("💰: %d(%d) --- %d/h\n", recipe.Price, recipe.Price+recipe.ExPrice, recipe.GoldEfficiency)
		msg += fmt.Sprintf("来源: %s\n", recipe.Origin)
		msg += fmt.Sprintf("单时间: %s\n", time)
		msg += fmt.Sprintf("总时间: %s (%d份)\n", allTime, recipe.Limit)
		msg += fmt.Sprintf("炒:%d 烤:%d 煮:%d\n", recipe.Stirfry, recipe.Bake, recipe.Boil)
		msg += fmt.Sprintf("蒸:%d 炸:%d 切:%d\n", recipe.Steam, recipe.Fry, recipe.Cut)
		msg += fmt.Sprintf("食材: %s\n", materials)
		msg += fmt.Sprintf("耗材效率: %d/h\n", recipe.MaterialEfficiency)
		msg += fmt.Sprintf("可解锁: %s\n", recipe.Unlock)
		msg += fmt.Sprintf("可合成: %s\n", recipe.Combo)
		msg += fmt.Sprintf("神级符文: %s\n", recipe.Gift)
		msg += fmt.Sprintf("贵客-符文: %s\n", recipe.GuestAntiques)
		msg += fmt.Sprintf("升阶贵客: %s", guests)
	} else {
		logger.Info("查询到多个菜谱")
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
	recipeMaterials := make([]database.RecipeMaterial, 0)
	material := new(database.Material)
	has, err := database.DB.Where("name = ?", arg).Get(material)
	if err != nil {
		return "", err
	}
	if !has {
		return "食材参数有误!", nil
	}

	if order == "食材效率" {
		err = database.DB.Where("material_id = ?", material.MaterialId).Desc("efficiency").Find(&recipeMaterials)
		if err != nil {
			return "", err
		}
		for _, recipeMaterial := range recipeMaterials {
			var recipe database.Recipe
			has, err := database.DB.Where("gallery_id = ?", recipeMaterial.RecipeGalleryId).Get(&recipe)
			if err != nil {
				return "", err
			}
			if !has {
				return "", errors.New(fmt.Sprintf("未查到图鉴Id %s 的菜谱", recipeMaterial.RecipeGalleryId))
			}
			recipe.MaterialEfficiency = recipeMaterial.Efficiency
			recipes = append(recipes, recipe)
		}
	} else {
		err = database.DB.Where("material_id = ?", material.MaterialId).Find(&recipeMaterials)
		if err != nil {
			return "", err
		}
		recipeIds := make([]string, 0)
		for _, recipeMaterial := range recipeMaterials {
			recipeIds = append(recipeIds, recipeMaterial.RecipeGalleryId)
		}
		orderStr, success := getRecipeOrderString(order)
		if !(success) {
			return "参数有误!", nil
		}
		err = database.DB.In("gallery_id", recipeIds).OrderBy(orderStr).Find(&recipes)
		if err != nil {
			return "", err
		}
	}

	msg := "查询到以下菜谱:\n"
	for p, recipe := range recipes {
		thirdInfo := getRecipeOrderInfo(recipe, order)
		msg += fmt.Sprintf("[%s]%s %s", recipe.GalleryId, recipe.Name, thirdInfo)
		if p != len(recipes)-1 {
			msg += "\n"
			if p == util.MaxSearchList-1 {
				msg += "......"
				break
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
