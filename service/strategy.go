package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/model/database"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"fmt"
	"strings"
)

func StrategyQuery(c *onebot.Context, args []string) {
	logger.Info("攻略查询, 参数:", args)

	if len(args) == 0 {
		//err := bot.SendMessage(c, strategyHelp())
		//if err != nil {
		//	logger.Error("发送信息失败!", err)
		//}
		return
	}

	if util.HasPrefixIn(args[0], "新增", "添加") {
		has, err := database.DB.Where("qq = ?", c.Sender.UserId).Exist(&database.Admin{})
		if err != nil {
			logger.Error("查询数据库出错", err)
			_ = bot.SendMessage(c, util.SystemErrorNote)
			return
		}
		if !has {
			_ = bot.SendMessage(c, util.PermissionDeniedNote)
			return
		}

		params := strings.Split(args[0], util.ArgsConnectCharacter)
		if len(params) < 3 {
			err := bot.SendMessage(c, "参数格式错误!")
			if err != nil {
				logger.Error("发送信息失败!", err)
			}
			return
		}
		keyword := params[1]
		value := params[2]
		note := createStrategy(keyword, value)
		if note != "" {
			err := bot.SendMessage(c, note)
			if err != nil {
				logger.Error("发送信息失败!", err)
			}
		} else {
			msg := "添加攻略成功!\n"
			msg += fmt.Sprintf("关键词:%s\n", keyword)
			msg += fmt.Sprintf("内容:%s", value)
			err := bot.SendMessage(c, msg)
			if err != nil {
				logger.Error("发送信息失败!", err)
			}
		}
		return
	}

	if util.HasPrefixIn(args[0], "更新", "修改") {
		has, err := database.DB.Where("qq = ?", c.Sender.UserId).Exist(&database.Admin{})
		if err != nil {
			logger.Error("查询数据库出错", err)
			_ = bot.SendMessage(c, util.SystemErrorNote)
			return
		}
		if !has {
			_ = bot.SendMessage(c, util.PermissionDeniedNote)
			return
		}

		params := strings.Split(args[0], util.ArgsConnectCharacter)
		if len(params) < 3 {
			err := bot.SendMessage(c, "参数格式错误!")
			if err != nil {
				logger.Error("发送信息失败!", err)
			}
			return
		}
		keyword := params[1]
		value := params[2]
		note := updateStrategy(keyword, value)
		if note != "" {
			err := bot.SendMessage(c, note)
			if err != nil {
				logger.Error("发送信息失败!", err)
			}
		} else {
			msg := fmt.Sprintf("更新攻略「%s」成功!\n", keyword)
			msg += fmt.Sprintf("更新后内容:%s", value)
			err := bot.SendMessage(c, msg)
			if err != nil {
				logger.Error("发送信息失败!", err)
			}
		}
		return
	}

	if util.HasPrefixIn(args[0], "删除", "移除") {
		has, err := database.DB.Where("qq = ?", c.Sender.UserId).Exist(&database.Admin{})
		if err != nil {
			logger.Error("查询数据库出错", err)
			_ = bot.SendMessage(c, util.SystemErrorNote)
			return
		}
		if !has {
			_ = bot.SendMessage(c, util.PermissionDeniedNote)
			return
		}

		params := strings.Split(args[0], util.ArgsConnectCharacter)
		if len(params) < 2 {
			err := bot.SendMessage(c, "参数格式错误!")
			if err != nil {
				logger.Error("发送信息失败!", err)
			}
			return
		}
		keyword := params[1]
		note := deleteStrategy(keyword)
		if note != "" {
			err := bot.SendMessage(c, note)
			if err != nil {
				logger.Error("发送信息失败!", err)
			}
		} else {
			msg := fmt.Sprintf("删除攻略「%s」成功!", keyword)
			err := bot.SendMessage(c, msg)
			if err != nil {
				logger.Error("发送信息失败!", err)
			}
		}
		return
	}

	strategies := make([]database.Strategy, 0)
	err := database.DB.Where("keyword like ?", "%"+args[0]+"%").Find(&strategies)
	if err != nil {
		logger.Error("数据库查询出错!", err)
		err := bot.SendMessage(c, util.SystemErrorNote)
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
		return
	}
	if len(strategies) == 0 {
		logger.Debug("未找到相关攻略!")
		err := bot.SendMessage(c, "这个有点难，我还没学会呢")
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
	} else if len(strategies) == 1 {
		logger.Debug("找到一个攻略:", strategies[0].Value)
		err := bot.SendMessage(c, strategies[0].Value)
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
	} else {
		logger.Debug("找到多条攻略:", strategies)
		msg := "这些攻略你想看哪条呀?\n"
		for _, strategy := range strategies {
			msg += fmt.Sprintf("%s  ", strategy.Keyword)
		}
		err := bot.SendMessage(c, msg)
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
	}
}

func createStrategy(keyword string, value string) string {
	if keyword == "" || value == "" {
		return "请填写完整关键词和内容"
	}

	strategy := new(database.Strategy)
	has, err := database.DB.Where("keyword = ?", keyword).Get(strategy)
	if err != nil {
		return util.SystemErrorNote
	}
	if has {
		return "与已有关键词重复了哦"
	}

	strategy = new(database.Strategy)
	strategy.Keyword = keyword
	strategy.Value = value
	_, err = database.DB.Insert(strategy)
	if err != nil {
		return util.SystemErrorNote
	}
	return ""
}

func updateStrategy(keyword string, value string) string {
	if keyword == "" || value == "" {
		return "请填写完整关键词和内容"
	}

	strategy := new(database.Strategy)
	has, err := database.DB.Where("keyword = ?", keyword).Get(strategy)
	if err != nil {
		return util.SystemErrorNote
	}
	if !has {
		return "关键词不存在!"
	}
	strategy.Value = value
	_, err = database.DB.Where("id = ?", strategy.Id).Cols("value").Update(strategy)
	if err != nil {
		return util.SystemErrorNote
	}
	return ""
}

func deleteStrategy(keyword string) string {
	if keyword == "" {
		return "请填写想要删除的关键词"
	}

	strategy := new(database.Strategy)
	affected, err := database.DB.Where("keyword = ?", keyword).Delete(strategy)
	if err != nil {
		return util.SystemErrorNote
	}
	if affected == 0 {
		return "删除失败！未找到符合要求的关键词"
	}
	return ""
}

// 把限时任务单独拎出来做快捷查询
func TimeLimitingQuestQuery(c *onebot.Context, args []string) {
	logger.Info("限时任务攻略查询")

	strategies := new(database.Strategy)
	has, err := database.DB.Where("keyword = ?", "限时任务").Get(strategies)
	if err != nil {
		logger.Error("数据库查询出错!", err)
		err := bot.SendMessage(c, util.SystemErrorNote)
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
		return
	}
	if !has {
		err := bot.SendMessage(c, "暂无限时任务攻略哦~")
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
	} else {
		logger.Debug("找到一个攻略:", strategies.Value)
		err := bot.SendMessage(c, strategies.Value)
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
	}
}
