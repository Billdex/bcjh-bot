package dao

import (
	"bcjh-bot/config"
	"bcjh-bot/model/database"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

var DB *xorm.Engine

// DSN Data Source Name
// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]

// InitDatabase 初始化数据库连接
func InitDatabase() error {
	var err error
	if config.AppConfig.DB.UseLocal {
		DB, err = xorm.NewEngine("sqlite", "./bot_data.db")
	} else {
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=Local",
			config.AppConfig.DB.User,
			config.AppConfig.DB.Password,
			config.AppConfig.DB.Host,
			config.AppConfig.DB.Database,
		)
		DB, err = xorm.NewEngine("mysql", dsn)
	}
	if err != nil {
		return fmt.Errorf("初始化数据库连接失败 %v", err)
	}

	DB.SetMaxIdleConns(8)
	DB.SetMaxOpenConns(16)

	err = database.Migration(DB)
	if err != nil {
		return fmt.Errorf("执行 migrate 失败 %v", err)
	}
	err = initDataImport()
	if err != nil {
		return fmt.Errorf("初始化数据导入失败 %v", err)
	}
	return nil
}
