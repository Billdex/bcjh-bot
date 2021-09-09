package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/config"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util/logger"
	"fmt"
)

// 图鉴网快捷访问
func GalleryWebsite(c *onebot.Context, args []string) {
	logger.Info("查询图鉴网")

	var msg string
	foodgame := "https://foodgame.gitee.io/"
	bcjh := "https://bcjh.gitee.io/"

	msg += fmt.Sprintf("L图鉴网: %s\n", foodgame)
	msg += fmt.Sprintf("白菜菊花手机图鉴网: %s", bcjh)

	err := bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}

// 游戏术语
func TermInfo(c *onebot.Context, args []string) {
	logger.Info("术语信息查询")
	//msg := termHelp()
	//err := bot.SendMessage(c, msg)
	//if err != nil {
	//	logger.Error("发送信息失败!", err)
	//}
}

// 白菜菊花App下载
func BCJHAppDownload(c *onebot.Context, args []string) {
	logger.Info("白菜菊花app下载")

	imgPath := config.AppConfig.Resource.Image + "/白菜菊花.jpg"
	var msg string
	msg += fmt.Sprintf("密码: bcjh\n")
	msg += bot.GetCQImage(imgPath, "file")

	err := bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}

// 计算器，不用第二遍解释
func Calculator(c *onebot.Context, args []string) {
	logger.Info("计算器查询")
	bcjh := "https://bcjh.gitee.io/"
	imgPath := config.AppConfig.Resource.Image + "/白菜菊花.jpg"
	var msg string
	msg += fmt.Sprintf("网页版计算器在白菜菊花图鉴网:%s\n", bcjh)
	msg += fmt.Sprintf("安卓用户支持使用白菜菊花app，扫描下图二维码下载，密码bcjh %s", bot.GetCQImage(imgPath, "file"))
	err := bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}

// 禁言转盘
func BanRandomGif(c *onebot.Context, arg []string) {
	logger.Info("禁言转盘")

	imgPath := config.AppConfig.Resource.Image + "/禁言转盘.gif"
	var msg string
	msg += fmt.Sprintf("客官这张图颜色够丰富吗?")
	msg += bot.GetCQImage(imgPath, "file")

	err := bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送消息失败!", err)
	}
}
