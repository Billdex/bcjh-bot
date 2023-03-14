package messageservice

import (
	. "bcjh-bot/middleware"
	"bcjh-bot/scheduler"
)

func Register(s *scheduler.Scheduler) {
	g := s.Group("#").Alias("＃")
	g.Use(Recovery)
	g.Use(QueryLog)
	g.Use(CheckBotState)
	g.Use(CheckBlackList)
	g.Use(MergeRepeatSpace)
	// 管理功能
	g.Bind("开机", MustAtSelf, MustAdmin, EnableBotInGroup)
	g.Bind("关机", MustAtSelf, MustAdmin, DisableBotInGroup)
	g.Bind("允许私聊", MustAtSelf, MustSuperAdmin, AllowPrivate).Alias("开启私聊")
	g.Bind("禁用私聊", MustAtSelf, MustSuperAdmin, DisablePrivate).Alias("关闭私聊", "禁止私聊")
	g.Bind("启用", MustAdmin, EnablePluginInGroup)
	g.Bind("停用", MustAdmin, DisablePluginInGroup)
	g.Bind("ban", MustAdmin, BanUser)
	g.Bind("allow", MustAdmin, AllowUser)
	g.Bind("更新", MustSuperAdmin, UpdateData)
	g.Bind("公告", MustSuperAdmin, CheckPluginState(true), PretreatedImage, PublicNotice)
	g.Bind("改命", MustSuperAdmin, ForceTarot).Alias("转运")

	// 常规查询
	g.Bind("帮助", CheckPluginState(true), HelpGuide).Alias("功能", "说明", "指引", "使用说明")
	g.Bind("反馈", CheckPluginState(true), Helper(feedbackHelp), Feedback).Alias("建议")
	g.Bind("厨师", CheckPluginState(true), Helper(chefHelp), ChefQuery).Alias("厨子")
	g.Bind("菜谱", CheckPluginState(true), Helper(recipeHelp), RecipeQuery).Alias("食谱")
	g.Bind("厨具", CheckPluginState(true), Helper(equipmentHelp), EquipmentQuery).Alias("装备", "道具")
	g.Bind("食材", CheckPluginState(true), Helper(materialHelp), MaterialQuery).Alias("材料")
	g.Bind("贵客", CheckPluginState(true), Helper(guestHelp), GuestQuery).Alias("稀有客人", "贵宾", "客人", "宾客", "稀客")
	g.Bind("符文", CheckPluginState(true), Helper(antiqueHelp), AntiqueQuery).Alias("礼物")
	g.Bind("调料", CheckPluginState(true), Helper(condimentHelp), CondimentQuery)
	g.Bind("任务", CheckPluginState(true), Helper(questHelp), TaskQuery).Alias("主线", "支线")
	g.Bind("限时任务", CheckPluginState(true), TimeLimitingTaskQuery).Alias("限时攻略", "限时支线")
	g.Bind("攻略", CheckPluginState(true), PretreatedImage, Helper(strategyHelp), StrategyQuery)
	g.Bind("碰瓷", CheckPluginState(false), Helper(upgradeGuestHelp), UpgradeGuestQuery).Alias("升阶贵客", "升级贵客")
	g.Bind("后厨", CheckPluginState(true), Helper(comboHelp), ComboQuery).Alias("合成")
	g.Bind("兑换码", CheckPluginState(true), PretreatedImage, ExchangeQuery).Alias("玉璧", "兑奖码")
	g.Bind("实验室", CheckPluginState(true), Helper(LaboratoryHelp), LaboratoryQuery).Alias("研究")
	g.Bind("修炼", CheckPluginState(false), Helper(ultimateHelp), UltimateQuery)
	g.Bind("个人数据导入", CheckPluginState(true), Helper(importDataHelp), ImportUserData)

	// 快捷查询
	g.Bind("图鉴网", GalleryWebsite).Alias("图鉴")
	g.Bind("白菜菊花", BCJHAppDownload)
	g.Bind("计算器", Calculator).Alias("计算机")
	g.Bind("游戏术语", TermInfo).Alias("黑话", "术语")

	// 快捷搜索
	g.Bind("#", QuickSearch).Alias("＃")

	// 娱乐功能
	g.Bind("抽签", CheckPluginState(false), Tarot).Alias("占卜", "求签", "运势", "卜卦", "占卦")
	g.Bind("随机个人图鉴", CheckPluginState(false), RandChefImg)
}
