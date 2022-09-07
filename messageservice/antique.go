package messageservice

import (
	"bcjh-bot/config"
	"bcjh-bot/dao"
	"bcjh-bot/model/database"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"fmt"
	"strconv"
	"strings"
)

func AntiqueQuery(c *scheduler.Context) {
	args := strings.Split(c.PretreatedMessage, " ")
	if args[0] == "" {
		_, _ = c.Reply(antiqueHelp())
		return
	}

	antique := args[0]

	guests := make([]database.GuestGift, 0)
	err := dao.DB.Where("antique like ?", "%"+antique+"%").Asc("total_time").Find(&guests)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}
	if len(guests) == 0 {
		_, _ = c.Reply("没有找到符文数据")
		return
	}

	antiqueMap := make(map[string]string)
	for _, guest := range guests {
		antiqueMap[guest.Antique] = guest.Antique
	}
	if len(antiqueMap) > 1 {
		msg := "你要找哪个符文的数据呢:"
		for _, v := range antiqueMap {
			msg += fmt.Sprintf("\n%s", v)
		}
		_, _ = c.Reply(msg)
		return
	}

	page := 1
	if len(args) > 1 {
		if util.HasPrefixIn(args[1], "p", "P") {
			num, err := strconv.Atoi(args[1][1:])
			if err != nil {
				_, _ = c.Reply("分页参数有误")
				return
			} else {
				if num > 0 {
					page = num
				}
			}
		}
	}

	var msg string
	listLength := config.AppConfig.Bot.GroupMsgMaxLen
	if c.GetMessageType() == onebot.MessageTypePrivate {
		listLength = config.AppConfig.Bot.PrivateMsgMaxLen
	}
	maxPage := (len(guests)-1)/listLength + 1
	if len(guests) > listLength {
		if page > maxPage {
			page = maxPage
		}
		msg += fmt.Sprintf("以下菜有概率得%s: (%d/%d)", guests[0].Antique, page, maxPage)
	} else {
		msg += fmt.Sprintf("以下菜有概率得%s:", guests[0].Antique)
	}
	for i := (page - 1) * listLength; i < page*listLength && i < len(guests); i++ {
		totalTime := util.FormatSecondToString(guests[i].TotalTime)
		msg += fmt.Sprintf("\n%s-%s %s", guests[i].Recipe, guests[i].GuestName, totalTime)
	}
	if page < maxPage {
		msg += "\n......"
	}

	_, _ = c.Reply(msg)
}
