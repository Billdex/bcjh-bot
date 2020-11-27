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

	guests := make([]database.Guest, 0)
	err := database.DB.Where("gallery_id = ?", args[0]).Asc("gallery_id").Find(&guests)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_ = bot.SendMessage(c, "查询数据失败!")
		return
	}

	if len(guests) == 0 {
		err = database.DB.Where("name like ?", "%"+args[0]+"%").Asc("gallery_id").Find(&guests)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			_ = bot.SendMessage(c, "查询数据失败!")
			return
		}
	}

	var msg string
	if len(guests) == 0 {
		msg = "哎呀，好像找不到呢!"
	} else if len(guests) == 1 {
		guest := guests[0]
		var gifts string
		for p, gift := range guest.Gifts {
			gifts += fmt.Sprintf("%s-%s", gift.Antique, gift.Recipe)
			if p != len(guest.Gifts)-1 {
				gifts += "\n"
			}
		}
		msg += fmt.Sprintf("%s %s\n", guest.GalleryId, guest.Name)
		msg += fmt.Sprintf("%s", gifts)
	} else {
		msg = "查询到以下贵客:\n"
		for p, guest := range guests {
			msg += fmt.Sprintf("%s %s", guest.GalleryId, guest.Name)
			if p != len(guests)-1 {
				msg += "\n"
				if p == util.MaxQueryListLength-1 {
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
