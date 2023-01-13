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

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(8)

	return nil
}

// LimitPaginate 限制分页
func LimitPaginate(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	} else if page > 500 {
		page = 500
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return page, pageSize
}
