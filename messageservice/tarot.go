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

// Tarot 抽签占卜
func Tarot(c *scheduler.Context) {
	y, m, d := time.Now().Date()
	timeSeed := time.Date(y, m, d, 0, 0, 0, 0, time.Local).Unix()
	tarots, err := dao.FindAllTarots()
	if err != nil {
		logger.Error("查询签文信息出错", err)
		_, _ = c.Reply(e.SystemErrorNote)
	}
	selfRand := rand.New(rand.NewSource(c.GetSenderId() + timeSeed))
	tarotId := selfRand.Int63n(int64(len(tarots))) + 1
	var tarot database.Tarot
	for i := range tarots {
		if int64(tarots[i].Id) == tarotId {
			tarot = tarots[i]
			break
		}
	}
	if tarot.Score == 99 && c.GetSenderId() != 1726688182 {
		score := selfRand.Int63n(98)
		for i := range tarots {
			if int64(tarots[i].Score) == score {
				tarot = tarots[i]
				break
			}
		}
	}
	msg := fmt.Sprintf("[%s]抽了一根签\n", c.GetSenderNickname())
	msg += fmt.Sprintf("运势指数 %d [%s]\n", tarot.Score, tarot.Level())
	msg += fmt.Sprintf("签上说:\n%s", tarot.Description)
	_, _ = c.Reply(msg)
	return
}

// ForceTarot 强抽一签
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
	tarots, err := dao.FindTarotsWithScore(num)
	if err != nil || len(tarots) == 0 {
		_, _ = c.Reply("改命失败")
		return
	}
	tarot := tarots[rand.Intn(len(tarots))]
	msg := fmt.Sprintf("[%s]抽了一根签\n", c.GetSenderNickname())
	msg += fmt.Sprintf("运势指数 %d [%s]\n", tarot.Score, tarot.Level())
	msg += fmt.Sprintf("签上说:\n%s", tarot.Description)
	_, _ = c.Reply(msg)
	return
}
