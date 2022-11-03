package messageservice

import (
	"bcjh-bot/config"
	"bcjh-bot/dao"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func UpgradeGuestQuery(c *scheduler.Context) {
	args := strings.Split(c.PretreatedMessage, " ")
	page := 1
	if len(args) == 2 {
		if util.HasPrefixIn(args[1], "p", "P") {
			num, err := strconv.Atoi(strings.Trim(args[1][1:], "-"))
			if err != nil {
				logger.Error("分页参数转数字出错!", err)
				_, _ = c.Reply("查询参数出错!")
				return
			} else {
				if num > 0 {
					page = num
				}
			}
		}
	}

	re, err := regexp.Compile(strings.ReplaceAll(args[0], "%", ".*"))
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
	mGuestNames := make(map[string]string)
	numId, _ := strconv.Atoi(args[0])
	guestId := fmt.Sprintf("%03d", numId)
	for i := range guestGifts {
		// 如果有贵客 id 或贵客名完全匹配，则视为查询该贵客
		if guestGifts[i].GuestId == guestId || guestGifts[i].GuestName == args[0] {
			guestName = guestGifts[i].GuestName
			mGuestNames = map[string]string{guestGifts[i].GuestId: guestGifts[i].GuestName}
			break
		}
		// 模糊匹配到则把结果加到列表里
		if re.MatchString(guestGifts[i].GuestName) {
			guestName = guestGifts[i].GuestName
			mGuestNames[guestGifts[i].GuestId] = guestGifts[i].GuestName
		}
	}

	if len(mGuestNames) == 0 {
		_, _ = c.Reply(fmt.Sprintf("唔, %s未曾光临本店呢", args[0]))
		return
	} else if len(mGuestNames) == 1 {
		// 筛选出包含该升阶贵客的菜谱
		allRecipes, err := dao.FindAllRecipes()
		if err != nil {
			logger.Error("获取菜谱数据失败", err)
			_, _ = c.Reply(e.SystemErrorNote)
			return
		}
		results := make([]string, 0)
		for _, recipe := range allRecipes {
			if upgrade, ok := recipe.HasUpgradeGuest(guestName); ok {
				results = append(results, fmt.Sprintf("%s: %s", upgrade, recipe.Name))
			}
		}

		listLength := config.AppConfig.Bot.GroupMsgMaxLen
		if c.GetRawMessage() == onebot.MessageTypePrivate {
			listLength = config.AppConfig.Bot.PrivateMsgMaxLen
		}
		msg := util.PaginationOutput(results, page, listLength,
			fmt.Sprintf("以下菜谱可碰瓷%s", guestName),
			func(s string) string {
				return s
			})
		_, _ = c.Reply(msg)
		return
	} else {
		msg := "想查哪个升阶贵客数据呢:"
		p := 0
		for k := range mGuestNames {
			msg += fmt.Sprintf("\n%s %s", k, mGuestNames[k])
			if p == config.AppConfig.Bot.GroupMsgMaxLen-1 {
				msg += "\n......"
				break
			}
			p++
		}
		_, _ = c.Reply(msg)
		return
	}
}
