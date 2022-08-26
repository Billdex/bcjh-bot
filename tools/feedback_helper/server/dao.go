package main

import (
	"fmt"
	"xorm.io/xorm"
)

var db *xorm.Engine

func InitDAO(dsn string) error {
	var err error
	db, err = xorm.NewEngine("mysql", dsn)

	if err != nil {
		return fmt.Errorf("连接数据库失败 %v", err)
	}

	db.SetMaxIdleConns(8)
	db.SetMaxOpenConns(16)

	return nil
}
