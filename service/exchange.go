package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/model/database"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"strconv"
	"strings"
)

func ExchangeQuery(c *onebot.Context, args []string) {
	logger.Info("兑换码查询:", args)

	arg := strings.Join(args, "")
	num, err := strconv.Atoi(arg)
	if err == nil {
		// 参数为数字则查询最新的n条
		exchangeCodes := make([]database.Exchange, 0)
		if num > util.MaxQueryListLength/2 {
			num = util.MaxQueryListLength / 2
		}
		err = database.DB.Limit(num).Desc("update_time").Find(&exchangeCodes)
		if err != nil {
			logger.Error("数据库查询出错!", err)
			_ = bot.SendMessage(c, util.SystemErrorNote)
			return
		}
		if len(exchangeCodes) == 0 {
			_ = bot.SendMessage(c, "没有查询到兑换码记录")
			return
		}
		msg := ""
		for _, exchangeCode := range exchangeCodes {
			msg += exchangeCode.Content
		}
		_ = bot.SendMessage(c, msg)
		return
	} else {
		// 管理员可操作新增与更新
		if prefix, has := util.WhatPrefixIn(arg, "新增", "添加"); has {
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
			content := strings.ReplaceAll(arg, prefix, "")
			exchange := new(database.Exchange)
			exchange.Content = content
			_, err = database.DB.Insert(exchange)
			if err != nil {
				logger.Error("数据库插入出错!", err)
				_ = bot.SendMessage(c, util.SystemErrorNote)
				return
			}
			_ = bot.SendMessage(c, "新增兑换码信息成功!"+content)
			return
		} else if prefix, has := util.WhatPrefixIn(arg, "更新", "修改"); has {
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
			exchangeCodes := make([]database.Exchange, 0)
			err = database.DB.Limit(1).Desc("update_time").Find(&exchangeCodes)
			if err != nil {
				logger.Error("数据库查询出错!", err)
				_ = bot.SendMessage(c, util.SystemErrorNote)
				return
			}
			if len(exchangeCodes) == 0 {
				_ = bot.SendMessage(c, "没有可更新的兑换码")
				return
			}
			exchangeCode := exchangeCodes[0]
			exchangeCode.Content = strings.ReplaceAll(arg, prefix, "")
			_, err = database.DB.Where("id = ?", exchangeCode.Id).Update(&exchangeCode)
			if err != nil {
				logger.Error("数据库更新出错!", err)
				_ = bot.SendMessage(c, util.SystemErrorNote)
				return
			}
			_ = bot.SendMessage(c, "已更新最新兑换码内容!")
			return
		} else if has := util.HasPrefixIn(arg, "删除", "移除"); has {
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
			exchangeCodes := make([]database.Exchange, 0)
			err = database.DB.Limit(1).Desc("update_time").Find(&exchangeCodes)
			if err != nil {
				logger.Error("数据库查询出错!", err)
				_ = bot.SendMessage(c, util.SystemErrorNote)
				return
			}
			if len(exchangeCodes) == 0 {
				_ = bot.SendMessage(c, "没有可删除的兑换码")
				return
			}
			_, err = database.DB.Delete(&exchangeCodes[0])
			if err != nil {
				logger.Error("数据库删除出错!", err)
				_ = bot.SendMessage(c, util.SystemErrorNote)
				return
			}
			_ = bot.SendMessage(c, "已删除最新兑换码内容!")
			return
		}
	}
	// 其他情况均视为查询最新的一条
	exchangeCodes := make([]database.Exchange, 0)
	err = database.DB.Limit(1).Desc("update_time").Find(&exchangeCodes)
	if err != nil {
		logger.Error("数据库查询出错!", err)
		_ = bot.SendMessage(c, util.SystemErrorNote)
		return
	}
	if len(exchangeCodes) == 0 {
		_ = bot.SendMessage(c, "没有查询到兑换码记录")
		return
	}
	_ = bot.SendMessage(c, exchangeCodes[0].Content)
	return
}
