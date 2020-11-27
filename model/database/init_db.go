package database

import (
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var DB *xorm.Engine

// DNS Data Source Name
// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
func InitDatabase(dsn string) error {
	var err error
	DB, err = xorm.NewEngine("mysql", dsn)
	if err != nil {
		return err
	}

	DB.SetMaxIdleConns(5)
	DB.SetMaxOpenConns(10)
	return Migration()
}
