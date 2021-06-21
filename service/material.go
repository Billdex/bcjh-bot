package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/model/database"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"fmt"
	"strconv"
	"strings"
)

func MaterialQuery(c *onebot.Context, args []string) {
	logger.Info("食材及效率查询:", args)

	if len(args) == 0 {
		err := bot.SendMessage(c, materialHelp())
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
		return
	}
	page := 1
	if len(args) > 1 {
		if strings.HasPrefix(args[1], "p") || strings.HasPrefix(args[1], "P") {
			num, err := strconv.Atoi(args[1][1:])
			if err != nil {
				logger.Error("字符串转int失败!", err)
			} else {
				if num < 1 {
					num = 1
				}
				page = num
			}
		}
	}

	materials := make([]database.Material, 0)
	err := database.DB.Where("name like ?", "%"+args[0]+"%").Find(&materials)
	if err != nil {
		logger.Error("数据库查询出错!")
		_ = bot.SendMessage(c, util.SystemErrorNote)
		return
	}
	if len(materials) == 0 {
		err := bot.SendMessage(c, fmt.Sprintf("没有找到叫%s的食材", args[0]))
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
		return
	}
	// 匹配大于1个时，如果有完全匹配的则直接使用该食材
	if len(materials) > 1 {
		match := false
		var msg string
		msg += "找到以下多个食材"
		for _, material := range materials {
			if material.Name == args[0] {
				match = true
				break
			}
			msg += fmt.Sprintf("\n%s %s", material.Name, material.Origin)
		}
		if !match {
			_ = bot.SendMessage(c, msg)
			return
		}
	}

	recipeMaterials := make([]database.RecipeMaterial, 0)
	err = database.DB.Where("material_id = ?", materials[0].MaterialId).Find(&recipeMaterials)
	if err != nil {
		logger.Error("数据库查询出错!")
		_ = bot.SendMessage(c, util.SystemErrorNote)
		return
	}
	if len(recipeMaterials) == 0 {
		err := bot.SendMessage(c, fmt.Sprintf("没有使用%s的菜谱哦~", args[0]))
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
		return
	}

	recipeMaterialMap := make(map[string]int)
	recipeGalleryIds := make([]string, 0)
	for _, recipeMaterial := range recipeMaterials {
		recipeGalleryIds = append(recipeGalleryIds, recipeMaterial.RecipeGalleryId)
		recipeMaterialMap[recipeMaterial.RecipeGalleryId] = recipeMaterial.Efficiency
	}

	// 根据查出的信息查询菜谱信息
	recipes := make([]database.Recipe, 0)
	err = database.DB.In("gallery_id", recipeGalleryIds).Find(&recipes)
	if err != nil {
		logger.Error("数据库查询出错!")
		_ = bot.SendMessage(c, util.SystemErrorNote)
		return
	}

	for i, _ := range recipes {
		recipes[i].MaterialEfficiency = recipeMaterialMap[recipes[i].GalleryId]
	}

	var note string
	recipes, note = orderRecipes(recipes, "耗材效率")
	if note != "" {
		_ = bot.SendMessage(c, note)
		return
	}

	// 处理消息
	var msg string
	listLength := util.MaxQueryListLength
	if c.MessageType == util.OneBotMessagePrivate {
		listLength = listLength * 2
	}
	maxPage := (len(recipes)-1)/listLength + 1
	if page > maxPage {
		page = maxPage
	}
	if len(recipes) > listLength {
		msg += fmt.Sprintf("以下菜谱使用了%s: (%d/%d)\n", materials[0].Name, page, maxPage)
	} else {
		msg += fmt.Sprintf("以下菜谱使用了%s:\n", materials[0].Name)
	}
	for i := (page - 1) * listLength; i < page*listLength && i < len(recipes); i++ {
		orderInfo := getRecipeInfoWithOrder(recipes[i], "耗材效率")
		msg += fmt.Sprintf("%s %s %s", recipes[i].GalleryId, recipes[i].Name, orderInfo)
		if i < page*listLength-1 && i < len(recipes)-1 {
			msg += "\n"
		}
	}
	if page < maxPage {
		msg += "\n......"
	}
	logger.Info("发送食材效率查询结果:", msg)
	err = bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}
