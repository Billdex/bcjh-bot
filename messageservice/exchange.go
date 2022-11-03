package messageservice

import (
	"bcjh-bot/config"
	"bcjh-bot/dao"
	"bcjh-bot/model/database"
	"bcjh-bot/scheduler"
	"bcjh-bot/util"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"strconv"
	"strings"
)

func ExchangeQuery(c *scheduler.Context) {
	arg := c.PretreatedMessage
	num, err := strconv.Atoi(arg)
	if err == nil {
		// 参数为数字则查询最新的n条
		exchangeCodes := make([]database.Exchange, 0)
		if num > config.AppConfig.Bot.ExchangeMsgMaxLen {
			num = config.AppConfig.Bot.ExchangeMsgMaxLen
		}
		err = dao.DB.Limit(num).Desc("update_time").Find(&exchangeCodes)
		if err != nil {
			logger.Error("数据库查询出错!", err)
			_, _ = c.Reply(e.SystemErrorNote)
			return
		}
		if len(exchangeCodes) == 0 {
			_, _ = c.Reply("没有查询到兑换码记录")
			return
		}
		msg := ""
		for _, exchangeCode := range exchangeCodes {
			msg += exchangeCode.Content
		}
		_, _ = c.Reply(msg)
		return
	} else {
		// 管理员可操作新增与更新
		if prefix, has := util.WhatPrefixIn(arg, "新增", "添加"); has {
			if !dao.IsSuperAdmin(c.GetSenderId()) {
				_, _ = c.Reply(e.PermissionDeniedNote)
				return
			}
			content := strings.ReplaceAll(arg, prefix, "")
			exchange := new(database.Exchange)
			exchange.Content = content
			_, err = dao.DB.Insert(exchange)
			if err != nil {
				logger.Error("数据库插入出错!", err)
				_, _ = c.Reply(e.SystemErrorNote)
				return
			}
			_, _ = c.Reply("新增兑换码信息成功!" + content)
			return
		} else if prefix, has := util.WhatPrefixIn(arg, "更新", "修改"); has {
			if !dao.IsSuperAdmin(c.GetSenderId()) {
				_, _ = c.Reply(e.PermissionDeniedNote)
				return
			}
			exchangeCodes := make([]database.Exchange, 0)
			err = dao.DB.Limit(1).Desc("update_time").Find(&exchangeCodes)
			if err != nil {
				logger.Error("数据库查询出错!", err)
				_, _ = c.Reply(e.SystemErrorNote)
				return
			}
			if len(exchangeCodes) == 0 {
				_, _ = c.Reply("没有可更新的兑换码")
				return
			}
			exchangeCode := exchangeCodes[0]
			exchangeCode.Content = strings.ReplaceAll(arg, prefix, "")
			_, err = dao.DB.Where("id = ?", exchangeCode.Id).Update(&exchangeCode)
			if err != nil {
				logger.Error("数据库更新出错!", err)
				_, _ = c.Reply(e.SystemErrorNote)
				return
			}
			_, _ = c.Reply("已更新最新兑换码内容!")
			return
		} else if has := util.HasPrefixIn(arg, "删除", "移除"); has {
			if !dao.IsSuperAdmin(c.GetSenderId()) {
				_, _ = c.Reply(e.PermissionDeniedNote)
				return
			}
			exchangeCodes := make([]database.Exchange, 0)
			err = dao.DB.Limit(1).Desc("update_time").Find(&exchangeCodes)
			if err != nil {
				logger.Error("数据库查询出错!", err)
				_, _ = c.Reply(e.SystemErrorNote)
				return
			}
			if len(exchangeCodes) == 0 {
				_, _ = c.Reply("没有可删除的兑换码")
				return
			}
			_, err = dao.DB.Delete(&exchangeCodes[0])
			if err != nil {
				logger.Error("数据库删除出错!", err)
				_, _ = c.Reply(e.SystemErrorNote)
				return
			}
			_, _ = c.Reply("已删除最新兑换码内容!")
			return
		}
	}
	// 其他情况均视为查询最新的一条
	exchangeCodes := make([]database.Exchange, 0)
	err = dao.DB.Limit(1).Desc("update_time").Find(&exchangeCodes)
	if err != nil {
		logger.Error("数据库查询出错!", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}
	if len(exchangeCodes) == 0 {
		_, _ = c.Reply("没有查询到兑换码记录")
		return
	}
	_, _ = c.Reply(exchangeCodes[0].Content)
	return
}
