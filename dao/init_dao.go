package dao

import (
	"bcjh-bot/config"
	"bcjh-bot/model/database"
	"fmt"
	"github.com/allegro/bigcache/v3"
	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
	"time"
	"xorm.io/xorm"
)

var DB *xorm.Engine
var Cache *bigcache.BigCache

func InitDao() error {
	if err := InitDatabase(); err != nil {
		return err
	}
	if err := InitCache(); err != nil {
		return err
	}
	initPluginAliasComparison()
	return nil
}

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
		return fmt.Errorf("连接数据库失败 %v", err)
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

// InitCache 初始化缓存
func InitCache() error {
	cacheCfg := bigcache.DefaultConfig(24 * time.Hour)
	cacheCfg.Shards = 16
	cacheCfg.HardMaxCacheSize = 64
	cache, err := bigcache.NewBigCache(cacheCfg)
	if err != nil {
		return fmt.Errorf("初始化缓存失败 %+v", err)
	}
	Cache = cache
	return nil
}
