package dao

import (
	"bcjh-bot/model/database"
	"bcjh-bot/util/logger"
	"fmt"
)

const (
	CacheKeyPluginState = "plugin_state_%d_%s"
)

var pluginList = map[string][]string{
	// 管理功能
	"公告": {},

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

// GetPluginState 获取某个群的某个功能启用状态
func GetPluginState(groupId int64, pluginName string, defaultState bool) (bool, error) {
	key := fmt.Sprintf(CacheKeyPluginState, groupId, pluginName)
	var state bool
	err := SimpleFindDataWithCache(key, &state, func(dest interface{}) error {
		var pluginState database.PluginState
		has, err := DB.Where("group_id = ? and plugin_name = ?", groupId, pluginName).Get(&pluginState)
		if has {
			*dest.(*bool) = pluginState.State
		} else {
			*dest.(*bool) = defaultState
		}
		return err
	})
	return state, err
}

// SetPluginState 设置某个群的某个功能启用状态
func SetPluginState(groupId int64, pluginName string, state bool) error {
	key := fmt.Sprintf(CacheKeyPluginState, groupId, pluginName)
	defer Cache.Delete(key)
	has, err := DB.Where("group_id = ? and plugin_name = ?", groupId, pluginName).Get(&database.PluginState{})
	if err != nil {
		logger.Error("查询数据库出错!", err)
		return err
	}
	if has {
		_, err := DB.Cols("state").Where("group_id = ? and plugin_name = ?", groupId, pluginName).Update(&database.PluginState{State: state})
		if err != nil {
			logger.Error("更新数据库出错!", err)
			return err
		}
	} else {
		_, err := DB.Insert(&database.PluginState{
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
