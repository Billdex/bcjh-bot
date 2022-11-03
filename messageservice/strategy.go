package messageservice

import (
	"bcjh-bot/dao"
	"bcjh-bot/model/database"
	"bcjh-bot/scheduler"
	"bcjh-bot/util"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"fmt"
	"strings"
)

func StrategyQuery(c *scheduler.Context) {
	arg := c.PretreatedMessage

	if util.HasPrefixIn(arg, "新增", "添加") {
		if !dao.IsSuperAdmin(c.GetSenderId()) {
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
		err := dao.CreateStrategy(keyword, value)
		if err != nil {
			_, _ = c.Reply(err.Error())
		} else {
			msg := fmt.Sprintf("添加攻略「%s」成功!\n", keyword)
			msg += fmt.Sprintf("内容:%s", value)
			_, _ = c.Reply(msg)
		}
		return
	}
	if util.HasPrefixIn(arg, "更新", "修改") {
		if !dao.IsSuperAdmin(c.GetSenderId()) {
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
		err := dao.UpdateStrategy(keyword, value)
		if err != nil {
			_, _ = c.Reply(err.Error())
			return
		} else {
			msg := fmt.Sprintf("更新攻略「%s」成功!\n", keyword)
			msg += fmt.Sprintf("更新后内容:%s", value)
			_, _ = c.Reply(msg)
		}
		return
	}
	if util.HasPrefixIn(arg, "删除", "移除") {
		if !dao.IsSuperAdmin(c.GetSenderId()) {
			_, _ = c.Reply(e.PermissionDeniedNote)
			return
		}
		params := strings.Split(arg, "-")
		if len(params) < 2 {
			_, _ = c.Reply("参数格式错误!")
			return
		}
		keyword := params[1]
		err := dao.DeleteStrategyByKeyword(keyword)
		if err != nil {
			_, _ = c.Reply(err.Error())
			return
		}

		_, _ = c.Reply(fmt.Sprintf("删除攻略「%s」成功!", keyword))
		return
	}

	keywords, err := dao.LoadStrategyKeywords()
	if err != nil {
		logger.Error("获取攻略关键词列表失败 %v", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}
	var matchList []string
	for _, keyword := range keywords {
		if strings.Contains(keyword, arg) {
			matchList = append(matchList, keyword)
		}
		if keyword == arg {
			matchList = []string{keyword}
			break
		}
	}

	if len(matchList) == 0 {
		_, _ = c.Reply("这个有点难，我还没学会呢")
	} else if len(matchList) == 1 {
		result, err := dao.GetStrategyByKeyword(matchList[0])
		if err != nil {
			_, _ = c.Reply(e.SystemErrorNote)
			return
		}
		_, _ = c.Reply(result.Value)
	} else {
		msg := "这些攻略你想看哪条呀?\n"
		msg += strings.Join(matchList, " ")
		_, _ = c.Reply(msg)
	}

	return
}

// TimeLimitingTaskQuery 把限时任务单独拎出来做快捷查询
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
