package dao

import (
	"bcjh-bot/model/database"
	"bcjh-bot/util/logger"
	"fmt"
)

const CacheKeyUserData = "user_data"

// FindUserDataWithUserId 根据 userId 查询用户数据
func FindUserDataWithUserId(userId int64) (database.UserData, error) {
	var data []database.UserData
	cacheKey := fmt.Sprintf("%s:%d", CacheKeyUserData, userId)
	err := SimpleFindDataWithCache(cacheKey, &data, func(dest interface{}) error {
		return DB.Where("qq = ?", userId).Find(dest)
	})
	if len(data) > 0 {
		return data[0], err
	}
	return database.UserData{}, err
}

// SetUserData 设置或更新用户数据
func SetUserData(data database.UserData) error {
	key := fmt.Sprintf("%s:%d", CacheKeyUserData, data.QQ)
	defer Cache.Delete(key)
	has, err := DB.Where("qq = ?", data.QQ).Get(&database.UserData{})
	if err != nil {
		logger.Error("查询数据库出错!", err)
		return err
	}
	if has {
		_, err := DB.Cols("bcjh_id", "data", "user").Where("qq = ?", data.QQ).Update(&database.UserData{
			User:   data.User,
			BcjhID: data.BcjhID,
			Data:   data.Data,
		})
		if err != nil {
			logger.Error("更新数据库出错!", err)
			return err
		}
	} else {
		_, err := DB.Insert(&data)
		if err != nil {
			logger.Error("数据库插入数据出错", err)
			return err
		}
	}
	return nil
}
