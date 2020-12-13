package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/config"
	"bcjh-bot/model/database"
	"bcjh-bot/model/gamedata"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// 更新数据
// 导出数据库数据->删库->插入新数据
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
	updateStart := time.Now().UnixNano()
	start := time.Now().UnixNano()
	dumpTime := time.Now().Format("200601021504")
	DumpFilePath := config.AppConfig.DB.ExportDir + "/DBDataDump" + dumpTime + ".sql"
	err = database.DB.DumpAllToFile(DumpFilePath)
	if err != nil {
		logger.Error("导出旧数据失败!", err)
		_ = bot.SendMessage(c, "导出旧数据失败!")
		return
	}
	dumpConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("导出旧数据完毕, 耗时%s", dumpConsume)

	start = time.Now().UnixNano()
	gameData, err := requestData()
	if err != nil {
		logger.Error("获取图鉴网数据失败!", err)
		_ = bot.SendMessage(c, "获取图鉴网数据失败!")
		return
	}
	requestConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("获取图鉴网数据完毕, 耗时%s", requestConsume)
	logger.Debug("数据内容为:%+v", gameData)

	// 更新数据
	// 更新厨师数据
	start = time.Now().UnixNano()
	err = updateChefs(gameData.Chefs)
	if err != nil {
		logger.Error("更新厨师数据出错!", err)
		_ = bot.SendMessage(c, "更新厨师数据出错!")
		return
	}
	chefConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新厨师数据完毕, 耗时%s", chefConsume)

	// 更新厨具数据
	start = time.Now().UnixNano()
	err = updateEquips(gameData.Equips)
	if err != nil {
		logger.Error("更新厨具数据出错!", err)
		_ = bot.SendMessage(c, "更新厨具数据出错!")
		return
	}
	equipConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新厨具数据完毕, 耗时%s", equipConsume)

	// 更新菜谱数据
	start = time.Now().UnixNano()
	err = updateRecipes(gameData.Recipes)
	if err != nil {
		logger.Error("更新菜谱数据出错!", err)
		_ = bot.SendMessage(c, "更新菜谱数据出错!")
		return
	}
	recipeConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新菜谱数据完毕, 耗时%s", recipeConsume)

	// 更新合成菜谱数据
	start = time.Now().UnixNano()
	err = updateCombos(gameData.Combos)
	if err != nil {
		logger.Error("更新后厨合成菜谱数据出错!", err)
		_ = bot.SendMessage(c, "更新后厨合成菜谱数据出错!")
		return
	}
	comboConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新后厨合成菜谱数据完毕, 耗时%s", comboConsume)

	// 更新贵客数据
	start = time.Now().UnixNano()
	err = updateGuests(gameData.Guests)
	if err != nil {
		logger.Error("更新贵客数据出错!", err)
		_ = bot.SendMessage(c, "更新贵客数据出错!")
		return
	}
	guestConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新贵客数据完毕, 耗时%s", guestConsume)

	// 更新食材数据
	start = time.Now().UnixNano()
	err = updateMaterials(gameData.Materials)
	if err != nil {
		logger.Error("更新食材数据出错!", err)
		_ = bot.SendMessage(c, "更新食材数据出错!")
		return
	}
	materialConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新食材数据完毕, 耗时%s", materialConsume)

	// 更新技能数据
	start = time.Now().UnixNano()
	err = updateSkills(gameData.Skills)
	if err != nil {
		logger.Error("更新技能数据出错!", err)
		_ = bot.SendMessage(c, "更新技能数据出错!")
		return
	}
	skillConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新技能数据完毕, 耗时%s", skillConsume)

	// 更新装修家具数据
	start = time.Now().UnixNano()
	err = updateDecorations(gameData.Decorations)
	if err != nil {
		logger.Error("更新装修家具数据出错!", err)
		_ = bot.SendMessage(c, "更新装修家具数据出错!")
		return
	}
	DecorationConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新装修家具数据完毕, 耗时%s", DecorationConsume)

	// 更新调料数据
	start = time.Now().UnixNano()
	err = updateCondiments(gameData.Condiments)
	if err != nil {
		logger.Error("更新调料数据出错!", err)
		_ = bot.SendMessage(c, "更新调料数据出错!")
		return
	}
	CondimentConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新调料数据完毕, 耗时%s", CondimentConsume)

	// 更新任务数据
	start = time.Now().UnixNano()
	err = updateQuests(gameData.Quests)
	if err != nil {
		logger.Error("更新任务数据出错!", err)
		_ = bot.SendMessage(c, "更新任务数据出错!")
		return
	}
	QuestConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新任务数据完毕, 耗时%s", QuestConsume)

	// 更新厨师图鉴图片数据
	start = time.Now().UnixNano()
	chefs := make([]database.Chef, 0)
	err = database.DB.Asc("gallery_id").Find(&chefs)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_ = bot.SendMessage(c, "更新厨师图鉴图片数据出错!")
		return
	}
	err = ChefInfoToImage(chefs)
	if err != nil {
		logger.Error("更新厨师图鉴图片数据出错!", err)
		_ = bot.SendMessage(c, "更新厨师图鉴图片数据出错!")
		return
	}
	chefImgConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新厨师图鉴图片数据完毕, 耗时%s", chefImgConsume)

	// 更新菜谱图鉴图片数据
	start = time.Now().UnixNano()
	recipes := make([]database.Recipe, 0)
	err = database.DB.Asc("gallery_id").Find(&recipes)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_ = bot.SendMessage(c, "更新菜谱图鉴图片数据出错!")
		return
	}
	err = RecipeInfoToImage(recipes)
	if err != nil {
		logger.Error("更新菜谱图鉴图片数据出错!", err)
		_ = bot.SendMessage(c, "更新菜谱图鉴图片数据出错!")
		return
	}
	recipeImgConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新菜谱图鉴图片数据完毕, 耗时%s", chefImgConsume)

	// 发送成功消息
	logger.Info("更新数据完毕")
	var strBdr = strings.Builder{}
	updateConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-updateStart)/1e9)
	strBdr.WriteString(fmt.Sprintf("更新数据完毕, 累计耗时%s\n", updateConsume))
	strBdr.WriteString(fmt.Sprintf("导出旧数据耗时%s\n", dumpConsume))
	strBdr.WriteString(fmt.Sprintf("抓取图鉴网数据耗时%s\n", requestConsume))
	strBdr.WriteString(fmt.Sprintf("更新厨师数据耗时%s\n", chefConsume))
	strBdr.WriteString(fmt.Sprintf("更新厨具数据耗时%s\n", equipConsume))
	strBdr.WriteString(fmt.Sprintf("更新菜谱数据耗时%s\n", recipeConsume))
	strBdr.WriteString(fmt.Sprintf("更新后厨合成菜谱数据耗时%s\n", comboConsume))
	strBdr.WriteString(fmt.Sprintf("更新贵客数据耗时%s\n", guestConsume))
	strBdr.WriteString(fmt.Sprintf("更新食材数据耗时%s\n", materialConsume))
	strBdr.WriteString(fmt.Sprintf("更新技能数据耗时%s\n", skillConsume))
	strBdr.WriteString(fmt.Sprintf("更新装修家具数据耗时%s\n", DecorationConsume))
	strBdr.WriteString(fmt.Sprintf("更新调料数据耗时%s\n", CondimentConsume))
	strBdr.WriteString(fmt.Sprintf("更新任务数据耗时%s\n", QuestConsume))
	strBdr.WriteString(fmt.Sprintf("更新厨师图鉴图片数据耗时%s", chefImgConsume))
	strBdr.WriteString(fmt.Sprintf("更新菜谱图鉴图片数据耗时%s", recipeImgConsume))
	err = bot.SendMessage(c, strBdr.String())
	if err != nil {
		logger.Error("发送消息失败!", err)
	}
}

