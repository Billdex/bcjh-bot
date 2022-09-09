package messageservice

import (
	"bcjh-bot/dao"
	"bcjh-bot/model/database"
	"bcjh-bot/scheduler"
	"bcjh-bot/util/logger"
)

const (
	close = iota
	open
	finish
)

func Feedback(c *scheduler.Context) {
	feedback := new(database.Feedback)
	feedback.Sender = c.GetSenderId()
	feedback.Nickname = c.GetSenderNickname()
	feedback.Message = c.PretreatedMessage
	feedback.Status = open
	affected, err := dao.DB.Insert(feedback)
	if err != nil {
		logger.Error("插入数据库出错:", err)
		_, _ = c.Reply("不小心反馈失败了呢!")
		return
	}
	if affected == 0 {
		_, _ = c.Reply("不小心反馈失败了呢!")
		return
	}
	_, _ = c.Reply("收到! 反馈信息已经记在小本本上啦~")
}
