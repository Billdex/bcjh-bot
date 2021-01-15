package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/model/database"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"strings"
)

func PublicNotice(c *onebot.Context, args []string) {
	logger.Info("发布公告:", args)

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

	groupList, err := bot.GetGroupList()
	if err != nil {
		logger.Error("获取群组列表失败!", err)
		_ = bot.SendMessage(c, "获取群组列表失败!")
		return
	}

	logger.Debugf("获取到%d个群组:%v", len(groupList), groupList)

	groups := make([]int, 0)
	for _, group := range groupList {
		groups = append(groups, group.GroupId)
	}

	notice := strings.Join(args, "")
	logger.Info("发送公告信息：", notice)
	bot.SendMassGroupMsg(notice, groups)
}
