package messageservice

import (
	"bcjh-bot/dao"
	"bcjh-bot/global"
	"bcjh-bot/model/database"
	"bcjh-bot/scheduler"
	"bcjh-bot/util"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"fmt"
	"strings"
)

func StrategyQuery(c *scheduler.Context) {
	arg := strings.TrimSpace(c.PretreatedMessage)

	if arg == "" {
		_, _ = c.Reply(strategyHelp())
		return
	}

	if util.HasPrefixIn(arg, "新增", "添加") {
		if !global.IsSuperAdmin(c.GetSenderId()) {
			_, _ = c.Reply(e.PermissionDeniedNote)
			return
		}

		params := strings.Split(arg, "-")
		if len(params) < 3 {
			_, _ = c.Reply("参数格式错误!")
			return
		}
		keyword := params[1]
		value := params[2]
		note := createStrategy(keyword, value)
		if note != "" {
			_, _ = c.Reply(note)
		} else {
			msg := "添加攻略成功!\n"
			msg += fmt.Sprintf("关键词:%s\n", keyword)
			msg += fmt.Sprintf("内容:%s", value)
			_, _ = c.Reply(msg)
		}
		return
	}

	if util.HasPrefixIn(arg, "更新", "修改") {
		if !global.IsSuperAdmin(c.GetSenderId()) {
			_, _ = c.Reply(e.PermissionDeniedNote)
			return
		}

		params := strings.Split(arg, "-")
		if len(params) < 3 {
			_, _ = c.Reply("参数格式错误!")
			return
		}
		keyword := params[1]
		value := params[2]
		note := updateStrategy(keyword, value)
		if note != "" {
			_, _ = c.Reply(note)
		} else {
			msg := fmt.Sprintf("更新攻略「%s」成功!\n", keyword)
			msg += fmt.Sprintf("更新后内容:%s", value)
			_, _ = c.Reply(msg)
		}
		return
	}

	if util.HasPrefixIn(arg, "删除", "移除") {
		if !global.IsSuperAdmin(c.GetSenderId()) {
			_, _ = c.Reply(e.PermissionDeniedNote)
			return
		}

		params := strings.Split(arg, "-")
		if len(params) < 3 {
			_, _ = c.Reply("参数格式错误!")
			return
		}

		keyword := params[1]
		note := deleteStrategy(keyword)
		if note != "" {
			_, _ = c.Reply(note)
		} else {
			msg := fmt.Sprintf("删除攻略「%s」成功!", keyword)
			_, _ = c.Reply(msg)
		}
		return
	}

	strategies := make([]database.Strategy, 0)
	err := dao.DB.Where("keyword like ?", "%"+arg+"%").Find(&strategies)
	if err != nil {
		logger.Error("数据库查询出错!", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}
	if len(strategies) == 0 {
		logger.Debug("未找到相关攻略!")
		_, _ = c.Reply("这个有点难，我还没学会呢")
	} else if len(strategies) == 1 {
		logger.Debug("找到一个攻略:", strategies[0].Value)
		_, _ = c.Reply(strategies[0].Value)
	} else {
		logger.Debug("找到多条攻略:", strategies)
		msg := "这些攻略你想看哪条呀?\n"
		for _, strategy := range strategies {
			if strategy.Keyword == arg {
				msg = strategy.Value
				break
			}
			msg += fmt.Sprintf("%s  ", strategy.Keyword)
		}
		_, _ = c.Reply(msg)
	}
}

func createStrategy(keyword string, value string) string {
	if keyword == "" || value == "" {
		return "请填写完整关键词和内容"
	}

	strategy := new(database.Strategy)
	has, err := dao.DB.Where("keyword = ?", keyword).Get(strategy)
	if err != nil {
		return e.SystemErrorNote
	}
	if has {
		return "与已有关键词重复了哦"
	}

	strategy = new(database.Strategy)
	strategy.Keyword = keyword
	strategy.Value = value
	_, err = dao.DB.Insert(strategy)
	if err != nil {
		logger.Error("数据库新增攻略出错", err)
		return e.SystemErrorNote
	}
	return ""
}

func updateStrategy(keyword string, value string) string {
	if keyword == "" || value == "" {
		return "请填写完整关键词和内容"
	}

	strategy := new(database.Strategy)
	has, err := dao.DB.Where("keyword = ?", keyword).Get(strategy)
	if err != nil {
		return e.SystemErrorNote
	}
	if !has {
		return "关键词不存在!"
	}
	strategy.Value = value
	_, err = dao.DB.Where("id = ?", strategy.Id).Cols("value").Update(strategy)
	if err != nil {
		logger.Error("数据库更新攻略出错", err)
		return e.SystemErrorNote
	}
	return ""
}

func deleteStrategy(keyword string) string {
	if keyword == "" {
		return "请填写想要删除的关键词"
	}

	strategy := new(database.Strategy)
	affected, err := dao.DB.Where("keyword = ?", keyword).Delete(strategy)
	if err != nil {
		return e.SystemErrorNote
	}
	if affected == 0 {
		return "删除失败！未找到符合要求的关键词"
	}
	return ""
}

// 把限时任务单独拎出来做快捷查询
func TimeLimitingTaskQuery(c *scheduler.Context) {
	logger.Info("限时任务攻略查询")

	strategies := new(database.Strategy)
	has, err := dao.DB.Where("keyword = ?", "限时任务").Get(strategies)
	if err != nil {
		logger.Error("数据库查询出错!", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}
	if !has {
		_, _ = c.Reply("暂无限时任务攻略哦~")
	} else {
		logger.Debug("找到一个攻略:", strategies.Value)
		_, _ = c.Reply(strategies.Value)
	}
}
