package messageservice

import (
	. "bcjh-bot/middleware"
	"bcjh-bot/scheduler"
	"fmt"
	"time"
)

func Register(s *scheduler.Scheduler) {
	g := s.Group("*")
	g.Use(CheckBotState)
	g.Use(CheckBlackList)
	// 管理功能
	g.Bind("开机", MustAdmin, EnableBotInGroup)
	g.Bind("关机", MustAdmin, DisableBotInGroup)
	g.Bind("启用", MustAdmin, EnablePluginInGroup)
	g.Bind("停用", MustAdmin, DisablePluginInGroup)
	g.Bind("ban", MustAdmin, BanUser)
	g.Bind("allow", MustAdmin, AllowUser)
	g.Bind("公告", MustSuperAdmin, CheckPluginState(true), PublicNotice)
	g.Bind("改命", MustSuperAdmin, ForceTarot).Alias("转运")

	// 其他查询功能
	g.Bind("反馈", CheckPluginState(true), Feedback).Alias("建议")
	g.Bind("抽签", CheckPluginState(false), Tarot).Alias("占卜", "求签", "运势", "卜卦")
	g.Bind("reply", replyMessage)
}

func replyMessage(c *scheduler.Context) {
	msg := fmt.Sprintf("[消息类型]: %s\n", c.GetMessageType())
	msg += fmt.Sprintf("[发送人]: %s\n", c.GetSenderNickname())
	msg += fmt.Sprintf("[消息内容]: %s\n", c.PretreatedMessage)
	msg += fmt.Sprintf("[发送时间]: %s", time.Unix(c.GetEventTime(), 0).Format("2006-01-02 15:04:05"))
	_, _ = c.Reply(msg)
}
