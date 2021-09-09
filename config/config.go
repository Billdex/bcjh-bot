package config

import (
	"gopkg.in/ini.v1"
)

type serverConfig struct {
	Port int `ini:"port"`
}

type botConfig struct {
	PrivateMsgLen int `ini:"private_msg_len"`
	GroupMsgLen   int `ini:"group_msg_len"`
}

type dbConfig struct {
	Host     string `ini:"host"`
	Database string `ini:"database"`
	User     string `ini:"user"`
	Password string `ini:"password"`
}

type resourceConfig struct {
	Image string `ini:"image"`
	Font  string `ini:"font"`
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
	cfg, err := ini.Load("./config/app.ini")
	if nil != err {
		return err
	}

	AppConfig = &appConfig{
		Server: serverConfig{
			Port: 5800,
		},
		Bot: botConfig{
			PrivateMsgLen: 20,
			GroupMsgLen:   10,
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
	err = cfg.MapTo(AppConfig)
	if nil != err {
		return err
	}

	return nil
}
