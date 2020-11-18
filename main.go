package main

import (
	"bcjh-bot/config"
	"bcjh-bot/logger"
	"bcjh-bot/server"
	"fmt"
	"strconv"
)

func main() {
	// 初始化配置信息
	err := config.InitConfig()
	if err != nil {
		fmt.Println("读取配置文件出错！", err)
		return
	}
	fmt.Println("配置文件加载完毕")

	err = logger.InitLog(config.AppConfig.Log.Style, config.AppConfig.Log.File, config.AppConfig.Log.Level)
	if err != nil {
		fmt.Println("初始化logger出错！", err)
		return
	}
	defer logger.Sync()
	logger.Info("初始化logger完毕")

	port := strconv.Itoa(config.AppConfig.Server.Port)
	err = server.Run(":" + port)
	if err != nil {
		fmt.Println("服务启动出错!", err)
		return
	}
}
