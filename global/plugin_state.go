package global

import (
	"bcjh-bot/model/database"
	"bcjh-bot/util/logger"
	"fmt"
	"sync"
)

const (
	pluginStateMapKey = "%d_%s"
)

var pluginStateMap = make(map[string]bool)
var pluginStateMux sync.Mutex

var pluginList = map[string][]string{
	// 管理功能
	"公告":  {},
	"防撤回": {},

	// 查询功能
	"帮助":   {"功能", "说明", "指引", "使用说明"},
	"反馈":   {"建议"},
	"厨师":   {"厨子"},
	"菜谱":   {"食谱", "菜单"},
	"厨具":   {"装备", "道具"},
	"食材":   {"材料"},
	"贵客":   {"稀有客人", "贵宾", "客人", "宾客", "稀客"},
	"符文":   {"礼物"},
	"调料":   {},
	"任务":   {"主线", "支线"},
	"限时任务": {"限时攻略", "限时支线"},
	"攻略":   {},
	"碰瓷":   {"升阶贵客", "升级贵客"},
	"后厨":   {"合成"},
	"兑换码":  {"玉璧", "兑奖码"},
	"实验室":  {"研究"},

	// 娱乐功能
	"抽签":     {"占卜", "求签", "运势", "卜卦", "占卦"},
	"随机个人图鉴": {},

	// 提醒功能
	"厨神提醒": {},
}

var pluginAliasComparison = make(map[string]string)

func getPluginState(key string) (bool, bool) {
	pluginStateMux.Lock()
	defer pluginStateMux.Unlock()
	value, ok := pluginStateMap[key]
	return value, ok
}

func updatePluginState(key string, value bool) {
	pluginStateMux.Lock()
	defer pluginStateMux.Unlock()
	pluginStateMap[key] = value
}

func deletePluginState(key string) {
	pluginStateMux.Lock()
	defer pluginStateMux.Unlock()
	delete(pluginStateMap, key)
}

func initPluginAliasComparison() {
	for key, pluginAliasList := range pluginList {
		pluginAliasComparison[key] = key
		for _, alias := range pluginAliasList {
			pluginAliasComparison[alias] = key
		}
	}
}

func GetPluginName(name string) (string, bool) {
	value, ok := pluginAliasComparison[name]
	return value, ok
}

func GetPluginState(groupId int64, pluginName string, defaultState bool) (bool, error) {
	key := fmt.Sprintf(pluginStateMapKey, groupId, pluginName)
	if pluginON, ok := getPluginState(key); ok {
		return pluginON, nil
	} else {
		pluginState := database.PluginState{}
		has, err := database.DB.Where("group_id = ? and plugin_name = ?", groupId, pluginName).Get(&pluginState)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return false, err
		}
		if has {
			updatePluginState(key, pluginState.State)
			return pluginState.State, nil
		} else {
			_, err := database.DB.Insert(&database.PluginState{
				GroupId:    groupId,
				PluginName: pluginName,
				State:      defaultState,
			})
			if err != nil {
				logger.Error("数据库插入数据出错", err)
				return false, err
			}
			deletePluginState(key)
			return defaultState, nil
		}
	}
}

func SetPluginState(groupId int64, pluginName string, state bool) error {
	key := fmt.Sprintf(pluginStateMapKey, groupId, pluginName)
	defer deletePluginState(key)
	has, err := database.DB.Where("group_id = ? and plugin_name = ?", groupId, pluginName).Get(&database.PluginState{})
	if err != nil {
		logger.Error("查询数据库出错!", err)
		return err
	}
	if has {
		_, err := database.DB.Cols("state").Where("group_id = ? and plugin_name = ?", groupId, pluginName).Update(&database.PluginState{State: state})
		if err != nil {
			logger.Error("更新数据库出错!", err)
			return err
		}
	} else {
		_, err := database.DB.Insert(&database.PluginState{
			GroupId:    groupId,
			PluginName: pluginName,
			State:      state,
		})
		if err != nil {
			logger.Error("数据库插入数据出错", err)
			return err
		}
	}
	return nil
}
