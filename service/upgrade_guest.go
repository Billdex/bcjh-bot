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

func UpgradeGuestQuery(c *onebot.Context, args []string) {
	logger.Info("碰瓷数据查询，参数:", args)

	if len(args) == 0 {
		err := bot.SendMessage(c, upgradeGuestHelp())
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
		return
	}
	page := 1
	if len(args) == 2 {
		if strings.HasPrefix(args[1], "p") || strings.HasPrefix(args[1], "P") {
			num, err := strconv.Atoi(args[1][1:])
			if err != nil {
				logger.Error("分页参数转数字出错!", err)
				_ = bot.SendMessage(c, "查询参数出错!")
				return
			} else {
				if num < 1 {
					num = 1
				}
				page = num
			}
		}
	}

	guests := make([]database.GuestGift, 0)
	err := database.DB.Where("guest_id = ?", args[0]).Find(&guests)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_ = bot.SendMessage(c, "查询数据失败!")
		return
	}

	if len(guests) == 0 {
		err = database.DB.Where("guest_name like ?", "%"+args[0]+"%").Find(&guests)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			_ = bot.SendMessage(c, "查询数据失败!")
			return
		}
	}

	guestInfo := make(map[string]string)
	for _, guest := range guests {
		key := guest.GuestId
		guestInfo[key] = guest.GuestName
	}
	var msg string
	if len(guestInfo) == 0 {
		msg = "哎呀，好像找不到这个贵客呢!"
	} else if len(guestInfo) == 1 {
		recipes := make([]database.Recipe, 0)
		var guestName string
		for _, guest := range guestInfo {
			guestName = guest
		}
		err = database.DB.Where("guests like ?", "%\""+guestName+"\"%").Asc("Time").Find(&recipes)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			_ = bot.SendMessage(c, "查询数据失败!")
			return
		}
		if len(recipes) == 0 {
			msg = fmt.Sprintf("%s没有碰瓷数据哦!", guestName)
		} else {
			results := make([]string, 0)
			for _, recipe := range recipes {
				var upgrade string
				for p, _ := range recipe.Guests {
					if recipe.Guests[p] == guestName {
						switch p {
						case 0:
							upgrade = "优"
						case 1:
							upgrade = "特"
						case 2:
							upgrade = "神"
						}
					}
				}
				results = append(results, fmt.Sprintf("%s: %s", upgrade, recipe.Name))
			}

			listLength := util.MaxQueryListLength
			maxPage := (len(results)-1)/listLength + 1
			if len(results) > listLength {
				if page > maxPage {
					page = maxPage
				}
				msg += fmt.Sprintf("以下菜谱可碰瓷%s: (%d/%d)", guestName, page, maxPage)
			} else {
				msg += fmt.Sprintf("以下菜谱可碰瓷%s:", guestName)
			}
			for i := (page - 1) * listLength; i < page*listLength && i < len(results); i++ {
				msg += fmt.Sprintf("\n%s", results[i])
			}
			if page < maxPage {
				msg += "\n......"
			}
		}
	} else {
		msg = "想查哪个升阶贵客数据呢:"
		p := 0
		for k, _ := range guestInfo {
			msg += fmt.Sprintf("\n%s %s", k, guestInfo[k])
			if p == util.MaxQueryListLength-1 {
				msg += "\n......"
				break
			}
			p++
		}
	}

	err = bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}
