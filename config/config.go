package config

import "gopkg.in/ini.v1"

type serverConfig struct {
	Port int `ini:"port"`
}

type cqhttpConfig struct {
	Host string `ini:"host"`
	Port int    `ini:"port"`
}

type dbConfig struct {
	Host      string `ini:"host"`
	Port      int    `ini:"port"`
	Database  string `ini:"database"`
	User      string `ini:"user"`
	Password  string `ini:"password"`
	ExportDir string `ini:"export_dir"`
}

type logConfig struct {
	Style string `ini:"style"`
	Level string `ini:"level"`
	File  string `ini:"file"`
}

type appConfig struct {
	Server   serverConfig `ini:"server"`
	CQHTTP   cqhttpConfig `ini:"cqhttp"`
	DBConfig dbConfig     `ini:"database"`
	Log      logConfig    `ini:"log"`
}

var AppConfig *appConfig

func InitConfig() error {
	cfg, err := ini.Load("./config/app.ini")
	if nil != err {
		return err
	}

	AppConfig = new(appConfig)
	err = cfg.MapTo(AppConfig)
	if nil != err {
		return err
	}

	return nil
}
