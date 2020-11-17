package main

import (
	"bcjh-bot/config"
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

	s := server.NewServer()

	s.Run(":" + strconv.Itoa(config.AppConfig.Server.Port))

}