// 从图鉴网爬取数据
func requestData() (gamedata.GameData, error) {
	var gameData gamedata.GameData
	r, err := http.Get(util.FoodGameDataURL)
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

// 更新厨师信息
func updateChefs(chefsData []gamedata.ChefData) error {
	session := database.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM `%s`", new(database.Chef).TableName())
	_, err = session.Exec(sql)
	if err != nil {
		session.Rollback()
		return err
	}
	chefs := make([]database.Chef, 0)
	for _, chefData := range chefsData {
		chef := database.Chef{
			ChefId:        chefData.ChefId,
			Name:          chefData.Name,
			Rarity:        chefData.Rarity,
			Origin:        strings.ReplaceAll(chefData.Origin, "<br>", ","),
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
			Sweet:         chefData.Sweet,
			Sour:          chefData.Sour,
			Spicy:         chefData.Spicy,
			Salty:         chefData.Salty,
			Bitter:        chefData.Bitter,
			Tasty:         chefData.Tasty,
			SkillId:       chefData.SkillId,
			UltimateGoals: chefData.UltimateGoals,
			UltimateSkill: chefData.UltimateSkill,
		}
		if len(chefData.Tags) > 0 {
			chef.Gender = chefData.Tags[0]
		}
		chefs = append(chefs, chef)
	}
	_, err = session.Insert(&chefs)
	if err != nil {
		session.Rollback()
		return err
	}
	err = session.Commit()
	return err
}

// 更新厨具信息
func updateEquips(equipsData []gamedata.EquipData) error {
	session := database.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM `%s`", new(database.Equip).TableName())
	_, err = session.Exec(sql)
	if err != nil {
		session.Rollback()
		return err
	}
	equips := make([]database.Equip, 0)
	for _, equipData := range equipsData {
		equips = append(equips, database.Equip{
			EquipId:   equipData.EquipId,
			Name:      equipData.Name,
			GalleryId: equipData.GalleryId,
			Origin:    strings.ReplaceAll(equipData.Origin, "<br>", ","),
			Rarity:    equipData.Rarity,
			Skills:    equipData.Skills,
		})
	}
	_, err = session.Insert(&equips)
	if err != nil {
		session.Rollback()
		return err
	}
	err = session.Commit()
	return err
}

// 更新菜谱信息
func updateRecipes(recipesData []gamedata.RecipeData) error {
	session := database.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	// 删除菜谱数据
	sql := fmt.Sprintf("DELETE FROM `%s`", new(database.Recipe).TableName())
	_, err = session.Exec(sql)
	if err != nil {
		session.Rollback()
		return err
	}
	// 删除菜谱-食材关系
	sql = fmt.Sprintf("DELETE FROM `%s`", new(database.RecipeMaterial).TableName())
	_, err = session.Exec(sql)
	if err != nil {
		session.Rollback()
		return err
	}
	recipes := make([]database.Recipe, 0)
	materials := make([]database.RecipeMaterial, 0)
	for _, recipeData := range recipesData {
		recipe := database.Recipe{
			RecipeId:       recipeData.RecipeId,
			Name:           recipeData.Name,
			GalleryId:      recipeData.GalleryId,
			Rarity:         recipeData.Rarity,
			Origin:         strings.ReplaceAll(recipeData.Origin, "<br>", ","),
			Stirfry:        recipeData.Stirfry,
			Bake:           recipeData.Bake,
			Boil:           recipeData.Boil,
			Steam:          recipeData.Steam,
			Fry:            recipeData.Fry,
			Cut:            recipeData.Cut,
			Condiment:      recipeData.Condiment,
			Price:          recipeData.Price,
			ExPrice:        recipeData.ExPrice,
			GoldEfficiency: recipeData.Price * 3600 / recipeData.Time,
			Gift:           recipeData.Gift,
			Time:           recipeData.Time,
			Limit:          recipeData.Limit,
			TotalTime:      recipeData.Time * recipeData.Limit,
			Unlock:         recipeData.Unlock,
			Combo:          "-",
		}
		// 插入升阶贵客信息
		guests := make([]string, 0)
		for _, guest := range recipeData.Guests {
			guests = append(guests, guest.Guest)
		}
		recipe.Guests = guests
		// 插入耗材信息
		materialSum := 0
		for _, materialData := range recipeData.Materials {
			materials = append(materials, database.RecipeMaterial{
				RecipeGalleryId: recipeData.GalleryId,
				MaterialId:      materialData.MaterialId,
				Quantity:        materialData.Quantity,
				Efficiency:      materialData.Quantity * 3600 / recipe.Time,
			})
			materialSum += materialData.Quantity
		}
		recipe.MaterialEfficiency = materialSum * 3600 / recipeData.Time
		recipes = append(recipes, recipe)
	}
	_, err = session.Insert(&materials)
	if err != nil {
		session.Rollback()
		return err
	}
	_, err = session.Insert(&recipes)
	if err != nil {
		session.Rollback()
		return err
	}
	err = session.Commit()
	return err
}

// 更新后厨合成菜信息
func updateCombos(combosData []gamedata.ComboData) error {
	session := database.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	recipes := new(database.Recipe)
	recipes.Combo = "-"
	_, err = session.Where("combo <> ?", "-").Update(recipes)
	if err != nil {
		session.Rollback()
		return err
	}
	for _, combo := range combosData {
		comboRecipe := new(database.Recipe)
		has, err := session.Where("recipe_id = ?", combo.RecipeId).Get(comboRecipe)
		if err != nil {
			session.Rollback()
			return err
		}
		if !has {
			session.Rollback()
			return errors.New(fmt.Sprintf("未查询到后厨合成菜谱%d信息", combo.RecipeId))
		}
		for _, recipeId := range combo.Recipes {
			recipe := new(database.Recipe)
			recipe.Combo = comboRecipe.Name
			_, err = session.Where("recipe_id = ?", recipeId).Update(recipe)
			if err != nil {
				session.Rollback()
				return err
			}
		}
	}
	err = session.Commit()
	return err
}

// 更新贵客信息
func updateGuests(guestsData []gamedata.GuestData) error {
	session := database.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM `%s`", new(database.GuestGift).TableName())
	_, err = session.Exec(sql)
	if err != nil {
		session.Rollback()
		return err
	}
	guests := make([]database.GuestGift, 0)
	for p, guestData := range guestsData {
		// 图鉴网未指明贵客id，只能按数据顺序排序，因此30后的图鉴编号有误!!!
		for _, gift := range guestData.Gifts {
			guest := database.GuestGift{
				GuestId:   fmt.Sprintf("%03d", p+1),
				GuestName: guestData.Name,
				Antique:   gift.Antique,
				Recipe:    gift.Recipe,
			}
			guests = append(guests, guest)
		}
	}
	_, err = session.Insert(&guests)
	if err != nil {
		session.Rollback()
		return err
	}
	err = session.Commit()
	return err
}

func updateMaterials(materialsData []gamedata.MaterialData) error {
	session := database.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM `%s`", new(database.Material).TableName())
	_, err = session.Exec(sql)
	if err != nil {
		session.Rollback()
		return err
	}
	materials := make([]database.Material, 0)
	for _, materialData := range materialsData {
		materials = append(materials, database.Material{
			MaterialId: materialData.MaterialId,
			Name:       materialData.Name,
			Origin:     materialData.Origin,
		})
	}
	_, err = session.Insert(&materials)
	if err != nil {
		session.Rollback()
		return err
	}
	err = session.Commit()
	return err
}

func updateSkills(skillsData []gamedata.SkillData) error {
	session := database.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM `%s`", new(database.Skill).TableName())
	_, err = session.Exec(sql)
	if err != nil {
		session.Rollback()
		return err
	}
	skills := make([]database.Skill, 0)
	for _, skillData := range skillsData {
		skill := database.Skill{
			SkillId:     skillData.SkillId,
			Description: strings.ReplaceAll(skillData.Description, "<br>", ","),
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
	_, err = session.Insert(&skills)
	if err != nil {
		session.Rollback()
		return err
	}
	err = session.Commit()
	return err
}

// 更新家具信息
func updateDecorations(decorationsData []gamedata.Decoration) error {
	session := database.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM `%s`", new(database.Decoration).TableName())
	_, err = session.Exec(sql)
	if err != nil {
		session.Rollback()
		return err
	}
	decorations := make([]database.Decoration, 0)
	for _, decorationData := range decorationsData {
		skill := database.Decoration{
			Id: decorationData.Id,
			//Icon: decorationDat.//,
			Name:     decorationData.Name,
			TipMin:   decorationData.TipMin,
			TipMax:   decorationData.TipMax,
			TipTime:  decorationData.TipTime,
			Gold:     decorationData.Gold,
			Position: decorationData.Position,
			Suit:     decorationData.Suit,
			SuitGold: decorationData.SuitGold,
			Origin:   decorationData.Origin,
		}
		decorations = append(decorations, skill)
	}
	_, err = session.Insert(&decorations)
	if err != nil {
		session.Rollback()
		return err
	}
	err = session.Commit()
	return err
}

// 更新调料信息
func updateCondiments(condimentsData []gamedata.Condiment) error {
	session := database.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM `%s`", new(database.Condiment).TableName())
	_, err = session.Exec(sql)
	if err != nil {
		session.Rollback()
		return err
	}
	condiments := make([]database.Condiment, 0)
	for _, condimentData := range condimentsData {
		skill := database.Condiment{
			CondimentId: condimentData.CondimentId,
			Name:        condimentData.Name,
			Rarity:      condimentData.Rarity,
			Skill:       condimentData.Skill,
			Origin:      condimentData.Origin,
		}
		condiments = append(condiments, skill)
	}
	_, err = session.Insert(&condiments)
	if err != nil {
		session.Rollback()
		return err
	}
	err = session.Commit()
	return err
}

// 更新任务信息
func updateQuests(questsData []gamedata.QuestData) error {
	session := database.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM `%s`", new(database.Quest).TableName())
	_, err = session.Exec(sql)
	if err != nil {
		session.Rollback()
		return err
	}
	quests := make([]database.Quest, 0)
	for _, questData := range questsData {
		quest := database.Quest{
			QuestId:     questData.QuestId,
			QuestIdDisp: fmt.Sprintf("%v", questData.QuestIdDisp),
			Type:        questData.Type,
			Goal:        questData.Goal,
		}
		rewards := make([]database.QuestRewards, 0)
		for _, rewardData := range questData.Rewards {
			rewards = append(rewards, database.QuestRewards{
				Name:     rewardData.Name,
				Quantity: rewardData.Quantity,
			})
		}
		quest.Rewards = rewards
		quests = append(quests, quest)
	}
	_, err = session.Insert(&quests)
	if err != nil {
		session.Rollback()
		return err
	}
	err = session.Commit()
	return err
}
