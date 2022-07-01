package messageservice

import (
	"bcjh-bot/config"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"fmt"
)

// 图鉴网快捷访问
func GalleryWebsite(c *scheduler.Context) {
	var msg string
	foodgame := "https://foodgame.github.io/"
	//bcjh := "https://bcjh.pages.dev/"
	bcjh := "https://www.bcjh.xyz/"

	msg += fmt.Sprintf("L图鉴网: %s\n", foodgame)
	//msg += fmt.Sprintf("白菜菊花手机图鉴网: %s\n", bcjh)
	msg += fmt.Sprintf("白菜菊花: %s", bcjh)

	_, _ = c.Reply(msg)
}

// 游戏术语
func TermInfo(c *scheduler.Context) {
	termImagePath := config.AppConfig.Resource.Shortcut + "/游戏术语.jpg"
	msg := onebot.GetCQImage(termImagePath, "file")
	_, _ = c.Reply(msg)
}

// 白菜菊花App下载
func BCJHAppDownload(c *scheduler.Context) {
	imgPath := config.AppConfig.Resource.Shortcut + "/白菜菊花.jpg"
	//imgPathXyz := config.AppConfig.Resource.Shortcut + "/白菜菊花xyz.jpg"
	var msg string
	msg += onebot.GetCQImage(imgPath, "file")
	//msg += onebot.GetCQImage(imgPathXyz, "file")

	_, _ = c.Reply(msg)
}

// 计算器，不用第二遍解释
func Calculator(c *scheduler.Context) {
	//bcjh := "https://bcjh.pages.dev/"
	bcjh := "https://www.bcjh.xyz/"
	//imgPathPage := config.AppConfig.Resource.Shortcut + "/白菜菊花page.jpg"
	//imgPathXyz := config.AppConfig.Resource.Shortcut + "/白菜菊花xyz.jpg"
	imgPath := config.AppConfig.Resource.Shortcut + "/白菜菊花.jpg"
	var msg string
	msg += fmt.Sprintf("网页版计算器在白菜菊花图鉴网:%s\n", bcjh)
	msg += fmt.Sprintf("安卓用户支持使用白菜菊花app，扫描下图二维码下载 %s",
		onebot.GetCQImage(imgPath, "file"))
	//onebot.GetCQImage(imgPathXyz, "file"))
	_, _ = c.Reply(msg)
}
