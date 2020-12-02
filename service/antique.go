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

func AntiqueQuery(c *onebot.Context, args []string) {
	logger.Info("符文查询, 参数:", args)
	if len(args) == 0 {
		err := bot.SendMessage(c, antiqueHelp())
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
		return
	}

	guests := make([]database.GuestGift, 0)
	err := database.DB.Where("antique like ?", "%"+args[0]+"%").Find(&guests)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_ = bot.SendMessage(c, util.SystemErrorNote)
		return
	}
	if len(guests) == 0 {
		_ = bot.SendMessage(c, "没有找到符文数据")
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

	recipesName := make([]string, 0)
	for _, guest := range guests {
		recipesName = append(recipesName, guest.Recipe)
	}
	recipes := make([]database.Recipe, 0)
	err = database.DB.In("name", recipesName).OrderBy("`total_time` ASC").Find(&recipes)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_ = bot.SendMessage(c, util.SystemErrorNote)
		return
	}

	var msg string
	listLength := util.MaxQueryListLength
	maxPage := (len(recipes)-1)/listLength + 1
	if len(recipes) > listLength {
		if page > maxPage {
			page = maxPage
		}
		msg += fmt.Sprintf("这里有你想点的菜吗: (%d/%d)\n", page, maxPage)
	} else {
		msg += "这里有你想点的菜吗:\n"
	}
	for i := (page - 1) * listLength; i < page*listLength && i < len(recipes); i++ {
		totalTime := util.FormatSecondToString(recipes[i].Time * recipes[i].Limit)
		msg += fmt.Sprintf("%s %s %s", recipes[i].GalleryId, recipes[i].Name, totalTime)
		if i < page*listLength-1 && i < len(recipes)-1 {
			msg += "\n"
		}
	}
	if page < maxPage {
		msg += "\n......"
	}

	logger.Info("发送菜谱查询结果:", msg)
	err = bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}

}
