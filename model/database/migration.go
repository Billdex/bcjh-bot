package database

import (
	"time"
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

var tables = []interface{}{
	&Feedback{},
	&Admin{},
	&Chef{},
	&Equip{},
	&Recipe{},
	&GuestGift{},
	&Material{},
	&Skill{},
	&RecipeMaterial{},
	&Decoration{},
	&Condiment{},
}

var migrations = []*migrate.Migration{
	{
		ID: time.Now().Format("20060102150405"), // 以时间作为迁移 ID，在程序启动时会进行一次数据库结构的更新
		Migrate: func(engine *xorm.Engine) error {
			return engine.Sync2(tables...)
		},
		Rollback: func(engine *xorm.Engine) error {
			return engine.DropTables(tables...)
		},
	},
}

func Migration() error {
	m := migrate.New(DB, migrate.DefaultOptions, migrations)
	return m.Migrate()
}
