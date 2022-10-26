package dao

import (
	"bcjh-bot/model/database"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"errors"
	"fmt"
)

const CacheKeyStrategyKeywords = "strategy_keywords"
const CacheKeyStrategyData = "strategy_data_%s"

// LoadStrategyKeywords 加载可用的攻略数据关键词
func LoadStrategyKeywords() ([]string, error) {
	var strategies []string
	err := SimpleFindDataWithCache(CacheKeyStrategyKeywords, &strategies, func(dest interface{}) error {
		results := make([]database.Strategy, 0)
		err := DB.Cols("keyword").Find(&results)
		if err != nil {
			return err
		}
		keywords := make([]string, 0, len(results))
		for i := range results {
			keywords = append(keywords, results[i].Keyword)
		}
		*dest.(*[]string) = keywords
		return nil
	})
	return strategies, err
}

// GetStrategyByKeyword 查询攻略数据
func GetStrategyByKeyword(keyword string) (string, error) {
	var strategy string
	key := fmt.Sprintf(CacheKeyStrategyData, keyword)
	err := SimpleFindDataWithCache(key, &strategy, func(dest interface{}) error {
		var result database.Strategy
		_, err := DB.Where("keyword = ?", keyword).Get(&result)
		if err != nil {
			return err
		}
		*dest.(*string) = result.Value
		return nil
	})
	return strategy, err
}

// HasStrategyKeyword 判断某个攻略关键词是否存在
func HasStrategyKeyword(keyword string) bool {
	keywords, err := LoadStrategyKeywords()
	if err != nil {
		logger.Errorf("载入关键词数据列表出错 %v", err)
		return false
	}
	for i := range keywords {
		if keywords[i] == keyword {
			return true
		}
	}
	return false
}

func CreateStrategy(keyword string, value string) error {
	if keyword == "" || value == "" {
		return errors.New("未填写关键词或内容")
	}
	if HasStrategyKeyword(keyword) {
		return errors.New("攻略关键词已存在")
	}
	_, err := DB.Insert(&database.Strategy{
		Keyword: keyword,
		Value:   value,
	})
	if err != nil {
		logger.Errorf("创建攻略 %s 失败 %v", keyword, err)
		return errors.New(e.SystemErrorNote)
	}
	Cache.Delete(CacheKeyStrategyKeywords)
	return nil
}

func UpdateStrategy(keyword string, value string) error {
	if keyword == "" || value == "" {
		return errors.New("未填写关键词或内容")
	}
	if !HasStrategyKeyword(keyword) {
		return errors.New("攻略不存在，无法更新")
	}
	affected, err := DB.Where("keyword = ?", keyword).Update(&database.Strategy{
		Keyword: keyword,
		Value:   value,
	})
	if err != nil {
		logger.Errorf("更新攻略 %s 失败 %v", keyword, err)
		return errors.New(e.SystemErrorNote)
	}
	if affected == 0 {
		return errors.New("攻略不存在")
	}
	Cache.Delete(CacheKeyStrategyKeywords)
	Cache.Delete(fmt.Sprintf(CacheKeyStrategyData, keyword))
	return nil
}

func DeleteStrategyByKeyword(keyword string) error {
	if keyword == "" {
		return errors.New("未填写要移除的攻略关键词")
	}
	if !HasStrategyKeyword(keyword) {
		return errors.New("攻略不存在，无法删除")
	}
	affected, err := DB.Where("keyword = ?", keyword).Delete(&database.Strategy{})
	if err != nil {
		logger.Errorf("删除攻略 %s 失败 %v", keyword, err)
		return errors.New(e.SystemErrorNote)
	}
	if affected == 0 {
		return errors.New("攻略不存在")
	}
	Cache.Delete(CacheKeyStrategyKeywords)
	Cache.Delete(fmt.Sprintf(CacheKeyStrategyData, keyword))
	return nil
}
