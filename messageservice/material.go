package messageservice

import (
	"bcjh-bot/config"
	"bcjh-bot/dao"
	"bcjh-bot/model/database"
	"bcjh-bot/scheduler"
	onebot2 "bcjh-bot/scheduler/onebot"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"fmt"
	"strconv"
	"strings"
)

func MaterialQuery(c *scheduler.Context) {
	if strings.TrimSpace(c.PretreatedMessage) == "" {
		_, _ = c.Reply(materialHelp())
		return
	}
	args := strings.Split(strings.TrimSpace(c.PretreatedMessage), " ")
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
	err := dao.DB.Where("name like ?", "%"+args[0]+"%").Find(&materials)
	if err != nil {
		logger.Error("数据库查询出错!")
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}
	if len(materials) == 0 {
		_, _ = c.Reply(fmt.Sprintf("没有找到叫%s的食材", args[0]))
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
			_, _ = c.Reply(msg)
			return
		}
	}

	recipeMaterials := make([]database.RecipeMaterial, 0)
	err = dao.DB.Where("material_id = ?", materials[0].MaterialId).Find(&recipeMaterials)
	if err != nil {
		logger.Error("数据库查询出错!")
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}
	if len(recipeMaterials) == 0 {
		_, _ = c.Reply(fmt.Sprintf("没有使用%s的菜谱哦~", args[0]))
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
	err = dao.DB.In("gallery_id", recipeGalleryIds).Find(&recipes)
	if err != nil {
		logger.Error("数据库查询出错!")
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}

	for i := range recipes {
		recipes[i].MaterialEfficiency = recipeMaterialMap[recipes[i].GalleryId]
	}

	var note string
	recipes, note = orderRecipes(recipes, "耗材效率")
	if note != "" {
		_, _ = c.Reply(note)
		return
	}

	// 处理消息
	var msg string
	listLength := config.AppConfig.Bot.GroupMsgMaxLen
	if c.GetMessageType() == onebot2.MessageTypePrivate {
		listLength = config.AppConfig.Bot.PrivateMsgMaxLen
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
	_, _ = c.Reply(msg)
}
