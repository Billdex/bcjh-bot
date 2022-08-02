package dao

import (
	"bcjh-bot/config"
	"bcjh-bot/model/database"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
	"strings"
	"xorm.io/xorm"
)

var DB *xorm.Engine

// DSN Data Source Name
// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]

// InitDatabase 初始化数据库连接
func InitDatabase(dsn string) error {
	var err error
	if config.AppConfig.DB.UseLocal {
		DB, err = xorm.NewEngine("sqlite", "./bot_data.db")
	} else {
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

// initDataImport 初始化数据导入
// 用于导入一些配置数据，以及图鉴网没有的数据
func initDataImport() error {
	// 初始化超级管理员信息
	if len(config.AppConfig.Bot.Admin) > 0 {
		var sql string
		if config.AppConfig.DB.UseLocal {
			sql = `INSERT OR IGNORE INTO admin (qq) values `
		} else {
			sql = `INSERT IGNORE INTO admin (qq) values `
		}
		qqValues := make([]string, 0, len(config.AppConfig.Bot.Admin))
		for _, adminId := range config.AppConfig.Bot.Admin {
			qqValues = append(qqValues, fmt.Sprintf("(%d)", adminId))
		}
		sql += strings.Join(qqValues, ",")
		_, err := DB.Exec(sql)
		if err != nil {
			return fmt.Errorf("配置导入超级管理员信息出错 %v", err)
		}
	}

	// 导入预配置数据
	tableMap := map[string]interface{}{
		"guest.sql":      database.Guest{},
		"laboratory.sql": database.Laboratory{},
	}
	for file, table := range tableMap {
		if total, err := DB.Count(table); err != nil {
			return err
		} else if total > 0 {
			continue
		} else {
			_, err = DB.ImportFile(fmt.Sprintf("%s/%s", config.AppConfig.Resource.Sql, file))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
