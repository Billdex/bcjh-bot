package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/config"
	"bcjh-bot/logger"
	"bcjh-bot/model/database"
	"bcjh-bot/model/gamedata"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"xorm.io/xorm"
)

//更新数据
//导出数据库数据->删库->重新同步表结构->插入数据
func UpdateData(c *onebot.Context, args []string) {
	logger.Info("更新数据, 参数:", args)

	has, err := database.DB.Where("qq = ?", c.Sender.UserId).Exist(&database.Admin{})
	if err != nil {
		logger.Error("查询数据库出错", err)
		return
	}
	if !has {
		_ = bot.SendMessage(c, "你没有权限!")
		return
	}
	_ = bot.SendMessage(c, "开始更新数据")
	dumpTime := time.Now().Format("2006010021504")
	DumpFilePath := config.AppConfig.DBConfig.ExportDir + "/DBDataDump" + dumpTime + ".sql"
	err = database.DB.DumpAllToFile(DumpFilePath)
	if err != nil {
		logger.Error("导出旧数据失败!", err)
		_ = bot.SendMessage(c, "导出旧数据失败!")
		return
	}
	logger.Info("导出旧数据完毕")

	gameData, err := requestData()
	if err != nil {
		logger.Error("获取图鉴网数据失败!", err)
		_ = bot.SendMessage(c, "获取图鉴网数据失败!")
		return
	}
	logger.Infof("获取到图鉴网数据%+v", gameData)

	//开启事务，删除原有数据
	session := database.DB.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		logger.Error("开启事务失败!", err)
		_ = bot.SendMessage(c, "开启事务失败!")
		return
	}

	//删除原数据
	tables := database.TablesName
	for _, table := range tables {
		sql := "DELETE FROM `" + table + "`"
		_, err = session.Exec(sql)
		if err != nil {
			logger.Error("删除旧数据出错!", err)
			session.Rollback()
			_ = bot.SendMessage(c, "删除旧数据出错!")
			return
		}
	}
	logger.Info("删除旧数据完毕!")

	//插入新数据
	//插入厨师数据
	err = updateChefs(session, gameData.Chefs)
	if err != nil {
		logger.Error("插入厨师数据出错!", err)
		session.Rollback()
		_ = bot.SendMessage(c, "更新厨师数据出错!")
		return
	}
	logger.Info("更新厨师数据完毕!")

	//更新厨具数据
	err = updateEquips(session, gameData.Equips)
	if err != nil {
		logger.Error("插入厨具数据出错!", err)
		session.Rollback()
		_ = bot.SendMessage(c, "更新厨具数据出错!")
		return
	}
	logger.Info("更新厨具数据完毕!")

	//更新菜谱数据
	err = updateRecipes(session, gameData.Recipes)
	if err != nil {
		logger.Error("插入菜谱数据出错!", err)
		session.Rollback()
		_ = bot.SendMessage(c, "更新菜谱数据出错!")
		return
	}
	logger.Info("更新菜谱数据完毕!")

	//更新贵客数据
	err = updateGuests(session, gameData.Guests)
	if err != nil {
		logger.Error("插入贵客数据出错!", err)
		session.Rollback()
		_ = bot.SendMessage(c, "更新贵客数据出错!")
		return
	}
	logger.Info("更新贵客数据完毕!")

	//更新食材数据
	err = updateMaterials(session, gameData.Materials)
	if err != nil {
		logger.Error("插入食材数据出错!", err)
		session.Rollback()
		_ = bot.SendMessage(c, "更新食材数据出错!")
		return
	}
	logger.Info("更新食材数据完毕!")

	//更新技能数据
	err = updateSkills(session, gameData.Skills)
	if err != nil {
		logger.Error("插入技能数据出错!", err)
		session.Rollback()
		_ = bot.SendMessage(c, "更新技能数据出错!")
		return
	}
	logger.Info("更新技能数据完毕!")

	err = session.Commit()
	if err != nil {
		logger.Error("提交事务失败!", err)
		_ = bot.SendMessage(c, "提交事务失败!")
		return
	}
	//关闭事务，发送成功消息
	logger.Info("更新数据完毕")
	err = bot.SendMessage(c, "更新数据完毕")
	if err != nil {
		logger.Error("发送消息失败!", err)
	}
}

