package model

import (
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var DB *xorm.Engine

func InitDatabase(connString string) error {
	db, err := xorm.NewEngine("mysql", connString)
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	DB = db
	Migration()
	return nil
}
