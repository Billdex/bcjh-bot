package dao

import (
	"bcjh-bot/config"
	"bcjh-bot/model/database"
	_ "embed"
	"fmt"
	"strings"
)

//go:embed sql/guest.sql
var guestSql string

//go:embed sql/laboratory.sql
var laboratorySql string

//go:embed sql/tarot.sql
var tarotSql string

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
	tableMap := map[string]string{
		database.Guest{}.TableName():      guestSql,
		database.Laboratory{}.TableName(): laboratorySql,
		database.Tarot{}.TableName():      tarotSql,
	}
	for tableName, sql := range tableMap {
		if total, err := DB.Table(tableName).Count(); err != nil {
			return err
		} else if total > 0 {
			continue
		} else {
			_, err = DB.Exec(sql)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
