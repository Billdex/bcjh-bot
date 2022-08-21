package messageservice

import (
	"bcjh-bot/dao"
	"bcjh-bot/model/database"
	"bcjh-bot/scheduler"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func Tarot(c *scheduler.Context) {
	now := time.Now()
	timeSeed := now.Unix()
	timeSeed -= int64(now.Hour() * 3600)
	timeSeed -= int64(now.Minute() * 60)
	timeSeed -= int64(now.Second())
	total, err := dao.DB.Count(&database.Tarot{})
	if err != nil {
		logger.Error("查询数据库出错", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}
	selfRand := rand.New(rand.NewSource(c.GetSenderId() + timeSeed))
	tarotId := selfRand.Int63n(total) + 1
	tarot := new(database.Tarot)
	_, err = dao.DB.Where("id = ?", tarotId).Get(tarot)
	if err != nil {
		logger.Error("查询数据库出错", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}
	if tarot.Score == 99 && c.GetSenderId() != 1726688182 {
		_, err = dao.DB.Where("score = ?", 97).OrderBy("id").Limit(1).Get(tarot)
		if err != nil {
			logger.Error("查询数据库出错", err)
			_, _ = c.Reply(e.SystemErrorNote)
			return
		}
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
	msg := fmt.Sprintf("[%s]抽了一根签\n", c.GetSenderNickname())
	msg += fmt.Sprintf("运势指数:%d [%s]\n", tarot.Score, level)
	msg += fmt.Sprintf("签上说:\n%s", tarot.Description)
	_, _ = c.Reply(msg)
	return
}

func ForceTarot(c *scheduler.Context) {
	num, err := strconv.Atoi(c.PretreatedMessage)
	if err != nil {
		_, _ = c.Reply("请输入正确的数字")
		return
	}
	if num == 99 && c.GetSenderId() != 1726688182 {
		_, _ = c.Reply("改命失败")
		return
	}
	tarotList := make([]database.Tarot, 0)
	err = dao.DB.Where("score = ?", num).Find(&tarotList)
	if err != nil || len(tarotList) == 0 {
		_, _ = c.Reply("改命失败")
		return
	}
	tarot := tarotList[rand.Intn(len(tarotList))]
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
	msg := fmt.Sprintf("[%s]抽了一根签\n", c.GetSenderNickname())
	msg += fmt.Sprintf("运势指数:%d [%s]\n", tarot.Score, level)
	msg += fmt.Sprintf("签上说:\n%s", tarot.Description)
	_, _ = c.Reply(msg)

}
