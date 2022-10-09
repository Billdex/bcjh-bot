package dao

import (
	"bcjh-bot/model/database"
	"bcjh-bot/util/logger"
	"fmt"
	"time"
)

const CacheKeyUserAllowState = "allow_state_%d_%d"

// GetUserAllowState 获取用户在某个群的使用权限状态
func GetUserAllowState(userId, groupId int64) (bool, error) {
	key := fmt.Sprintf(CacheKeyUserAllowState, userId, groupId)
	var endTimestamp int64
	err := SimpleFindDataWithCache(key, &endTimestamp, func(dest interface{}) error {
		var banState database.BlackList
		_, err := DB.Cols("end_time").Where("qq = ? and group_id = ?", userId, groupId).Get(&banState)
		// 未查询到数据则说明用户未被禁用过, 可以直接写入零值
		*dest.(*int64) = banState.EndTime
		return err
	})
	if err != nil {
		return false, err
	}
	return endTimestamp < time.Now().Unix(), nil
}

// SetUserBanTime 设置用户禁用时间
func SetUserBanTime(userId int64, groupId int64, endTime int64) error {
	key := fmt.Sprintf(CacheKeyUserAllowState, userId, groupId)
	defer Cache.Delete(key)
	has, err := DB.Where("qq = ? and group_id = ?", userId, groupId).Get(&database.BlackList{})
	if err != nil {
		logger.Error("查询数据库出错!", err)
		return err
	}
	if has {
		_, err := DB.Cols("end_time", "end_datetime").Where("qq = ? and group_id = ?", userId, groupId).
			Update(&database.BlackList{
				EndTime:     endTime,
				EndDatetime: time.Unix(endTime, 0),
			})
		if err != nil {
			logger.Error("更新数据库出错!", err)
			return err
		}
	} else {
		_, err := DB.Insert(&database.BlackList{
			QQ:          userId,
			GroupId:     groupId,
			EndTime:     endTime,
			EndDatetime: time.Unix(endTime, 0),
		})
		if err != nil {
			logger.Error("数据库插入数据出错", err)
			return err
		}
	}
	return nil
}
