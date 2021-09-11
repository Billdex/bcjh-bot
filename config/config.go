package config

import (
	"bcjh-bot/util"
	"fmt"
	"gopkg.in/ini.v1"
)

type serverConfig struct {
	Port int `ini:"port"`
}

type botConfig struct {
	PrivateMsgMaxLen  int `ini:"private_msg_max_len"`
	GroupMsgMaxLen    int `ini:"group_msg_max_len"`
	ExchangeMsgMaxLen int `ini:"exchange_msg_max_len"`
}

type dbConfig struct {
	Host     string `ini:"host"`
	Database string `ini:"database"`
	User     string `ini:"user"`
	Password string `ini:"password"`
}

type resourceConfig struct {
	Image    string `ini:"image"`
	Font     string `ini:"font"`
	Shortcut string `ini:"shortcut"`
	Sql      string `ini:"sql"`
}

type logConfig struct {
	Style   string `ini:"style"`
	Level   string `ini:"level"`
	OutPath string `ini:"out_path"`
}

type appConfig struct {
	Server   serverConfig   `ini:"server"`
	Bot      botConfig      `ini:"bot"`
	DB       dbConfig       `ini:"database"`
	Resource resourceConfig `ini:"resource"`
	Log      logConfig      `ini:"log"`
}

var AppConfig *appConfig

func InitConfig() error {
	AppConfig = &appConfig{
		Server: serverConfig{
			Port: 5800,
		},
		Bot: botConfig{
			PrivateMsgMaxLen: 20,
			GroupMsgMaxLen:   10,
		},
		DB: dbConfig{
			Host:     "127.0.0.1:3306",
			Database: "bcjh",
			User:     "root",
			Password: "",
		},
		Resource: resourceConfig{
			Image: "/home/bcjh-bot/resource/image/",
			Font:  "/home/bcjh-bot/resource/font",
		},
		Log: logConfig{
			Style:   "CONSOLE",
			Level:   "INFO",
			OutPath: "./logs/bcjh-bot.log",
		},
	}
	path := "./config/app.ini"
	has, err := util.PathExists(path)
	if !has {
		err := initDefaultConfig(path)
		if err != nil {
			return fmt.Errorf("未找到配置文件, 生成默认配置文件出错! %s", err)
		}
		return fmt.Errorf("未找到配置文件, 已生成默认配置文件")
	}
	cfg, err := ini.Load(path)
	if nil != err {
		return fmt.Errorf("加载配置文件出错! %s", err)
	}

	err = cfg.MapTo(AppConfig)
	if nil != err {
		return err
	}

	return nil
}

func initDefaultConfig(path string) error {
	defaultConfig := &appConfig{
		Server: serverConfig{
			Port: 5800,
		},
		Bot: botConfig{
			PrivateMsgMaxLen:  20,
			GroupMsgMaxLen:    10,
			ExchangeMsgMaxLen: 3,
		},
		DB: dbConfig{
			Host:     "127.0.0.1:3306",
			Database: "bcjh",
			User:     "root",
			Password: "",
		},
		Resource: resourceConfig{
			Image: "/home/bcjh-bot/resource/image/",
			Font:  "/home/bcjh-bot/resource/font",
		},
		Log: logConfig{
			Style:   "CONSOLE",
			Level:   "INFO",
			OutPath: "./logs/bcjh-bot.log",
		},
	}
	cfg := ini.Empty()
	err := ini.ReflectFrom(cfg, defaultConfig)
	if err != nil {
		return err
	}
	err = cfg.SaveTo(path)
	if err != nil {
		return err
	}
	return nil
}
