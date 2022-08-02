package config

import (
	"bcjh-bot/util"
	"fmt"
	"gopkg.in/ini.v1"
)

type serverConfig struct {
	Port int `ini:"port"`
}

// botConfig 机器人相关配置
type botConfig struct {
	PrivateMsgMaxLen  int     `ini:"private_msg_max_len"`  // 私聊信息最大长度，对于一些列表查询用此参数进行分页
	GroupMsgMaxLen    int     `ini:"group_msg_max_len"`    // 群聊消息最大长度
	ExchangeMsgMaxLen int     `ini:"exchange_msg_max_len"` // 每次最多可查询的兑换码长度
	Admin             []int64 `ini:"admin"`                // 超级管理员列表，填 QQ 号
}

// dbConfig 数据库相关的配置
type dbConfig struct {
	UseLocal bool   `int:"use_local"` // 是否使用本地储存，开启后将在程序目录创建一个 sqlite 数据库文件。如果不使用则读取下列配置连接 mysql
	Host     string `ini:"host"`      // 数据库 host，如果非域名默认端口，还需要带上端口号
	Database string `ini:"database"`  // 数据库 database 名
	User     string `ini:"user"`      // 登录 user
	Password string `ini:"password"`  // 登录密码
}

// resourceConfig 资源配置，主要是所需要的资源路径
type resourceConfig struct {
	Image    string `ini:"image"`    // 各项图片资源存放和生成用的路径
	Font     string `ini:"font"`     // 字体路径
	Shortcut string `ini:"shortcut"` // 部分快捷用于所需图片的存放路径
	Sql      string `ini:"sql"`      // 需要导入的预配置 sql 数据文件目录
}

// logConfig 日志配置
type logConfig struct {
	Style   string `ini:"style"`    // 日志风格，可用 console 或 json
	Level   string `ini:"level"`    // 日志级别
	OutPath string `ini:"out_path"` // 日志输出路径
}

// appConfig 应用一级设置，配置文件解析到该结构对象
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
			PrivateMsgMaxLen:  20,
			GroupMsgMaxLen:    10,
			ExchangeMsgMaxLen: 3,
		},
		DB: dbConfig{
			UseLocal: true,
			Host:     "127.0.0.1:3306",
			Database: "bcjh",
			User:     "bcjh",
			Password: "",
		},
		Resource: resourceConfig{
			Image:    "/home/bcjh-bot/resource/image/",
			Font:     "/home/bcjh-bot/resource/font/",
			Shortcut: "/home/bcjh-bot/resource/shortcut/",
			Sql:      "/home/bcjh-bot/resource/sql/",
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
			UseLocal: true,
			Host:     "127.0.0.1:3306",
			Database: "bcjh",
			User:     "root",
			Password: "",
		},
		Resource: resourceConfig{
			Image:    "/home/bcjh-bot/resource/image/",
			Font:     "/home/bcjh-bot/resource/font/",
			Shortcut: "/home/bcjh-bot/resource/shortcut/",
			Sql:      "/home/bcjh-bot/resource/sql/",
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
