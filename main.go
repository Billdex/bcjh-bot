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
	"flag"
	"log"
	"strconv"
	"time"
)

func main() {
	// 加在启动参数
	cfgPath := flag.String("cfg", "config.ini", "配置文件路径")
	flag.Parse()

	// 初始化配置文件
	err := config.InitConfig(*cfgPath)
	if err != nil {
		log.Println(err)
		time.Sleep(5 * time.Second)
		return
	}
	log.Println("已加载配置文件")

	// 初始化logger
	err = logger.InitLog(logger.EncodeStyleConsole, config.AppConfig.Log.OutPath, config.AppConfig.Log.Level)
	if err != nil {
		log.Println("初始化日志组件出错！", err)
		time.Sleep(5 * time.Second)
		return
	}
	defer logger.Sync()
	log.Println("日志组件初始化完毕")

	// 初始化 Dao 层
	err = dao.InitDao()
	if err != nil {
		log.Println("初始化数据库或缓存出错!", err)
		time.Sleep(5 * time.Second)
		return
	}
	log.Println("初始化数据库与缓存完毕")

	// 注册插件与启动服务
	handler := &onebot.Handler{}
	s := scheduler.New()
	// 消息处理
	messageservice.Register(s)
	// 事件处理
	noticeservice.Register(handler)

	// 定时任务
	crontab.Register(s)
	port := strconv.Itoa(config.AppConfig.Server.Port)
	_ = s.Serve(":"+port, "/", handler)
}
