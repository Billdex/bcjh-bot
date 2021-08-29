package messageservice

import (
	. "bcjh-bot/middleware"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util/logger"
	"fmt"
	"time"
)

func Register(s *scheduler.Scheduler) {
	g := s.Group("*")
	g.Use(CheckBotState)
	g.Bind("开机", MustAdmin, EnableBotInGroup)
	g.Bind("关机", MustAdmin, DisableBotInGroup)
	g.Bind("reply", replyMessage)
}

func replyMessage(c *scheduler.Context) {
	msg := fmt.Sprintf("[消息类型]: %s\n", c.GetMessageType())
	msg += fmt.Sprintf("[发送人]: %s\n", c.GetSenderNickname())
	msg += fmt.Sprintf("[消息内容]: %s\n", c.PretreatedMessage)
	msg += fmt.Sprintf("[发送时间]: %s", time.Unix(c.GetEventTime(), 0).Format("2006-01-02 15:04:05"))
	_, _ = c.Reply(msg)
}

func EnableBotInGroup(c *scheduler.Context) {
	if c.GetMessageType() != onebot.MessageTypeGroup || c.GetGroupEvent() == nil {
		return
	}
	atList := c.GetAtList()
	for _, at := range atList {
		if at == c.GetBot().BotId {
			event := c.GetGroupEvent()
			if err := SetBotState(event.SelfId, event.GroupId, true); err != nil {
				logger.Error("设置群内机器人启动出错:", err)
				return
			}
			_, _ = c.Reply("已开机")
			break
		}
	}
}

func DisableBotInGroup(c *scheduler.Context) {
	if c.GetMessageType() != onebot.MessageTypeGroup || c.GetGroupEvent() == nil {
		return
	}
	atList := c.GetAtList()
	for _, at := range atList {
		if at == c.GetBot().BotId {
			event := c.GetGroupEvent()
			if err := SetBotState(event.SelfId, event.GroupId, false); err != nil {
				logger.Error("设置群内机器人关闭出错:", err)
				return
			}
			_, _ = c.Reply("已关机")
			break
		}
	}

}
