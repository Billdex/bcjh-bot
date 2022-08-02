package main

import (
	"bcjh-bot/config"
	"bcjh-bot/crontab"
	"bcjh-bot/dao"
	"bcjh-bot/messageservice"
	"bcjh-bot/noticeservice"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util/logger"
	"fmt"
	"strconv"
)

func main() {
	// 初始化配置文件
	err := config.InitConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("配置文件加载完毕")

	// 初始化logger
	err = logger.InitLog(config.AppConfig.Log.Style, config.AppConfig.Log.OutPath, config.AppConfig.Log.Level)
	if err != nil {
		fmt.Println("初始化logger出错！", err)
		return
	}
	defer logger.Sync()
	logger.Info("初始化logger完毕")

	// 初始化数据库引擎
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=Local",
		config.AppConfig.DB.User,
		config.AppConfig.DB.Password,
		config.AppConfig.DB.Host,
		config.AppConfig.DB.Database,
	)

	err = dao.InitDatabase(connStr)
	if err != nil {
		logger.Error("初始化数据库配置出错!", err)
		return
	}
	logger.Info("初始化数据库引擎完毕")

	// 注册插件与启动服务
	handler := &onebot.Handler{}
	s := scheduler.New()
	// 消息处理
	messageservice.Register(s)
	// 时间处理
	noticeservice.Register(handler)

	// 定时任务
	crontab.Register(s)
	port := strconv.Itoa(config.AppConfig.Server.Port)
	_ = s.Serve(":"+port, "/", handler)
}
