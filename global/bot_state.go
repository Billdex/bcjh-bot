package global

import (
	"bcjh-bot/model/database"
	"bcjh-bot/util/logger"
	"fmt"
	"sync"
)

const (
	botStateMapKey = "%d_%d"
)

var botStateMap = make(map[string]bool)
var botStateMux sync.Mutex

func getBotState(key string) (bool, bool) {
	botStateMux.Lock()
	defer botStateMux.Unlock()
	value, ok := botStateMap[key]
	return value, ok
}

func updateBotState(key string, value bool) {
	botStateMux.Lock()
	defer botStateMux.Unlock()
	botStateMap[key] = value
}

func deleteBotState(key string) {
	botStateMux.Lock()
	defer botStateMux.Unlock()
	delete(botStateMap, key)
}

func GetBotState(botId int64, groupId int64) (bool, error) {
	key := fmt.Sprintf(botStateMapKey, botId, groupId)
	if botON, ok := getBotState(key); ok {
		return botON, nil
	} else {
		botState := database.BotState{}
		has, err := database.DB.Where("bot_id = ? and group_id = ?", botId, groupId).Get(&botState)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return false, err
		}
		if has {
			updateBotState(key, botState.State)
			return botState.State, nil
		} else {
			_, err := database.DB.Insert(&database.BotState{
				BotId:   botId,
				GroupId: groupId,
				State:   false,
			})
			if err != nil {
				logger.Error("数据库插入数据出错", err)
				return false, err
			}
			deleteBotState(key)
			return false, nil
		}
	}
}

func SetBotState(botId int64, groupId int64, state bool) error {
	key := fmt.Sprintf(botStateMapKey, botId, groupId)
	defer deleteBotState(key)
	has, err := database.DB.Where("bot_id = ? and group_id = ?", botId, groupId).Get(&database.BotState{})
	if err != nil {
		logger.Error("查询数据库出错!", err)
		return err
	}
	if has {
		_, err := database.DB.Cols("state").Where("bot_id = ? and group_id = ?", botId, groupId).Update(&database.BotState{State: state})
		if err != nil {
			logger.Error("更新数据库出错!", err)
			return err
		}
	} else {
		_, err := database.DB.Insert(&database.BotState{
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
