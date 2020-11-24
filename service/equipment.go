package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/model/database"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"fmt"
)

func EquipmentQuery(c *onebot.Context, args []string) {
	logger.Info("厨具查询，参数:", args)
	if len(args) == 0 {
		err := bot.SendMessage(c,
			fmt.Sprintf("指令示例:\n"+
				"%s厨具 金烤叉", util.PrefixCharacter))
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

	equips := make([]database.Equip, 0)
	err := database.DB.Where("gallery_id = ?", args[0]).Asc("gallery_id").Find(&equips)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_ = bot.SendMessage(c, "查询数据失败!")
		return
	}
	if len(equips) == 0 {
		err = database.DB.Where("name like ?", "%"+args[0]+"%").Asc("gallery_id").Find(&equips)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			_ = bot.SendMessage(c, "查询数据失败!")
			return
		}
	}

	var msg string
	if len(equips) == 0 {
		msg = "哎呀，好像找不到呢!"
	} else if len(equips) == 1 {
		equip := equips[0]
		rarity := ""
		for i := 0; i < equip.Rarity; i++ {
			rarity += "🔥"
		}
		skills := ""
		for p, skillId := range equip.Skills {
			skill := new(database.Skill)
			has, err := database.DB.Where("skill_id = ?", skillId).Get(skill)
			if err != nil {
				logger.Error("查询数据库出错!", err)
				_ = bot.SendMessage(c, "查询数据失败!")
				return
			}
			if has {
				skills += skill.Description
				if p != len(equip.Skills)-1 {
					skills += ","
				}
			}
		}
		msg += fmt.Sprintf("%s %s\n", equip.GalleryId, equip.Name)
		msg += fmt.Sprintf("%s\n", rarity)
		msg += fmt.Sprintf("来源: %s\n", equip.Origin)
		msg += fmt.Sprintf("效果: %s", skills)

	} else {
		msg = "查询到以下厨具:\n"
		for p, equip := range equips {
			msg += fmt.Sprintf("%s %s", equip.GalleryId, equip.Name)
			if p != len(equips)-1 {
				msg += "\n"
				if p == util.MaxSearchList-1 {
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
