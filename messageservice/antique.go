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
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func AntiqueQuery(c *scheduler.Context) {
	args := strings.Split(c.PretreatedMessage, " ")
	antique := args[0]

	re, err := regexp.Compile(strings.ReplaceAll(antique, "%", ".*"))
	if err != nil {
		_, _ = c.Reply("查询格式有误")
		return
	}

	guestGifts, err := dao.FindAllGuestGifts()
	if err != nil {
		logger.Error("查询贵客礼物数据出错!")
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}

	antiqueMap := make(map[string]struct{})
	gifts := make([]database.GuestGift, 0)
	for _, guestGift := range guestGifts {
		if re.MatchString(guestGift.Antique) {
			antiqueMap[guestGift.Antique] = struct{}{}
			gifts = append(gifts, guestGift)
		}
	}
	if len(antiqueMap) == 0 {
		_, _ = c.Reply("没有找到符文数据")
		return
	}
	if len(antiqueMap) > 1 {
		msg := "你要找哪个符文的数据呢"
		for antique := range antiqueMap {
			msg += fmt.Sprintf("\n%s", antique)
		}
		_, _ = c.Reply(msg)
		return
	}
	// 最终结果按照从菜谱总耗时升序排列
	sort.Slice(gifts, func(i, j int) bool {
		return gifts[i].TotalTime < gifts[j].TotalTime
	})

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

	listLength := config.AppConfig.Bot.GroupMsgMaxLen
	if c.GetMessageType() == onebot.MessageTypePrivate {
		listLength = config.AppConfig.Bot.PrivateMsgMaxLen
	}
	msg := util.PaginationOutput(gifts, page, listLength,
		fmt.Sprintf("以下菜谱概率获得%s", gifts[0].Antique),
		func(gift database.GuestGift) string {
			return fmt.Sprintf("%s-%s %s", gift.Recipe, gift.GuestName, util.FormatSecondToString(gift.TotalTime))
		})

	_, _ = c.Reply(msg)
}
