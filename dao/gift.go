package dao

import "bcjh-bot/model/database"

// FindGuestGiftsByRecipeName 根据菜谱名字查询对应的贵客礼物数据
func FindGuestGiftsByRecipeName(recipe string) ([]database.GuestGift, error) {
	var gifts []database.GuestGift
	err := DB.Where("recipe = ?", recipe).Find(&gifts)
	return gifts, err
}
