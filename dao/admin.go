package dao

import (
	"bcjh-bot/model/database"
	"bcjh-bot/util/logger"
)

// IsSuperAdmin 判断用户是否为超管
func IsSuperAdmin(qq int64) bool {
	has, err := DB.Where("qq = ?", qq).Exist(&database.Admin{})
	if err != nil {
		logger.Error("查询数据库出错:", err)
		return false
	}
	return has
}
