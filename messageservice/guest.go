package messageservice

import (
	"bcjh-bot/config"
	"bcjh-bot/dao"
	"bcjh-bot/model/database"
	"bcjh-bot/scheduler"
	"bcjh-bot/util"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func GuestQuery(c *scheduler.Context) {
	arg := c.PretreatedMessage
	re, err := regexp.Compile(strings.ReplaceAll(arg, "%", ".*"))
	if err != nil {
		_, _ = c.Reply("查询格式有误")
		return
	}
	guestGifts, err := dao.FindAllGuestGifts()
	if err != nil {
		logger.Error("查询贵客礼物数据失败", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}
	var guestName string
	mGuests := make(map[string]database.GuestGift)
	numId, _ := strconv.Atoi(arg)
	guestId := fmt.Sprintf("%03d", numId)
	for i := range guestGifts {
		// 如果有贵客 id 或贵客名完全匹配，则视为查询该贵客
		if guestGifts[i].GuestId == guestId || guestGifts[i].GuestName == arg {
			guestName = guestGifts[i].GuestName
			break
		}
		// 模糊匹配到则把结果加到列表里
		if re.MatchString(guestGifts[i].GuestName) {
			mGuests[guestGifts[i].GuestName] = guestGifts[i]
		}
	}

	if guestName == "" {
		if len(mGuests) == 0 {
			_, _ = c.Reply(fmt.Sprintf("唔, %s未曾光临本店呢", arg))
			return
		} else if len(mGuests) == 1 {
			for _, v := range mGuests {
				guestName = v.GuestName
			}
		} else {
			msg := "查询到以下贵客"
			cnt := 0
			for _, guest := range mGuests {
				cnt++
				if cnt > config.AppConfig.Bot.GroupMsgMaxLen {
					msg += "\n..."
					break
				}
				msg += fmt.Sprintf("\n%s %s", guest.GuestId, guest.GuestName)
			}
			_, _ = c.Reply(msg)
			return
		}
	}

	guests := make([]database.GuestGift, 0)
	for _, guest := range guestGifts {
		if guest.GuestName == guestName {
			guests = append(guests, guest)
		}
	}
	sort.Slice(guests, func(i, j int) bool {
		return guests[i].TotalTime < guests[j].TotalTime
	})
	msg := fmt.Sprintf("%s %s", guests[0].GuestId, guests[0].GuestName)
	for _, guestInfo := range guests {
		msg += fmt.Sprintf("\n%s-%s %s", guestInfo.Recipe, guestInfo.Antique, util.FormatSecondToString(guestInfo.TotalTime))
	}
	_, _ = c.Reply(msg)
}
