package service

import (
	"bcjh-bot/config"
	"bcjh-bot/logger"
	"bcjh-bot/model"
	"bcjh-bot/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

//更新数据
//导出数据库数据->删库->重新同步表结构->插入数据
func UpdateData(msg *model.OneBotMsg, args []string) {
	logger.Info("更新数据, 参数:", args)
	dumpTime := time.Now().Format("2006010021504")
	DumpFilePath := config.AppConfig.DBConfig.ExportDir + "/DBDataDump" + dumpTime + ".sql"
	err := model.DB.DumpAllToFile(DumpFilePath)
	if err != nil {
		logger.Error("导出数据库数据失败!", err)
		return
	}
	logger.Info("导出旧数据完毕")

	gameData, err := RequestData()
	if err != nil {
		logger.Error("获取图鉴网数据失败!", err)
		return
	}
	logger.Infof("获取到图鉴网数据%+v", gameData)

	//开启事务，删除原有数据
	session := model.DB.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		logger.Error("更新数据失败!", err)
		return
	}

	//删除原数据
	tables := model.TablesName

	for _, table := range tables {
		sql := "DELETE FROM `" + table + "`"
		_, err = session.Exec(sql)
		if err != nil {
			logger.Error("删除旧数据出错!", err)
			session.Rollback()
			return
		}
	}

	//插入新数据
	//插入厨师数据
	chefsData := gameData.Chefs
	chefs := make([]model.Chef, 0)
	for _, chefData := range chefsData {
		chef := model.Chef{
			ChefId:        chefData.ChefId,
			Name:          chefData.Name,
			Rarity:        chefData.Rarity,
			Origin:        chefData.Origin,
			GalleryId:     chefData.GalleryId,
			Stirfry:       chefData.Stirfry,
			Bake:          chefData.Bake,
			Boil:          chefData.Boil,
			Steam:         chefData.Steam,
			Fry:           chefData.Fry,
			Cut:           chefData.Cut,
			Meat:          chefData.Meat,
			Flour:         chefData.Flour,
			Fish:          chefData.Fish,
			Vegetable:     chefData.Vegetable,
			SkillId:       chefData.SkillId,
			UltimateGoal:  chefData.UltimateGoal,
			UltimateSkill: chefData.UltimateSkill,
		}
		if len(chefData.Tags) > 0 {
			chef.Gender = chefData.Tags[0]
		}
		chefs = append(chefs, chef)
	}
	_, err = session.Insert(&chefs)
	if err != nil {
		logger.Error("插入厨师数据出错!", err)
		session.Rollback()
		return
	}
	logger.Info("更新厨师数据完毕!")

	//更新厨具数据

	//更新菜谱数据

	err = session.Commit()
	if err != nil {
		logger.Error("更新数据失败!", err)
		return
	}
	//关闭事务，发送成功消息
	logger.Info("更新数据完毕")
}

//从图鉴网爬取数据
func RequestData() (model.GameData, error) {
	var gameData model.GameData
	r, err := http.Get(util.FoodGameDataUrl)
	if err != nil {
		return gameData, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return gameData, err
	}
	err = json.Unmarshal(body, &gameData)
	return gameData, err
}
