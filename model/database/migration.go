package database

import (
	"time"
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

var tables = []interface{}{
	&Admin{},
	&BlackList{},
	&BotState{},
	&Chef{},
	&Condiment{},
	&Decoration{},
	&Equip{},
	&Exchange{},
	&Feedback{},
	&Guest{},
	&GuestGift{},
	&Laboratory{},
	&Material{},
	&PluginState{},
	&Quest{},
	&Recipe{},
	&RecipeMaterial{},
	&Skill{},
	&Strategy{},
	&Tarot{},
	&UserData{},
}

var migrations = []*migrate.Migration{
	{
		ID: time.Now().Format("20060102150405"), // 以时间作为迁移 ID，在程序启动时会进行一次数据库结构的更新
		Migrate: func(engine *xorm.Engine) error {
			return engine.Sync(tables...)
		},
		Rollback: func(engine *xorm.Engine) error {
			return engine.DropTables(tables...)
		},
	},
}

func Migration(db *xorm.Engine) error {
	m := migrate.New(db, migrate.DefaultOptions, migrations)
	return m.Migrate()
}