//从图鉴网爬取数据
func requestData() (gamedata.GameData, error) {
	var gameData gamedata.GameData
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

func updateChefs(session *xorm.Session, chefsData []gamedata.ChefData) error {
	chefs := make([]database.Chef, 0)
	for _, chefData := range chefsData {
		chef := database.Chef{
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
			UltimateGoals: chefData.UltimateGoals,
			UltimateSkill: chefData.UltimateSkill,
		}
		if len(chefData.Tags) > 0 {
			chef.Gender = chefData.Tags[0]
		}
		chefs = append(chefs, chef)
	}
	_, err := session.Insert(&chefs)
	return err
}

func updateEquips(session *xorm.Session, equipsData []gamedata.EquipData) error {
	equips := make([]database.Equip, 0)
	for _, equipData := range equipsData {
		equips = append(equips, database.Equip{
			EquipId:   equipData.EquipId,
			Name:      equipData.Name,
			GalleryId: equipData.GalleryId,
			Origin:    equipData.Origin,
			Rarity:    equipData.Rarity,
			Skills:    equipData.Skills,
		})
	}
	_, err := session.Insert(&equips)
	return err
}

func updateRecipes(session *xorm.Session, recipesData []gamedata.RecipeData) error {
	recipes := make([]database.Recipe, 0)
	for _, recipeData := range recipesData {
		recipe := database.Recipe{
			RecipeId:  recipeData.RecipeId,
			Name:      recipeData.Name,
			GalleryId: recipeData.GalleryId,
			Rarity:    recipeData.Rarity,
			Origin:    recipeData.Origin,
			Stirfry:   recipeData.Stirfry,
			Bake:      recipeData.Bake,
			Boil:      recipeData.Boil,
			Steam:     recipeData.Steam,
			Fry:       recipeData.Fry,
			Cut:       recipeData.Cut,
			Price:     recipeData.Price,
			ExPrice:   recipeData.ExPrice,
			Gift:      recipeData.Gift,
			Limit:     recipeData.Limit,
			Time:      recipeData.Time,
			Unlock:    recipeData.Unlock,
		}
		guests := make([]string, 0)
		for _, guest := range recipeData.Guests {
			guests = append(guests, guest.Guest)
		}
		recipe.Guests = guests
		recipes = append(recipes, recipe)
		materials := make([]database.RecipeMaterial, 0)
		for _, materialData := range recipeData.Materials {
			materials = append(materials, database.RecipeMaterial{
				MaterialId: materialData.MaterialId,
				Quantity:   materialData.Quantity,
			})
		}
		recipe.Materials = materials
	}
	_, err := session.Insert(&recipes)
	return err
}

func updateGuests(session *xorm.Session, guestsData []gamedata.GuestData) error {
	guests := make([]database.Guest, 0)
	for p, guestData := range guestsData {
		guest := database.Guest{
			GuestId:   p + 1,
			Name:      guestData.Name,
			GalleryId: fmt.Sprintf("%03d", p+1),
		}
		gifts := make([]database.GuestGift, 0)
		for _, gift := range guestData.Gifts {
			gifts = append(gifts, database.GuestGift{
				Antique: gift.Antique,
				Recipe:  gift.Recipe,
			})
		}
		guests = append(guests, guest)
	}
	_, err := session.Insert(&guests)
	return err
}

func updateMaterials(session *xorm.Session, materialsData []gamedata.MaterialData) error {
	materials := make([]database.Material, 0)
	for _, materialData := range materialsData {
		materials = append(materials, database.Material{
			MaterialId: materialData.MaterialId,
			Name:       materialData.Name,
			Origin:     materialData.Origin,
		})
	}
	_, err := session.Insert(&materials)
	return err
}

func updateSkills(session *xorm.Session, skillsData []gamedata.SkillData) error {
	skills := make([]database.Skill, 0)
	for _, skillData := range skillsData {
		skill := database.Skill{
			SkillId:     skillData.SkillId,
			Description: skillData.Description,
		}
		effects := make([]database.SkillEffect, 0)
		for _, effectData := range skillData.Effects {
			effects = append(effects, database.SkillEffect{
				Calculation: effectData.Calculation,
				Type:        effectData.Type,
				Condition:   effectData.Condition,
				Tag:         effectData.Tag,
				Value:       effectData.Value,
			})
		}
		skill.Effects = effects
		skills = append(skills, skill)
	}
	_, err := session.Insert(&skills)
	return err
}
