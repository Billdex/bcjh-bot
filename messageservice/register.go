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
	g.Bind("更新", MustSuperAdmin, UpdateData)
	g.Bind("公告", MustSuperAdmin, CheckPluginState(true), PublicNotice)
	g.Bind("改命", MustSuperAdmin, ForceTarot).Alias("转运")

	// 常规查询
	g.Bind("帮助", CheckPluginState(true), HelpGuide).Alias("功能", "说明", "指引", "使用说明")
	g.Bind("反馈", CheckPluginState(true), Feedback).Alias("建议")
	g.Bind("厨师", CheckPluginState(true), ChefQuery).Alias("厨子")
	g.Bind("菜谱", CheckPluginState(true), RecipeQuery).Alias("食谱")
	g.Bind("厨具", CheckPluginState(true), EquipmentQuery).Alias("装备", "道具")
	g.Bind("食材", CheckPluginState(true), MaterialQuery).Alias("材料")
	g.Bind("贵客", CheckPluginState(true), GuestQuery).Alias("稀有客人", "贵宾", "客人", "宾客", "稀客")

	// 快捷查询
	g.Bind("图鉴网", GalleryWebsite).Alias("图鉴")
	g.Bind("白菜菊花", BCJHAppDownload)
	g.Bind("计算器", Calculator).Alias("计算机")
	g.Bind("游戏术语", TermInfo).Alias("黑话", "术语")

	// 娱乐功能
	g.Bind("抽签", CheckPluginState(false), Tarot).Alias("占卜", "求签", "运势", "卜卦", "占卦")
	g.Bind("reply", replyMessage)
}

// 测试用, 后续会删掉
func replyMessage(c *scheduler.Context) {
	msg := fmt.Sprintf("[消息类型]: %s\n", c.GetMessageType())
	msg += fmt.Sprintf("[发送人]: %s\n", c.GetSenderNickname())
	msg += fmt.Sprintf("[消息内容]: %s\n", c.PretreatedMessage)
	msg += fmt.Sprintf("[发送时间]: %s", time.Unix(c.GetEventTime(), 0).Format("2006-01-02 15:04:05"))
	_, _ = c.Reply(msg)
}
