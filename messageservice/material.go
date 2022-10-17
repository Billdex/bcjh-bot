package messageservice

import (
	"bcjh-bot/config"
	"bcjh-bot/dao"
	"bcjh-bot/model/database"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"fmt"
	"strconv"
	"strings"
)

func MaterialQuery(c *scheduler.Context) {
	args := strings.Split(c.PretreatedMessage, " ")
	page := 1
	if len(args) > 1 {
		if strings.HasPrefix(args[1], "p") || strings.HasPrefix(args[1], "P") {
			num, err := strconv.Atoi(args[1][1:])
			if err != nil {
				_, _ = c.Reply("分页参数错误")
			} else {
				if num < 1 {
					num = 1
				}
				page = num
			}
		}
	}

	// 查出所有食材，假设存在完全匹配的则只使用该食材筛选。
	materials, err := dao.SearchMaterialsWithName(args[0])
	if err != nil {
		logger.Error("根据名称搜索食材失败", err)
		_, _ = c.Reply(e.SystemErrorNote)
	}
	if len(materials) == 0 {
		_, _ = c.Reply(fmt.Sprintf("%s是什么食材呀", args[0]))
	}
	// 查询到多个食材时检查是否有完整匹配的，有则直接按照原值筛选菜谱，没有则返回食材列表
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
	// 查出所有菜谱，根据菜谱数据取结果
	allRecipes, err := dao.FindAllRecipes()
	if err != nil {
		logger.Error("查询菜谱数据失败", err)
		_, _ = c.Reply(e.SystemErrorNote)
	}
	recipes := make([]database.Recipe, 0)
	for _, recipe := range allRecipes {
		for i := range recipe.Materials {
			if recipe.Materials[i].Material.Name == args[0] {
				recipe.MaterialEfficiency = recipe.Materials[i].Efficiency
				recipes = append(recipes, recipe)
				break
			}
		}
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
	if c.GetMessageType() == onebot.MessageTypePrivate {
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
	_, _ = c.Reply(msg)
}
