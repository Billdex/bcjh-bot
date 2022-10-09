package dao

import (
	"bcjh-bot/util/logger"
	"encoding/json"
)

// SimpleFindDataWithCache 查询数据并使用 cache 通用封装
func SimpleFindDataWithCache(key string, result interface{}, dbFunc func(dest interface{}) error) error {
	// 尝试从缓存读取数据
	data, err := Cache.Get(key)
	if err == nil {
		err = json.Unmarshal(data, result)
		if err != nil {
			logger.Errorf("json 反序列化 cache %s 结果数据出错 %+v", key, err)
		} else {
			// 读取到数据且反序列化正常则直接返回
			return nil
		}
	}
	// 未读到缓存或反序列化失败则执行传入的 db 查询方法
	err = dbFunc(result)
	if err != nil {
		return err
	}
	b, _ := json.Marshal(result)
	_ = Cache.Set(key, b)

	return nil
}
