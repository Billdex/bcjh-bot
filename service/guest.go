package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/model/database"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"fmt"
)

func GuestQuery(c *onebot.Context, args []string) {
	logger.Info("贵客查询，参数:", args)

	if len(args) == 0 {
		err := bot.SendMessage(c,
			fmt.Sprintf("指令示例:\n"+
				"%s贵客 木良", util.PrefixCharacter))
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
		key := fmt.Sprintf("%s %s", guest.GuestId, guest.GuestName)
		value, hasKey := guestInfo[key]
		if hasKey {
			value += "\n"
			value += fmt.Sprintf("%s-%s", guest.Antique, guest.Recipe)
		} else {
			value = fmt.Sprintf("%s %s\n", guest.GuestId, guest.GuestName)
			value += fmt.Sprintf("%s-%s", guest.Antique, guest.Recipe)
		}
		guestInfo[key] = value
	}
	var msg string
	if len(guestInfo) == 0 {
		msg = "哎呀，好像找不到呢!"
	} else if len(guestInfo) == 1 {
		key := fmt.Sprintf("%s %s", guests[0].GuestId, guests[0].GuestName)
		msg = guestInfo[key]
	} else {
		msg = "查询到以下贵客:"
		p := 0
		for k, _ := range guestInfo {
			msg += fmt.Sprintf("\n%s", k)
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
