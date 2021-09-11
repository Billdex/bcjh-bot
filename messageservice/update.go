package messageservice

import (
	"bcjh-bot/config"
	"bcjh-bot/model/database"
	"bcjh-bot/model/gamedata"
	"bcjh-bot/scheduler"
	"bcjh-bot/util/logger"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"xorm.io/xorm"
)

const (
	foodGameGithubURLBase = "https://foodgame.github.io"
	foodGameGiteeURLBase  = "https://foodgame.gitee.io"
	bcjhURLBase           = "https://bcjh.gitee.io"

	dataURI        = "/data/data.min.json"
	imageCSSURI    = "/css/image.css"
	chefImageURI   = "/images/chef_retina.png"
	recipeImageURI = "/images/recipe_retina.png"
	equipImageURI  = "/images/equip_retina.png"
)

var updateState = false
var updateMux sync.Mutex

func setUpdateState(state bool) {
	updateMux.Lock()
	defer updateMux.Unlock()
	updateState = state
}

func getUpdateState() bool {
	updateMux.Lock()
	defer updateMux.Unlock()
	return updateState
}

func UpdateData(c *scheduler.Context) {
	// 防止在未更新完毕的情况下调用更新
	if getUpdateState() {
		_, _ = c.Reply("数据正在更新中")
		return
	}
	setUpdateState(true)
	defer setUpdateState(false)
	var baseURL string
	switch strings.TrimSpace(c.PretreatedMessage) {
	case "github":
		baseURL = foodGameGithubURLBase
	case "gitee":
		baseURL = foodGameGiteeURLBase
	case "白菜菊花":
		baseURL = bcjhURLBase
	default:
		baseURL = foodGameGiteeURLBase
	}
	_, _ = c.Reply(fmt.Sprintf("开始导入数据, 数据源:\n%s", baseURL))
	updateStart := time.Now().UnixNano()

	start := time.Now().UnixNano()
	gameData, err := requestData(baseURL + dataURI)
	if err != nil {
		logger.Error("获取图鉴网数据失败!", err)
		_, _ = c.Reply("获取图鉴网数据失败!")
		return
	}
	requestConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("获取图鉴网数据完毕, 耗时%s", requestConsume)

	// 导入sql数据
	start = time.Now().UnixNano()
	err = importDirAllSqlFile(database.DB, config.AppConfig.Resource.Sql)
	if err != nil {
		logger.Error("导入预配置sql数据出错!", err)
		_, _ = c.Reply("导入预配置sql数据出错!")
		return
	}
	importSqlConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("导入预配置sql数据完毕, 耗时%s", importSqlConsume)

	// 更新数据
	// 更新厨师数据
	start = time.Now().UnixNano()
	err = updateChefs(gameData.Chefs)
	if err != nil {
		logger.Error("更新厨师数据出错!", err)
		_, _ = c.Reply("更新厨师数据出错!")
		return
	}
	chefConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新厨师数据完毕, 耗时%s", chefConsume)

	// 更新厨具数据
	start = time.Now().UnixNano()
	err = updateEquips(gameData.Equips)
	if err != nil {
		logger.Error("更新厨具数据出错!", err)
		_, _ = c.Reply("更新厨具数据出错!")
		return
	}
	equipConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新厨具数据完毕, 耗时%s", equipConsume)

	// 更新菜谱数据
	start = time.Now().UnixNano()
	err = updateRecipes(gameData.Recipes)
	if err != nil {
		logger.Error("更新菜谱数据出错!", err)
		_, _ = c.Reply("更新菜谱数据出错!")
		return
	}
	recipeConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新菜谱数据完毕, 耗时%s", recipeConsume)

	// 更新合成菜谱数据
	start = time.Now().UnixNano()
	err = updateCombos(gameData.Combos)
	if err != nil {
		logger.Error("更新后厨合成菜谱数据出错!", err)
		_, _ = c.Reply("更新后厨合成菜谱数据出错!")
		return
	}
	comboConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新后厨合成菜谱数据完毕, 耗时%s", comboConsume)

	// 更新贵客数据
	start = time.Now().UnixNano()
	err = updateGuests(gameData.Guests)
	if err != nil {
		logger.Error("更新贵客数据出错!", err)
		_, _ = c.Reply("更新贵客数据出错!")
		return
	}
	guestConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新贵客数据完毕, 耗时%s", guestConsume)

	// 更新食材数据
	start = time.Now().UnixNano()
	err = updateMaterials(gameData.Materials)
	if err != nil {
		logger.Error("更新食材数据出错!", err)
		_, _ = c.Reply("更新食材数据出错!")
		return
	}
	materialConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新食材数据完毕, 耗时%s", materialConsume)

	// 更新技能数据
	start = time.Now().UnixNano()
	err = updateSkills(gameData.Skills)
	if err != nil {
		logger.Error("更新技能数据出错!", err)
		_, _ = c.Reply("更新技能数据出错!")
		return
	}
	skillConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新技能数据完毕, 耗时%s", skillConsume)

	// 更新装修家具数据
	start = time.Now().UnixNano()
	err = updateDecorations(gameData.Decorations)
	if err != nil {
		logger.Error("更新装修家具数据出错!", err)
		_, _ = c.Reply("更新装修家具数据出错!")
		return
	}
	DecorationConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新装修家具数据完毕, 耗时%s", DecorationConsume)

	// 更新调料数据
	start = time.Now().UnixNano()
	err = updateCondiments(gameData.Condiments)
	if err != nil {
		logger.Error("更新调料数据出错!", err)
		_, _ = c.Reply("更新调料数据出错!")
		return
	}
	CondimentConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新调料数据完毕, 耗时%s", CondimentConsume)

	// 更新任务数据
	start = time.Now().UnixNano()
	err = updateQuests(gameData.Quests)
	if err != nil {
		logger.Error("更新任务数据出错!", err)
		_, _ = c.Reply("更新任务数据出错!")
		return
	}
	QuestConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新任务数据完毕, 耗时%s", QuestConsume)

	// 解析ImgCSS数据
	start = time.Now().UnixNano()
	imgCSS, err := ResolvingImgCSS(baseURL + imageCSSURI)
	if err != nil {
		logger.Error("解析ImgCSS出错!", err)
		_, _ = c.Reply("解析ImgCSS出错!")
		return
	}
	imgCSSConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("解析ImgCSS完毕, 耗时%s", imgCSSConsume)

	// 更新厨师图鉴图片数据
	start = time.Now().UnixNano()
	chefs := make([]database.Chef, 0)
	err = database.DB.Asc("gallery_id").Find(&chefs)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_, _ = c.Reply("更新厨师图鉴图片数据出错!")
		return
	}
	err = ChefInfoToImage(chefs, baseURL+chefImageURI, imgCSS)
	if err != nil {
		logger.Error("更新厨师图鉴图片数据出错!", err)
		_, _ = c.Reply("更新厨师图鉴图片数据出错!")
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
		_, _ = c.Reply("更新菜谱图鉴图片数据出错!")
		return
	}
	err = RecipeInfoToImage(recipes, baseURL+recipeImageURI, imgCSS)
	if err != nil {
		logger.Error("更新菜谱图鉴图片数据出错!", err)
		_, _ = c.Reply("更新菜谱图鉴图片数据出错!")
		return
	}
	recipeImgConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新菜谱图鉴图片数据完毕, 耗时%s", recipeImgConsume)

	// 更新厨具图鉴图片数据
	start = time.Now().UnixNano()
	equips := make([]database.Equip, 0)
	err = database.DB.Asc("gallery_id").Find(&equips)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_, _ = c.Reply("更新厨具图鉴图片数据出错!")
		return
	}
	err = EquipmentInfoToImage(equips, baseURL+equipImageURI, imgCSS)
	if err != nil {
		logger.Error("更新厨具图鉴图片数据出错!", err)
		_, _ = c.Reply("更新厨具图鉴图片数据出错!")
		return
	}
	equipImgConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-start)/1e9)
	logger.Infof("更新厨具图鉴图片数据完毕, 耗时%s", equipImgConsume)

	// 发送成功消息
	logger.Info("更新数据完毕")
	var strBdr = strings.Builder{}
	updateConsume := fmt.Sprintf("%.2fs", (float64)(time.Now().UnixNano()-updateStart)/1e9)
	strBdr.WriteString(fmt.Sprintf("更新数据完毕, 累计耗时%s\n", updateConsume))
	strBdr.WriteString(fmt.Sprintf("抓取图鉴网数据耗时%s\n", requestConsume))
	strBdr.WriteString(fmt.Sprintf("导入预配置sql数据耗时%s\n", importSqlConsume))
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
	strBdr.WriteString(fmt.Sprintf("解析ImgCSS数据耗时%s\n", imgCSSConsume))
	strBdr.WriteString(fmt.Sprintf("更新厨师图鉴图片数据耗时%s\n", chefImgConsume))
	strBdr.WriteString(fmt.Sprintf("更新菜谱图鉴图片数据耗时%s\n", recipeImgConsume))
	strBdr.WriteString(fmt.Sprintf("更新厨具图鉴图片数据耗时%s", equipImgConsume))
	_, _ = c.Reply(strBdr.String())
}

