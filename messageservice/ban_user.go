package messageservice

import (
	"bcjh-bot/dao"
	"bcjh-bot/global"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func BanUser(c *scheduler.Context) {
	if c.GetMessageType() != onebot.MessageTypeGroup || c.GetGroupEvent() == nil {
		return
	}
	atList := c.GetAtList()
	if len(atList) == 0 {
		_, _ = c.Reply("请at需要禁用的用户")
		return
	}
	stringBanTime := matchStringTime(c.PretreatedMessage)
	if stringBanTime == "" {
		_, _ = c.Reply("未填写禁用时长或格式错误")
		return
	}
	banTime := stringTimeToSecond(stringBanTime)
	if banTime < 0 {
		_, _ = c.Reply("未填写禁用时长或格式错误")
		return
	}
	if banTime > 30*24*60*60 {
		_, _ = c.Reply("禁用时长不能超过三十天")
		return
	}
	endTime := time.Now().Unix() + banTime
	successList := make([]int64, 0)
	failList := make([]int64, 0)
	for _, qq := range atList {
		if global.IsSuperAdmin(qq) {
			_, _ = c.Reply(fmt.Sprintf(e.PermissionDeniedNote))
		} else {
			err := dao.SetUserBanTime(qq, c.GetGroupId(), endTime)
			if err != nil {
				logger.Errorf("%d 加入黑名单失败 %v", qq, err)
				failList = append(failList, qq)
			} else {
				successList = append(successList, qq)
			}
		}
	}
	var successMsg, failMsg string
	if len(successList) > 0 {
		successMsg = fmt.Sprintf("%+v 已被禁用至%s", successList, time.Unix(endTime, 0).Format("2006-01-02 15:04:05"))
	}
	if len(failList) > 0 {
		failMsg = fmt.Sprintf("%+v 加入禁用名单失败", failList)
	}
	_, _ = c.Reply(strings.Join([]string{successMsg, failMsg}, "\n"))
}

func AllowUser(c *scheduler.Context) {
	if c.GetMessageType() != onebot.MessageTypeGroup || c.GetGroupEvent() == nil {
		return
	}
	atList := c.GetAtList()
	if len(atList) == 0 {
		_, _ = c.Reply("请at需要移出禁用名单的用户")
		return
	}
	successList := make([]int64, 0)
	failList := make([]int64, 0)
	for _, qq := range atList {
		// 禁用时间设为当前时间即可视为移出黑名单
		err := dao.SetUserBanTime(qq, c.GetGroupId(), time.Now().Unix())
		if err != nil {
			logger.Errorf("%d 移出禁用名单失败 %v", qq, err)
			failList = append(failList, qq)
		} else {
			successList = append(successList, qq)
		}
	}
	var successMsg, failMsg string
	if len(successList) > 0 {
		successMsg = fmt.Sprintf("%+v 已移出禁用名单", successList)
	}
	if len(failList) > 0 {
		failMsg = fmt.Sprintf("%+v 移出禁用名单失败", failList)
	}
	_, _ = c.Reply(strings.Join([]string{successMsg, failMsg}, "\n"))
}

func matchStringTime(s string) string {
	reg := `(\d+d\d+h\d+m)|(\d+d\d+h)|(\d+d\d+m)|(\d+h\d+m)|(\d+d)|(\d+h)|(\d+m)`
	pattern := regexp.MustCompile(reg)
	allIndexes := pattern.FindAllSubmatch([]byte(s), -1)
	for _, loc := range allIndexes {
		if len(loc[0]) != 0 {
			return string(loc[0])
		}
	}
	return ""
}

func stringTimeToSecond(s string) int64 {
	sumTime := 0
	buf := ""
	for _, c := range s {
		if c == 'd' {
			num, err := strconv.Atoi(buf)
			if err != nil {
				return -1
			}
			sumTime += num * 60 * 60 * 24
			buf = ""
		} else if c == 'h' {
			num, err := strconv.Atoi(buf)
			if err != nil {
				return -1
			}
			sumTime += num * 60 * 60
			buf = ""

		} else if c == 'm' {
			num, err := strconv.Atoi(buf)
			if err != nil {
				return -1
			}
			sumTime += num * 60
			buf = ""
		} else {
			buf += string(c)
		}
	}
	return int64(sumTime)
}
