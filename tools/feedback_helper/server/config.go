package main

import (
	"fmt"
	"gopkg.in/ini.v1"
)

type Config struct {
	// 端口号
	Port string `ini:"port"`
	// 数据库连接 dsn
	MysqlDSN string `ini:"mysql_dsn"`
	// 日志配置
	LogLevel string `ini:"log_level"`
	LogPath  string `ini:"log_path"`
}

var cfg Config

func LoadConfig(path string) error {
	c, err := ini.Load(path)
	if err != nil {
		return fmt.Errorf("读取配置文件失败 %v", err)
	}
	err = c.MapTo(&cfg)
	if err != nil {
		return fmt.Errorf("获取配置数据失败 %v", err)
	}
	return nil
}
