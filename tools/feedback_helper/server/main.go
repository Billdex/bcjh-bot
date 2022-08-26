package main

import (
	"bcjh-bot/util/logger"
	"flag"
	"fmt"
)

func main() {
	cfgPath := flag.String("cfg", "config.ini", "配置文件路径")
	flag.Parse()

	// 初始化配置文件
	err := LoadConfig(*cfgPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 初始化logger
	err = logger.InitLog(logger.EncodeStyleConsole, cfg.LogPath, cfg.LogLevel)
	if err != nil {
		fmt.Println("初始化日志组件出错！", err)
		return
	}
	defer logger.Sync()

	// 初始化数据库
	err = InitDAO(cfg.MysqlDSN)
	if err != nil {
		fmt.Println("初始化数据库连接出错！", err)
		return
	}

	// 初始化 router
	r := InitRouter()
	// 启动 http 服务
	if err = r.Run(":" + cfg.Port); err != nil {
		logger.Error("启动 http 服务失败", err)
	}

}
