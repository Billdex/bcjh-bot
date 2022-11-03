package dao

import (
	"bcjh-bot/model/database"
	"bcjh-bot/util/logger"
	"fmt"
)

const CacheKeyBotState = "bot_state_%d_%d"

// GetBotState 获取机器人账号在某个群的启用状态
func GetBotState(botId, groupId int64) (bool, error) {
	key := fmt.Sprintf(CacheKeyBotState, botId, groupId)
	var state bool
	err := SimpleFindDataWithCache(key, &state, func(dest interface{}) error {
		var botState database.BotState
		_, err := DB.Where("bot_id = ? and group_id = ?", botId, groupId).Get(&botState)
		// 未查询到结果 与 查询到了但状态为否 的清空都视为未启用
		*dest.(*bool) = botState.State
		return err
	})
	return state, err
}

// SetBotState 设置机器人账号在群内的启用状态
func SetBotState(botId int64, groupId int64, state bool) error {
	key := fmt.Sprintf(CacheKeyBotState, botId, groupId)
	defer Cache.Delete(key)
	has, err := DB.Where("bot_id = ? and group_id = ?", botId, groupId).Get(&database.BotState{})
	if err != nil {
		logger.Error("查询数据库出错!", err)
		return err
	}
	if has {
		_, err := DB.Cols("state").Where("bot_id = ? and group_id = ?", botId, groupId).Update(&database.BotState{State: state})
		if err != nil {
			logger.Error("更新数据库出错!", err)
			return err
		}
	} else {
		_, err := DB.Insert(&database.BotState{
			BotId:   botId,
			GroupId: groupId,
			State:   state,
		})
		if err != nil {
			logger.Error("数据库插入数据出错", err)
			return err
		}
	}
	return nil

}
