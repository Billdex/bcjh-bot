package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/model/database"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"fmt"
	"math/rand"
	"time"
)

func Tarot(c *onebot.Context, _ []string) {
	if c.MessageType == util.OneBotMessageGroup {
		has, err := database.DB.Where("plugin = ? and group_id = ?", "tarot", c.GroupId).Exist(&database.WhiteList{})
		if err != nil {
			logger.Error("查询数据库出错", err)
			return
		}
		if !has {
			return
		}
	}
	sender := c.Sender.UserId
	now := time.Now()
	timeSeed := now.Unix()
	timeSeed -= int64(now.Hour() * 3600)
	timeSeed -= int64(now.Minute() * 60)
	timeSeed -= int64(now.Second())
	rand.Seed(int64(sender) + timeSeed)
	total, err := database.DB.Count(&database.Tarot{})
	if err != nil {
		logger.Error("查询数据库出错", err)
		_ = bot.SendMessage(c, util.SystemErrorNote)
		return
	}
	tarotId := rand.Int63n(total) + 1
	tarot := new(database.Tarot)
	_, err = database.DB.Where("id = ?", tarotId).Get(tarot)
	if err != nil {
		logger.Error("查询数据库出错", err)
		_ = bot.SendMessage(c, util.SystemErrorNote)
		return
	}
	var level string
	switch {
	case tarot.Score == 0:
		level = "不知道吉不吉"
	case 0 < tarot.Score && tarot.Score < 15:
		level = "小小吉"
	case 15 <= tarot.Score && tarot.Score < 40:
		level = "小吉"
	case 40 <= tarot.Score && tarot.Score < 60:
		level = "中吉"
	case 60 <= tarot.Score && tarot.Score < 85:
		level = "大吉"
	case 85 <= tarot.Score && tarot.Score < 100:
		level = "大大吉"
	case tarot.Score == 100:
		level = "超吉"
	default:
		level = "?"
	}
	msg := fmt.Sprintf("[%s]抽了一根签\n", c.Sender.Nickname)
	msg += fmt.Sprintf("运势指数:%d [%s]\n", tarot.Score, level)
	msg += fmt.Sprintf("签上说:\n%s", tarot.Description)
	_ = bot.SendMessage(c, msg)
	return
}
