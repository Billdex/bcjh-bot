package messageservice

import (
	. "bcjh-bot/middleware"
	"bcjh-bot/scheduler"
)

func Register(s *scheduler.Scheduler) {
	g := s.Group("#").Alias("＃")
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
	g.Bind("符文", CheckPluginState(true), AntiqueQuery).Alias("礼物")
	g.Bind("调料", CheckPluginState(true), CondimentQuery)
	g.Bind("任务", CheckPluginState(true), TaskQuery).Alias("主线", "支线")
	g.Bind("限时任务", CheckPluginState(true), TimeLimitingTaskQuery).Alias("限时攻略", "限时支线")
	g.Bind("攻略", CheckPluginState(true), StrategyQuery)
	g.Bind("碰瓷", CheckPluginState(false), UpgradeGuestQuery).Alias("升阶贵客", "升级贵客")
	g.Bind("后厨", CheckPluginState(true), ComboQuery).Alias("合成")
	g.Bind("兑换码", CheckPluginState(true), ExchangeQuery).Alias("玉璧")
	g.Bind("实验室", CheckPluginState(true), LaboratoryQuery).Alias("研究")

	// 快捷查询
	g.Bind("图鉴网", GalleryWebsite).Alias("图鉴")
	g.Bind("白菜菊花", BCJHAppDownload)
	g.Bind("计算器", Calculator).Alias("计算机")
	g.Bind("游戏术语", TermInfo).Alias("黑话", "术语")

	// 娱乐功能
	g.Bind("抽签", CheckPluginState(false), Tarot).Alias("占卜", "求签", "运势", "卜卦", "占卦")
}
