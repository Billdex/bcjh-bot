package messageservice

import (
	"bcjh-bot/config"
	"bcjh-bot/dao"
	"bcjh-bot/model/database"
	"bcjh-bot/model/gamedata"
	"bcjh-bot/scheduler"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	foodGameGithubURLBase = "https://foodgame.github.io"
	foodGameGiteeURLBase  = "https://foodgame.gitee.io"
	bcjhCfPageURLBase     = "https://bcjh.pages.dev"
	bcjhURLBase           = "https://bcjh.xyz"

	dataURI        = "/data/data.min.json"
	imageCSSURI    = "/css/image.css"
	chefImageURI   = "/images/chef_retina.png"
	recipeImageURI = "/images/recipe_retina.png"
	equipImageURI  = "/images/equip_retina.png"
)

var updateState = false
var updateMux sync.Mutex

func UpdateData(c *scheduler.Context) {
	// 防止在未更新完毕的情况下再次调用更新
	updateMux.Lock()
	if updateState == true {
		_, _ = c.Reply("数据正在更新中")
		updateMux.Unlock()
		return
	}
	updateState = true
	updateMux.Unlock()
	defer func() { updateState = false }()

	var baseURL string
	switch c.PretreatedMessage {
	case "lgithub":
		baseURL = foodGameGithubURLBase
	case "lgitee":
		baseURL = foodGameGiteeURLBase
	case "白菜菊花cf":
		baseURL = bcjhCfPageURLBase
	case "白菜菊花":
		baseURL = bcjhURLBase
	default:
		if util.HasPrefixIn(strings.TrimSpace(c.PretreatedMessage), "http://", "https://") {
			baseURL = strings.TrimSpace(c.PretreatedMessage)
		} else {
			baseURL = bcjhCfPageURLBase
		}
	}
	_, _ = c.Reply(fmt.Sprintf("开始导入数据, 数据源:\n%s", baseURL))
	updateStart := time.Now()
	msg := ""

	// 获取图鉴网数据
	stepStart := time.Now()
	gameData, err := requestData(baseURL + dataURI)
	if err != nil {
		logger.Error("获取图鉴网数据失败!", err)
		_, _ = c.Reply("获取图鉴网数据失败!")
		return
	}
	stepTime := time.Now().Sub(stepStart).Round(time.Millisecond).String()
	logger.Infof("获取图鉴网数据完毕, 耗时%s", stepTime)
	msg += fmt.Sprintf("获取图鉴网数据耗时%s\n", stepTime)

	// 更新数据
	// 更新厨师数据
	stepStart = time.Now()
	err = updateChefs(gameData.Chefs)
	if err != nil {
		logger.Error("更新厨师数据出错!", err)
		_, _ = c.Reply("更新厨师数据出错!")
		return
	}
	stepTime = time.Now().Sub(stepStart).Round(time.Millisecond).String()
	logger.Infof("更新厨师数据完毕, 耗时%s", stepTime)
	msg += fmt.Sprintf("更新厨师数据耗时%s\n", stepTime)

	// 更新厨具数据
	stepStart = time.Now()
	err = updateEquips(gameData.Equips)
	if err != nil {
		logger.Error("更新厨具数据出错!", err)
		_, _ = c.Reply("更新厨具数据出错!")
		return
	}
	stepTime = time.Now().Sub(stepStart).Round(time.Millisecond).String()
	logger.Infof("更新厨具数据完毕, 耗时%s", stepTime)
	msg += fmt.Sprintf("更新厨具数据耗时%s\n", stepTime)

	// 更新菜谱数据
	stepStart = time.Now()
	err = updateRecipes(gameData.Recipes, gameData.Combos)
	if err != nil {
		logger.Error("更新菜谱数据出错!", err)
		_, _ = c.Reply("更新菜谱数据出错!")
		return
	}
	stepTime = time.Now().Sub(stepStart).Round(time.Millisecond).String()
	logger.Infof("更新菜谱数据完毕, 耗时%s", stepTime)
	msg += fmt.Sprintf("更新菜谱数据耗时%s\n", stepTime)

	// 更新贵客数据
	stepStart = time.Now()
	err = updateGuests(gameData.Guests)
	if err != nil {
		logger.Error("更新贵客数据出错!", err)
		_, _ = c.Reply("更新贵客数据出错!")
		return
	}
	stepTime = time.Now().Sub(stepStart).Round(time.Millisecond).String()
	logger.Infof("更新贵客数据完毕, 耗时%s", stepTime)
	msg += fmt.Sprintf("更新贵客数据耗时%s\n", stepTime)

	// 更新食材数据
	stepStart = time.Now()
	err = updateMaterials(gameData.Materials)
	if err != nil {
		logger.Error("更新食材数据出错!", err)
		_, _ = c.Reply("更新食材数据出错!")
		return
	}
	stepTime = time.Now().Sub(stepStart).Round(time.Millisecond).String()
	logger.Infof("更新食材数据完毕, 耗时%s", stepTime)
	msg += fmt.Sprintf("更新食材数据耗时%s\n", stepTime)

	// 更新技能数据
	stepStart = time.Now()
	err = updateSkills(gameData.Skills)
	if err != nil {
		logger.Error("更新技能数据出错!", err)
		_, _ = c.Reply("更新技能数据出错!")
		return
	}
	stepTime = time.Now().Sub(stepStart).Round(time.Millisecond).String()
	logger.Infof("更新技能数据完毕, 耗时%s", stepTime)
	msg += fmt.Sprintf("更新技能数据耗时%s\n", stepTime)

	// 更新装修家具数据
	stepStart = time.Now()
	err = updateDecorations(gameData.Decorations)
	if err != nil {
		logger.Error("更新装修家具数据出错!", err)
		_, _ = c.Reply("更新装修家具数据出错!")
		return
	}
	stepTime = time.Now().Sub(stepStart).Round(time.Millisecond).String()
	logger.Infof("更新装修家具数据完毕, 耗时%s", stepTime)
	msg += fmt.Sprintf("更新装修家具数据耗时%s\n", stepTime)

	// 更新调料数据
	stepStart = time.Now()
	err = updateCondiments(gameData.Condiments)
	if err != nil {
		logger.Error("更新调料数据出错!", err)
		_, _ = c.Reply("更新调料数据出错!")
		return
	}
	stepTime = time.Now().Sub(stepStart).Round(time.Millisecond).String()
	logger.Infof("更新调料数据完毕, 耗时%s", stepTime)
	msg += fmt.Sprintf("更新调料数据耗时%s\n", stepTime)

	// 更新任务数据
	stepStart = time.Now()
	err = updateQuests(gameData.Quests)
	if err != nil {
		logger.Error("更新任务数据出错!", err)
		_, _ = c.Reply("更新任务数据出错!")
		return
	}
	stepTime = time.Now().Sub(stepStart).Round(time.Millisecond).String()
	logger.Infof("更新任务数据完毕, 耗时%s", stepTime)
	msg += fmt.Sprintf("更新任务数据耗时%s\n", stepTime)

	// 开始绘制之前清除所有缓存数据，防止数据不一致
	ClearGameDataCache()

	// 解析ImgCSS数据
	stepStart = time.Now()
	imgCSS, err := ResolvingImgCSS(baseURL + imageCSSURI)
	if err != nil {
		logger.Error("解析ImgCSS出错!", err)
		_, _ = c.Reply("解析ImgCSS出错!")
		return
	}
	stepTime = time.Now().Sub(stepStart).Round(time.Millisecond).String()
	logger.Infof("解析ImgCSS图片位置信息完毕, 耗时%s", stepTime)
	msg += fmt.Sprintf("解析ImgCSS图片位置信息耗时%s\n", stepTime)

	// 下载图鉴网厨师、菜谱、厨具图片数据
	imgResourceDir := config.AppConfig.Resource.Image
	stepStart = time.Now()
	var chefImage image.Image
	var recipeImage image.Image
	var equipImage image.Image
	var eg errgroup.Group
	eg.Go(func() error {
		chefImage, err = DownloadAndLoadImage(baseURL+chefImageURI, fmt.Sprintf("%s/chef/chef_gallery.png", imgResourceDir))
		if err != nil {
			return fmt.Errorf("下载图鉴网厨师图片出错! %v", err)
		}
		return nil
	})
	eg.Go(func() error {
		recipeImage, err = DownloadAndLoadImage(baseURL+recipeImageURI, fmt.Sprintf("%s/recipe/recipe_gallery.png", imgResourceDir))
		if err != nil {
			return fmt.Errorf("下载图鉴网菜谱图片出错! %v", err)
		}
		return nil
	})
	eg.Go(func() error {
		equipImage, err = DownloadAndLoadImage(baseURL+equipImageURI, fmt.Sprintf("%s/equip/equip_gallery.png", imgResourceDir))
		if err != nil {
			return fmt.Errorf("下载图鉴网厨具图片出错! %v", err)
		}
		return nil
	})
	err = eg.Wait()
	if err != nil {
		logger.Error(err)
		_, _ = c.Reply(err.Error())
		return
	}

	stepTime = time.Now().Sub(stepStart).Round(time.Millisecond).String()
	logger.Infof("下载图鉴网厨师、菜谱、厨具图片完毕, 耗时%s", stepTime)
	msg += fmt.Sprintf("下载图鉴网厨师、菜谱、厨具图片耗时%s\n", stepTime)

	// 绘制厨师图鉴图片数据
	eg.Go(func() error {
		stepStart := time.Now()
		chefs, err := dao.FindAllChefs()
		if err != nil {
			return fmt.Errorf("绘制厨师图鉴图片出错! %v", err)
		}
		err = GenerateAllChefsImages(chefs, chefImage, imgCSS)
		if err != nil {
			return fmt.Errorf("绘制厨师图鉴图片出错! %v", err)
		}
		stepTime = time.Now().Sub(stepStart).Round(time.Millisecond).String()
		logger.Infof("绘制厨师图鉴图片完毕, 耗时%s", stepTime)
		msg += fmt.Sprintf("绘制厨师图鉴图片耗时%s\n", stepTime)
		return nil
	})

	// 绘制菜谱图鉴图片数据
	eg.Go(func() error {
		stepStart := time.Now()
		recipes, err := dao.FindAllRecipes()
		if err != nil {
			return fmt.Errorf("绘制菜谱图鉴图片出错! %v", err)
		}
		err = GenerateAllRecipesImages(recipes, recipeImage, imgCSS)
		if err != nil {
			return fmt.Errorf("绘制菜谱图鉴图片出错! %v", err)
		}
		stepTime = time.Now().Sub(stepStart).Round(time.Millisecond).String()
		logger.Infof("绘制菜谱图鉴图片完毕, 耗时%s", stepTime)
		msg += fmt.Sprintf("绘制菜谱图鉴图片耗时%s\n", stepTime)
		return nil
	})

	// 绘制厨具图鉴图片数据
	eg.Go(func() error {
		stepStart := time.Now()
		equips, err := dao.FindAllEquips()
		if err != nil {
			return fmt.Errorf("绘制厨具图鉴图片出错! %v", err)
		}
		err = GenerateAllEquipmentsImages(equips, equipImage, imgCSS)
		if err != nil {
			return fmt.Errorf("绘制厨具图鉴图片出错! %v", err)
		}
		stepTime = time.Now().Sub(stepStart).Round(time.Millisecond).String()
		logger.Infof("绘制厨具图鉴图片完毕, 耗时%s", stepTime)
		msg += fmt.Sprintf("绘制厨具图鉴图片耗时%s\n", stepTime)
		return nil
	})
	err = eg.Wait()
	if err != nil {
		logger.Error(err)
		_, _ = c.Reply(err.Error())
		return
	}

	// 发送成功消息
	logger.Info("更新数据完毕")
	msg = fmt.Sprintf("更新数据完毕, 累计耗时%s\n", time.Now().Sub(updateStart).Round(time.Millisecond).String()) + msg
	msg = strings.TrimSuffix(msg, "\n") // 去除结尾的换行
	_, _ = c.Reply(msg)
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

// DownloadAndLoadImage 下载图鉴图片并导出 image.Image 对象
func DownloadAndLoadImage(url string, path string) (image.Image, error) {
	// 下载图片
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// 保存到文件
	out, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(out, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// 导出image实例
	img, err := png.Decode(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return img, nil
}

// 更新厨师信息
func updateChefs(chefsData []gamedata.ChefData) error {
	session := dao.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM `%s`", database.Chef{}.TableName())
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
	session := dao.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM `%s`", database.Equip{}.TableName())
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
func updateRecipes(recipesData []gamedata.RecipeData, combosData []gamedata.ComboData) error {
	session := dao.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	// 删除菜谱数据
	sql := fmt.Sprintf("DELETE FROM `%s`", database.Recipe{}.TableName())
	_, err = session.Exec(sql)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	// 删除菜谱-食材关系
	sql = fmt.Sprintf("DELETE FROM `%s`", database.RecipeMaterial{}.TableName())
	_, err = session.Exec(sql)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	// 准备后厨合成菜数据
	mIdToNameCombo := make(map[int]struct {
		Name   string
		Combos []string
	})
	for i := range recipesData {
		mIdToNameCombo[recipesData[i].RecipeId] = struct {
			Name   string
			Combos []string
		}{Name: recipesData[i].Name, Combos: []string{}}
	}
	for _, combo := range combosData {
		for _, recipeId := range combo.Recipes {
			nameComboData := mIdToNameCombo[recipeId]
			nameComboData.Combos = append(nameComboData.Combos, mIdToNameCombo[combo.RecipeId].Name)
			mIdToNameCombo[recipeId] = nameComboData
		}
	}

	// 生成菜谱数据
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
			Combo:          mIdToNameCombo[recipeData.RecipeId].Combos,
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

// 更新贵客信息
func updateGuests(guestsData []gamedata.GuestData) error {
	session := dao.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM `%s`", database.GuestGift{}.TableName())
	_, err = session.Exec(sql)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	// 从guest表中获取预设的贵客编号
	guestInfo := make([]database.Guest, 0)
	err = dao.DB.Find(&guestInfo)
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
	err = dao.DB.Find(&recipeInfo)
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
	session := dao.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM `%s`", database.Material{}.TableName())
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
	session := dao.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM `%s`", database.Skill{}.TableName())
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
	session := dao.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM `%s`", database.Decoration{}.TableName())
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
	session := dao.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM `%s`", database.Condiment{}.TableName())
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
	session := dao.DB.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf("DELETE FROM `%s`", database.Quest{}.TableName())
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

// ClearGameDataCache 清除图鉴网相关的游戏数据缓存
func ClearGameDataCache() {
	dao.ClearChefsCache()
	dao.ClearEquipsCache()
	dao.ClearRecipesCache()
	dao.ClearGuestGiftsCache()
	dao.ClearMaterialsCache()
	dao.ClearSkillsCache()
	dao.ClearQuestsCache()
}