// 从图鉴网爬取数据
func requestData(url string) (gamedata.GameData, error) {
	var gameData gamedata.GameData
	r, err := http.Get(url)
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

// 导入预配置sql
func importDirAllSqlFile(engine *xorm.Engine, dir string) error {
	tableMap := map[string]interface{}{
		"guest.sql":      database.Guest{},
		"laboratory.sql": database.Laboratory{},
	}
	for file, table := range tableMap {
		if exist, err := engine.IsTableExist(table); err != nil {
			return err
		} else if exist {
			continue
		} else {
			_, err = engine.ImportFile(fmt.Sprintf("%s/%s", dir, file))
			if err != nil {
				return err
			}
		}
	}
	return nil
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
		_ = session.Rollback()
		return err
	}
	chefs := make([]database.Chef, 0)
	for _, chefData := range chefsData {
		chef := database.Chef{
			ChefId:        chefData.ChefId,
			Name:          chefData.Name,
			Rarity:        chefData.Rarity,
			Origin:        strings.ReplaceAll(chefData.Origin, "<br>", ", "),
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
		_ = session.Rollback()
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
		_ = session.Rollback()
		return err
	}
	equips := make([]database.Equip, 0)
	for _, equipData := range equipsData {
		equips = append(equips, database.Equip{
			EquipId:   equipData.EquipId,
			Name:      equipData.Name,
			GalleryId: equipData.GalleryId,
			Origin:    strings.ReplaceAll(equipData.Origin, "<br>", ", "),
			Rarity:    equipData.Rarity,
			Skills:    equipData.Skills,
		})
	}
	_, err = session.Insert(&equips)
	if err != nil {
		_ = session.Rollback()
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
		_ = session.Rollback()
		return err
	}
	// 删除菜谱-食材关系
	sql = fmt.Sprintf("DELETE FROM `%s`", new(database.RecipeMaterial).TableName())
	_, err = session.Exec(sql)
	if err != nil {
		_ = session.Rollback()
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
			Origin:         strings.ReplaceAll(recipeData.Origin, "<br>", ", "),
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
		_ = session.Rollback()
		return err
	}
	_, err = session.Insert(&recipes)
	if err != nil {
		_ = session.Rollback()
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
		_ = session.Rollback()
		return err
	}
	for _, combo := range combosData {
		comboRecipe := new(database.Recipe)
		has, err := session.Where("recipe_id = ?", combo.RecipeId).Get(comboRecipe)
		if err != nil {
			_ = session.Rollback()
			return err
		}
		if !has {
			_ = session.Rollback()
			return errors.New(fmt.Sprintf("未查询到后厨合成菜谱%d信息", combo.RecipeId))
		}
		for _, recipeId := range combo.Recipes {
			recipe := new(database.Recipe)
			recipe.Combo = comboRecipe.Name
			_, err = session.Where("recipe_id = ?", recipeId).Update(recipe)
			if err != nil {
				_ = session.Rollback()
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
		_ = session.Rollback()
		return err
	}
	// 从guest表中获取预设的贵客编号
	guestInfo := make([]database.Guest, 0)
	err = database.DB.Find(&guestInfo)
	if err != nil {
		logger.Error("数据库查询出错", err)
		_ = session.Rollback()
		return err
	}
	guestMap := make(map[string]string)
	for _, guest := range guestInfo {
		guestMap[guest.GuestName] = guest.GuestId
	}
	// 从菜谱表中获取菜谱相关信息
	recipeInfo := make([]database.Recipe, 0)
	err = database.DB.Find(&recipeInfo)
	if err != nil {
		logger.Error("数据库查询出错", err)
		_ = session.Rollback()
		return err
	}
	recipeMap := make(map[string]database.Recipe)
	for _, recipe := range recipeInfo {
		recipeMap[recipe.Name] = recipe
	}

	guests := make([]database.GuestGift, 0)
	for _, guestData := range guestsData {
		for _, gift := range guestData.Gifts {
			guest := database.GuestGift{
				GuestId:   guestMap[guestData.Name],
				GuestName: guestData.Name,
				Antique:   gift.Antique,
				Recipe:    gift.Recipe,
				TotalTime: recipeMap[gift.Recipe].TotalTime,
			}
			guests = append(guests, guest)
		}
	}
	_, err = session.Insert(&guests)
	if err != nil {
		_ = session.Rollback()
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
		_ = session.Rollback()
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
		_ = session.Rollback()
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
		_ = session.Rollback()
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
		_ = session.Rollback()
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
		_ = session.Rollback()
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
		_ = session.Rollback()
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
		_ = session.Rollback()
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
		_ = session.Rollback()
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
		_ = session.Rollback()
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
		_ = session.Rollback()
		return err
	}
	err = session.Commit()
	return err
}

// 解析ImgCSS数据
func ResolvingImgCSS(cssURL string) (*gamedata.ImgCSS, error) {
	imgCSS := new(gamedata.ImgCSS)
	imgCSS.ChefImg = make(map[int]gamedata.ObjImgInfo)
	imgCSS.RecipeImg = make(map[int]gamedata.ObjImgInfo)
	imgCSS.EquipImg = make(map[int]gamedata.ObjImgInfo)

	r, err := http.Get(cssURL)
	if err != nil {
		return imgCSS, err
	}
	defer r.Body.Close()

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return imgCSS, err
	}
	chefRegStr := ".icon-chef.chef_([0-9]+?){background-position:(.?[0-9]+)px (.?[0-9]+)px;width:(.?[0-9]+)px;height:(.?[0-9]+)px;}"
	pattern := regexp.MustCompile(chefRegStr)
	allIndexes := pattern.FindAllSubmatchIndex(buf, -1)
	for _, loc := range allIndexes {
		id, _ := strconv.Atoi(string(buf[loc[2]:loc[3]]))
		x, _ := strconv.Atoi(string(buf[loc[4]:loc[5]]))
		y, _ := strconv.Atoi(string(buf[loc[6]:loc[7]]))
		w, _ := strconv.Atoi(string(buf[loc[8]:loc[9]]))
		h, _ := strconv.Atoi(string(buf[loc[10]:loc[11]]))
		imgCSS.ChefImg[id] = gamedata.ObjImgInfo{
			Id:     id,
			X:      -x,
			Y:      -y,
			Width:  w,
			Height: h,
		}
	}
	recipeRegStr := ".icon-recipe.recipe_([0-9]+?){background-position:(-?[0-9]+)px (-?[0-9]+)px;width:([0-9]+)px;height:([0-9]+)px;"
	pattern = regexp.MustCompile(recipeRegStr)
	allIndexes = pattern.FindAllSubmatchIndex(buf, -1)
	for _, loc := range allIndexes {
		id, _ := strconv.Atoi(string(buf[loc[2]:loc[3]]))
		x, _ := strconv.Atoi(string(buf[loc[4]:loc[5]]))
		y, _ := strconv.Atoi(string(buf[loc[6]:loc[7]]))
		w, _ := strconv.Atoi(string(buf[loc[8]:loc[9]]))
		h, _ := strconv.Atoi(string(buf[loc[10]:loc[11]]))
		imgCSS.RecipeImg[id] = gamedata.ObjImgInfo{
			Id:     id,
			X:      -x,
			Y:      -y,
			Width:  w,
			Height: h,
		}
	}
	equipRegStr := ".icon-equip.equip_([0-9]+?){background-position:(-?[0-9]+)px (-?[0-9]+)px;width:([0-9]+)px;height:([0-9]+)px;"
	pattern = regexp.MustCompile(equipRegStr)
	allIndexes = pattern.FindAllSubmatchIndex(buf, -1)
	for _, loc := range allIndexes {
		id, _ := strconv.Atoi(string(buf[loc[2]:loc[3]]))
		x, _ := strconv.Atoi(string(buf[loc[4]:loc[5]]))
		y, _ := strconv.Atoi(string(buf[loc[6]:loc[7]]))
		w, _ := strconv.Atoi(string(buf[loc[8]:loc[9]]))
		h, _ := strconv.Atoi(string(buf[loc[10]:loc[11]]))
		imgCSS.EquipImg[id] = gamedata.ObjImgInfo{
			Id:     id,
			X:      -x,
			Y:      -y,
			Width:  w,
			Height: h,
		}
	}

	return imgCSS, nil
}
