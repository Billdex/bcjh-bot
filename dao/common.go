package dao

import (
	"bcjh-bot/util/logger"
	"encoding/json"
)

// SimpleFindDataWithCache 查询数据并使用 cache 通用封装
func SimpleFindDataWithCache(key string, result interface{}) error {
	data, err := Cache.Get(key)
	if err == nil {
		err = json.Unmarshal(data, result)
		if err != nil {
			logger.Errorf("json 反序列化 cache %s 结果数据出错 %+v", key, err)
		} else {
			return nil
		}
	}
	err = DB.Find(result)
	if err != nil {
		return err
	}
	b, _ := json.Marshal(result)
	_ = Cache.Set(key, b)

	return nil
}
