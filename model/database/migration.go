package database

import (
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

var tables = []interface{}{
	&Feedback{},
	&Admin{},
	&Chef{},
	&Equip{},
	&Recipe{},
	&Guest{},
	&Material{},
	&Skill{},
	&RecipeMaterial{},
	&Decoration{},
}

var migrations = []*migrate.Migration{
	{
		ID: "202011260341",
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
