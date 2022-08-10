package global

import (
	"bcjh-bot/model/database"
	"bcjh-bot/util/logger"
	"fmt"
	"sync"
	"time"
)

const (
	blackListMapKey = "%d_%d"
)

var blackListMap = make(map[string]int64)
var blackListMux sync.Mutex

func getUserState(key string) (int64, bool) {
	blackListMux.Lock()
	defer blackListMux.Unlock()
	value, ok := blackListMap[key]
	return value, ok
}

func updateUserState(key string, value int64) {
	blackListMux.Lock()
	defer blackListMux.Unlock()
	blackListMap[key] = value
}

func deleteUserState(key string) {
	blackListMux.Lock()
	defer blackListMux.Unlock()
	delete(blackListMap, key)
}

func GetUserAllowState(userId int64, groupId int64) bool {
	now := time.Now().Unix()
	key := fmt.Sprintf(botStateMapKey, userId, groupId)
	if endTime, ok := getUserState(key); ok {
		if endTime < now {
			_, err := database.DB.Where("qq = ? and group_id = ?", userId, groupId).Delete(&database.BlackList{})
			if err != nil {
				logger.Error("数据库执行删除出错!", err)
			}
			deleteUserState(key)
			return true
		} else {
			return false
		}
	} else {
		userState := database.BlackList{}
		has, err := database.DB.Where("qq = ? and group_id = ?", userId, groupId).Get(&userState)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return true
		}
		if has {
			if userState.EndTime < now {
				_, err := database.DB.Where("qq = ? and group_id = ?", userId, groupId).Delete(&database.BlackList{})
				if err != nil {
					logger.Error("数据库执行删除出错!", err)
				}
				deleteUserState(key)
				return true
			} else {
				updateUserState(key, userState.EndTime)
				return false
			}
		} else {
			deleteUserState(key)
			return true
		}
	}
}

func PullUserBlackList(userId int64, groupId int64, endTime int64) error {
	key := fmt.Sprintf(blackListMapKey, userId, groupId)
	defer deleteUserState(key)
	has, err := database.DB.Where("qq = ? and group_id = ?", userId, groupId).Get(&database.BlackList{})
	if err != nil {
		logger.Error("查询数据库出错!", err)
		return err
	}
	if has {
		_, err := database.DB.Cols("end_time").Where("qq = ? and group_id = ?", userId, groupId).
			Update(&database.BlackList{EndTime: endTime})
		if err != nil {
			logger.Error("更新数据库出错!", err)
			return err
		}
	} else {
		_, err := database.DB.Insert(&database.BlackList{
			QQ:      userId,
			GroupId: groupId,
			EndTime: endTime,
		})
		if err != nil {
			logger.Error("数据库插入数据出错", err)
			return err
		}
	}
	return nil
}

func RemoveUserFromBlackList(userId int64, groupId int64) error {
	key := fmt.Sprintf(blackListMapKey, userId, groupId)
	defer deleteUserState(key)
	_, err := database.DB.Where("qq = ? and group_id = ?", userId, groupId).Delete(&database.BlackList{})
	if err != nil {
		logger.Error("删除数据库数据出错!", err)
		return err
	}
	return nil
}
