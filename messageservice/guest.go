package messageservice

import (
	"bcjh-bot/config"
	"bcjh-bot/model/database"
	"bcjh-bot/scheduler"
	"bcjh-bot/util"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"fmt"
	"strconv"
	"strings"
)

func GuestQuery(c *scheduler.Context) {
	arg := strings.TrimSpace(c.PretreatedMessage)
	if arg == "" {
		_, _ = c.Reply(guestHelp())
		return
	}

	argType := "guest_id"
	guests := make([]database.Guest, 0)
	numId, err := strconv.Atoi(arg)
	if err == nil {
		guestId := fmt.Sprintf("%03d", numId)
		err := database.DB.Where("guest_id = ?", guestId).Find(&guests)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			_, _ = c.Reply(e.SystemErrorNote)
			return
		}
	}
	if len(guests) == 0 {
		err = database.DB.Where("guest_name like ?", "%"+arg+"%").Find(&guests)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			_, _ = c.Reply(e.SystemErrorNote)
			return
		}
		argType = "guest_name"
	}

	// 查询到多个贵客时返回贵客列表
	if len(guests) > 1 {
		msg := "查询到以下贵客"
		for i, guest := range guests {
			if i > config.AppConfig.Bot.GroupMsgMaxLen-1 {
				msg += "\n..."
				break
			}
			msg += fmt.Sprintf("\n%s %s", guest.GuestId, guest.GuestName)
		}
		_, _ = c.Reply(msg)
		return
	} else if len(guests) == 0 {
		_, _ = c.Reply(fmt.Sprintf("%s是什么神秘贵客呀？", arg))
		return
	}

	guestsInfo := make([]database.GuestGift, 0)
	switch argType {
	case "guest_id":
		err = database.DB.Where("guest_id = ?", guests[0].GuestId).Asc("total_time").Find(&guestsInfo)
	case "guest_name":
		err = database.DB.Where("guest_name = ?", guests[0].GuestName).Asc("total_time").Find(&guestsInfo)
	default:
		err = database.DB.Where("guest_name = ?", guests[0].GuestName).Asc("total_time").Find(&guestsInfo)
	}
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}
	if len(guestsInfo) == 0 {
		_, _ = c.Reply(fmt.Sprintf("未查询到%s的贵客信息", arg))
		return
	}
	msg := fmt.Sprintf("%s %s", guests[0].GuestId, guests[0].GuestName)
	for _, guestInfo := range guestsInfo {
		msg += fmt.Sprintf("\n%s-%s %s", guestInfo.Recipe, guestInfo.Antique, util.FormatSecondToString(guestInfo.TotalTime))
	}
	_, _ = c.Reply(msg)
}
