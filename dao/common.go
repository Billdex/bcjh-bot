package dao

import (
	"bcjh-bot/util/logger"
	"bytes"
	"encoding/gob"
	"github.com/allegro/bigcache/v3"
)

// SimpleFindDataWithCache 查询数据并使用 cache 通用封装
func SimpleFindDataWithCache(key string, result interface{}, dbFunc func(dest interface{}) error) error {
	// 尝试从缓存读取数据
	data, reason, err := Cache.GetWithInfo(key)
	if err == nil {
		err = gob.NewDecoder(bytes.NewReader(data)).Decode(result)
		if err != nil {
			logger.Errorf("解码 cache %s 结果数据出错 %+v", key, err)
		} else {
			// 读取到数据且反序列化正常则直接返回
			return nil
		}
	}
	if err != bigcache.ErrEntryNotFound {
		logger.Infof("%s 读取缓存失败，原因：%+v, err: %+v", key, reason, err)
	}
	// 未读到缓存或反序列化失败则执行传入的 db 查询方法
	err = dbFunc(result)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	_ = gob.NewEncoder(&buf).Encode(result)
	err = Cache.Set(key, buf.Bytes())
	if err != nil {
		logger.Warnf("%s Cache Set Fail. err: %+v", key, err)
	}

	return nil
}
