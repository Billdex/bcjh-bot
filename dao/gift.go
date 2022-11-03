package dao

import (
	"bcjh-bot/model/database"
)

const CacheKeyGuestGiftList = "guest_gift_list"

// ClearGuestGiftsCache 清除贵客礼物缓存
func ClearGuestGiftsCache() {
	Cache.Delete(CacheKeyGuestGiftList)
}

// FindAllGuestGifts 查询全部贵客礼物数据
func FindAllGuestGifts() ([]database.GuestGift, error) {
	var gifts []database.GuestGift
	err := SimpleFindDataWithCache(CacheKeyGuestGiftList, &gifts, func(dest interface{}) error {
		return DB.OrderBy("guest_id").Find(dest)
	})
	return gifts, err
}

// GetRecipeGuestGiftsMap 获取 map 格式的贵客礼物关联数据，key 为菜谱名
func GetRecipeGuestGiftsMap() (map[string][]database.GuestGift, error) {
	gifts, err := FindAllGuestGifts()
	if err != nil {
		return nil, err
	}
	mResult := make(map[string][]database.GuestGift)
	for _, gift := range gifts {
		mResult[gift.Recipe] = append(mResult[gift.Recipe], gift)
	}
	return mResult, nil
}
