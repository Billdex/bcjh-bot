package messageservice

import (
	"bcjh-bot/config"
	"bcjh-bot/dao"
	"bcjh-bot/model/database"
	"bcjh-bot/model/gamedata"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func RecipeQuery(c *scheduler.Context) {
	order := "单时间"
	page := 1
	var note string
	recipes, err := dao.FindAllRecipes()
	if err != nil {
		logger.Error("查询菜谱数据出错!", err)
		_, _ = c.Reply(e.SystemErrorNote)
	}
	args := strings.Split(c.PretreatedMessage, " ")
	argCount := 0
	for _, arg := range args {
		if arg == "" {
			continue
		}
		switch arg {
		case "图鉴序", "时间", "单时间", "总时间", "单价", "售价", "金币效率", "耗材效率", "稀有度":
			order = arg
		case "1火", "1星", "一火", "一星":
			recipes, note = filterRecipesByRarity(recipes, 1, true)
		case "2火", "2星", "二火", "二星", "两火", "两星":
			recipes, note = filterRecipesByRarity(recipes, 2, true)
		case "3火", "3星", "三火", "三星":
			recipes, note = filterRecipesByRarity(recipes, 3, true)
		case "4火", "4星", "四火", "四星":
			recipes, note = filterRecipesByRarity(recipes, 4, true)
		case "5火", "5星", "五火", "五星":
			recipes, note = filterRecipesByRarity(recipes, 5, true)
		case "仅1火", "仅1星", "仅一火", "仅一星":
			recipes, note = filterRecipesByRarity(recipes, 1, false)
		case "仅2火", "仅2星", "仅二火", "仅二星", "仅两火", "仅两星":
			recipes, note = filterRecipesByRarity(recipes, 2, false)
		case "仅3火", "仅3星", "仅三火", "仅三星":
			recipes, note = filterRecipesByRarity(recipes, 3, false)
		case "仅4火", "仅4星", "仅四火", "仅四星":
			recipes, note = filterRecipesByRarity(recipes, 4, false)
		case "仅5火", "仅5星", "仅五火", "仅五星":
			recipes, note = filterRecipesByRarity(recipes, 5, false)
		case "炒技法", "烤技法", "煮技法", "蒸技法", "炸技法", "切技法":
			recipes, note = filterRecipesBySkill(recipes, strings.TrimSuffix(arg, "技法"))
		case "甜味", "酸味", "辣味", "咸味", "苦味", "鲜味":
			recipes, note = filterRecipesByCondiment(recipes, strings.TrimSuffix(arg, "味"))
		default:
			if pre, has := util.WhatPrefixIn(arg, "食材", "材料"); has {
				materials := strings.Split(strings.TrimLeft(arg, pre), "-")
				recipes, note = filterRecipesByMaterials(recipes, materials)
			} else if pre, has := util.WhatPrefixIn(arg, "技法"); has {
				skills := strings.Split(strings.TrimLeft(arg, pre), "-")
				recipes, note = filterRecipesBySkills(recipes, skills)
			} else if pre, has := util.WhatPrefixIn(arg, "贵客", "稀有客人", "客人", "贵宾", "宾客", "稀客"); has {
				guests := strings.Split(strings.TrimLeft(arg, pre), "-")
				recipes, note = filterRecipesByGuests(recipes, guests[1:])
			} else if pre, has := util.WhatPrefixIn(arg, "符文", "礼物"); has {
				antiques := strings.Split(strings.TrimLeft(arg, pre), "-")
				if len(antiques) > 1 {
					recipes, note = filterRecipesByAntique(recipes, antiques[1])
				} else {
					recipes, note = filterRecipesByAntique(recipes, antiques[0])
				}
			} else if util.HasPrefixIn(arg, "神级符文", "神级奖励") {
				antiques := strings.Split(arg, "-")
				if len(antiques) > 1 {
					recipes, note = filterRecipesByUpgradeAntique(recipes, antiques[1])
				}
			} else if util.HasPrefixIn(arg, "来源") {
				origins := strings.Split(arg, "-")
				if len(origins) > 1 {
					recipes, note = filterRecipesByOrigin(recipes, origins[1])
				}
			} else if util.HasPrefixIn(arg, "调料", "调味", "味道") {
				condiments := strings.Split(arg, "-")
				if len(condiments) > 1 {
					recipes, note = filterRecipesByCondiment(recipes, condiments[1])
				}
			} else if util.HasPrefixIn(arg, "$", "＄", "￥") {
				num, err := strconv.Atoi(arg[1:])
				if err != nil {
					note = "单价筛选参数有误"
				} else {
					recipes, note = filterRecipesByPrice(recipes, num)
				}
			} else if util.HasPrefixIn(arg, "p", "P") {
				pageNum, err := strconv.Atoi(strings.Trim(arg[1:], "-"))
				if err != nil {
					note = "分页参数有误"
				} else {
					if pageNum > 0 {
						page = pageNum
					}
				}
			} else {
				recipes, note = filterRecipesByName(recipes, arg)
			}
		}

		if note != "" {
			logger.Info("菜谱查询失败:", note)
			_, _ = c.Reply(note)
			return
		}
		argCount++
	}

	// 对菜谱查询结果排序
	recipes, note = orderRecipes(recipes, order)
	if note != "" {
		logger.Info("菜谱排序失败:", note)
		_, _ = c.Reply(note)
		return
	}
	// 根据结果翻页并发送消息
	msg := echoRecipesMessage(recipes, order, page, c.GetMessageType() == onebot.MessageTypePrivate)
	logger.Info("发送菜谱查询结果:", msg)
	_, _ = c.Reply(msg)
}

// 根据具体稀有度筛选菜谱, gte 参数为 true 时筛选大于等于，否则筛选仅等于
func filterRecipesByRarity(recipes []database.Recipe, rarity int, gte bool) ([]database.Recipe, string) {
	if len(recipes) == 0 {
		return recipes, ""
	}
	result := make([]database.Recipe, 0)
	for i := range recipes {
		if (gte && recipes[i].Rarity >= rarity) || recipes[i].Rarity == rarity {
			result = append(result, recipes[i])
		}
	}
	return result, ""
}

// 根据食材筛选菜谱
func filterRecipesByMaterial(recipes []database.Recipe, material string) ([]database.Recipe, string) {
	if len(recipes) == 0 || material == "" {
		return recipes, ""
	}
	var origins []string
	result := make([]database.Recipe, 0)
	// 符合下列特征的关键词视为根据来源筛选食材
	switch material {
	case "鱼类", "水产", "水产类", "海鲜", "海鲜类", "池塘":
		origins = []string{"池塘"}
	case "蔬菜", "蔬菜类", "菜类":
		origins = []string{"菜棚", "菜地", "森林"}
	case "肉类":
		origins = []string{"牧场", "鸡舍", "猪圈"}
	case "面类", "加工类", "作坊":
		origins = []string{"作坊"}
	}
	if len(origins) > 0 {
		for i := range recipes {
			if recipes[i].HasMaterialOrigins(origins) {
				result = append(result, recipes[i])
			}
		}
	} else {
		// 查出所有食材，假设存在完全匹配的则只使用该食材筛选。
		materials, err := dao.SearchMaterialsWithName(material)
		if err != nil {
			logger.Error("根据名称搜索食材失败", err)
			return nil, e.SystemErrorNote
		}
		if len(materials) == 0 {
			return nil, fmt.Sprintf("厨师长说没有用%s做过菜", material)
		}
		if len(materials) > 1 {
			for i := range materials {
				if materials[i].Name == material {
					materials = []database.Material{materials[i]}
					break
				}
			}
		}
		materialNames := make([]string, 0, len(materials))
		for i := range materials {
			materialNames = append(materialNames, materials[i].Name)
		}
		for i := range recipes {
			if recipes[i].UsedMaterials(materialNames) {
				result = append(result, recipes[i])
			}
		}
	}
	return result, ""
}

// 根据食材列表筛选菜谱
func filterRecipesByMaterials(recipes []database.Recipe, materials []string) ([]database.Recipe, string) {
	if len(materials) == 0 {
		return nil, "你想筛选什么食材呀? 食材参数格式为「食材-食材名」"
	}
	if len(recipes) == 0 {
		return recipes, ""
	}
	result := recipes
	var note string
	for _, material := range materials {
		result, note = filterRecipesByMaterial(result, material)
		if note != "" {
			return nil, note
		}
	}

	return result, ""
}

// 根据技法筛选菜谱
func filterRecipesBySkill(recipes []database.Recipe, skill string) ([]database.Recipe, string) {
	if len(recipes) == 0 || skill == "" {
		return recipes, ""
	}
	result := make([]database.Recipe, 0)
	for _, recipe := range recipes {
		need, err := recipe.NeedSkill(skill)
		if err != nil {
			return nil, err.Error()
		} else if need {
			result = append(result, recipe)
		}
	}
	return result, ""
}

// 根据技法列表筛选菜谱
func filterRecipesBySkills(recipes []database.Recipe, skills []string) ([]database.Recipe, string) {
	if len(skills) == 0 {
		return nil, "你想筛选什么技法呀? 技法参数格式为「技法-技法名」或「X技法」"
	}
	if len(recipes) == 0 {
		return recipes, ""
	}
	result := recipes
	var note string
	for _, skill := range skills {
		if skill == "" {
			continue
		} else {
			result, note = filterRecipesBySkill(result, skill)
			if note != "" {
				return nil, note
			}
		}
	}
	return result, ""
}

// 根据贵客筛选菜谱
func filterRecipeByGuest(recipes []database.Recipe, guest string) ([]database.Recipe, string) {
	if len(recipes) == 0 || guest == "" {
		return recipes, ""
	}
	result := make([]database.Recipe, 0)
	for i := range recipes {
		if recipes[i].HasGuest(guest) {
			result = append(result, recipes[i])
		}
	}
	return result, ""
}

// 根据贵客列表查询菜谱
func filterRecipesByGuests(recipes []database.Recipe, guests []string) ([]database.Recipe, string) {
	if len(guests) == 0 {
		return nil, "你想筛选哪位贵客呀? 贵客参数格式为「贵客-贵客名」"
	}
	result := recipes
	var note string
	for _, guest := range guests {
		result, note = filterRecipeByGuest(result, guest)
		if note != "" {
			return nil, note
		}
	}

	return result, ""
}

// 根据符文礼物查询菜谱
func filterRecipesByAntique(recipes []database.Recipe, antique string) ([]database.Recipe, string) {
	if len(recipes) == 0 {
		return recipes, ""
	}
	result := make([]database.Recipe, 0)
	for i := range recipes {
		if recipes[i].HasAntique(antique) {
			result = append(result, recipes[i])
		}
	}
	return result, ""
}

// 根据菜谱神级符文查询菜谱
func filterRecipesByUpgradeAntique(recipes []database.Recipe, antique string) ([]database.Recipe, string) {
	if len(recipes) == 0 || antique == "" {
		return recipes, ""
	}
	result := make([]database.Recipe, 0)
	pattern := strings.ReplaceAll(antique, "%", ".*")
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, "神级符文查询格式有误"
	}
	for i := range recipes {
		if re.MatchString(recipes[i].Gift) {
			result = append(result, recipes[i])
		}
	}
	return result, ""
}

// 根据来源筛选菜谱
func filterRecipesByOrigin(recipes []database.Recipe, origin string) ([]database.Recipe, string) {
	if len(recipes) == 0 {
		return recipes, ""
	}
	result := make([]database.Recipe, 0)
	pattern := strings.ReplaceAll(origin, "%", ".*")
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, "菜谱来源查询格式有误"
	}
	for i := range recipes {
		if re.MatchString(recipes[i].Origin) {
			result = append(result, recipes[i])
		}
	}
	return result, ""
}

// 根据调料筛选菜谱
func filterRecipesByCondiment(recipes []database.Recipe, condiment string) ([]database.Recipe, string) {
	if len(recipes) == 0 || condiment == "" {
		return recipes, ""
	}
	result := make([]database.Recipe, 0)
	switch condiment {
	case "甜":
		condiment = "Sweet"
	case "酸":
		condiment = "Sour"
	case "辣":
		condiment = "Spicy"
	case "咸":
		condiment = "Salty"
	case "苦":
		condiment = "Bitter"
	case "鲜":
		condiment = "Tasty"
	default:
		return nil, fmt.Sprintf("%s是什么味道呀", condiment)
	}
	for i := range recipes {
		if recipes[i].Condiment == condiment {
			result = append(result, recipes[i])
		}
	}
	return result, ""
}

// 根据名字或图鉴ID筛选菜谱
func filterRecipesByName(recipes []database.Recipe, name string) ([]database.Recipe, string) {
	result := make([]database.Recipe, 0)
	numId, err := strconv.Atoi(name)
	if err != nil {
		pattern := strings.ReplaceAll(name, "%", ".*")
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, "查询格式有误"
		}
		for i := range recipes {
			if recipes[i].Name == name {
				return []database.Recipe{recipes[i]}, ""
			}
			if re.MatchString(recipes[i].Name) {
				result = append(result, recipes[i])
			}
		}
	} else {
		for i := range recipes {
			if recipes[i].RecipeId == numId {
				result = append(result, recipes[i])
			}
		}
	}
	return result, ""
}

// 根据菜谱单价筛选菜谱
func filterRecipesByPrice(recipes []database.Recipe, price int) ([]database.Recipe, string) {
	result := make([]database.Recipe, 0)
	for i := range recipes {
		if recipes[i].Price >= price {
			result = append(result, recipes[i])
		}
	}
	return result, ""
}

// 根据排序参数排序菜谱
func orderRecipes(recipes []database.Recipe, order string) ([]database.Recipe, string) {
	if len(recipes) == 0 {
		return recipes, ""
	}
	switch order {
	case "图鉴序":
		sort.Slice(recipes, func(i, j int) bool {
			return recipes[i].RecipeId < recipes[j].RecipeId
		})
	case "单时间":
		sort.Slice(recipes, func(i, j int) bool {
			return recipes[i].Time == recipes[j].Time && recipes[i].RecipeId < recipes[j].RecipeId ||
				recipes[i].Time < recipes[j].Time
		})
	case "总时间":
		sort.Slice(recipes, func(i, j int) bool {
			return recipes[i].TotalTime == recipes[j].TotalTime && recipes[i].RecipeId < recipes[j].RecipeId ||
				recipes[i].TotalTime < recipes[j].TotalTime
		})
	case "单价", "售价":
		sort.Slice(recipes, func(i, j int) bool {
			return recipes[i].Price == recipes[j].Price && recipes[i].RecipeId < recipes[j].RecipeId ||
				recipes[i].Price > recipes[j].Price
		})
	case "金币效率":
		sort.Slice(recipes, func(i, j int) bool {
			return recipes[i].GoldEfficiency == recipes[j].GoldEfficiency && recipes[i].RecipeId < recipes[j].RecipeId ||
				recipes[i].GoldEfficiency > recipes[j].GoldEfficiency
		})
	case "耗材效率":
		sort.Slice(recipes, func(i, j int) bool {
			return recipes[i].MaterialEfficiency == recipes[j].MaterialEfficiency && recipes[i].RecipeId < recipes[j].RecipeId ||
				recipes[i].MaterialEfficiency > recipes[j].MaterialEfficiency
		})
	case "稀有度":
		sort.Slice(recipes, func(i, j int) bool {
			return recipes[i].Rarity == recipes[j].Rarity && recipes[i].RecipeId < recipes[j].RecipeId ||
				recipes[i].Rarity > recipes[j].Rarity
		})
	default:
		return nil, "排序参数有误"
	}
	return recipes, ""
}

// 输出单菜谱消息数据
func echoRecipeMessage(recipe database.Recipe) string {
	// 尝试寻找图片文件，未找到则按照文字格式发送
	resourceImageDir := config.AppConfig.Resource.Image + "/recipe"
	imagePath := fmt.Sprintf("%s/recipe_%s_%s.png", resourceImageDir, recipe.GalleryId, strings.ReplaceAll(recipe.Name, " ", "_"))
	logger.Debug("imagePath:", imagePath)
	var msg string
	if has, err := util.PathExists(imagePath); has {
		msg = onebot.GetCQImage(imagePath, "file")
	} else {
		if err != nil {
			logger.Debugf("无法确定文件是否存在!", err)
		}
		logger.Info("未找到菜谱图鉴图片, 以文字格式发送数据")
		// 菜谱所需技法数据
		recipeSkill := ""
		if recipe.Stirfry > 0 {
			recipeSkill += fmt.Sprintf("炒: %d  ", recipe.Stirfry)
		}
		if recipe.Bake > 0 {
			recipeSkill += fmt.Sprintf("烤: %d  ", recipe.Bake)
		}
		if recipe.Boil > 0 {
			recipeSkill += fmt.Sprintf("煮: %d  ", recipe.Boil)
		}
		if recipe.Steam > 0 {
			recipeSkill += fmt.Sprintf("蒸: %d  ", recipe.Steam)
		}
		if recipe.Fry > 0 {
			recipeSkill += fmt.Sprintf("炸: %d  ", recipe.Fry)
		}
		if recipe.Cut > 0 {
			recipeSkill += fmt.Sprintf("切: %d  ", recipe.Cut)
		}
		// 食材数据
		materials := make([]string, 0, len(recipe.Materials))
		for _, material := range recipe.Materials {
			materials = append(materials, fmt.Sprintf("%s*%d", material.Material.Name, material.Quantity))
		}
		// 贵客礼物数据
		giftInfo := make([]string, 0, len(recipe.GuestGifts))
		for _, guestGift := range recipe.GuestGifts {
			giftInfo = append(giftInfo, fmt.Sprintf("%s-%s", guestGift.GuestName, guestGift.Antique))
		}
		// 升阶贵客数据
		guests := ""
		if len(recipe.Guests) > 0 && recipe.Guests[0] != "" {
			guests += fmt.Sprintf("优-%s, ", recipe.Guests[0])
		} else {
			guests += fmt.Sprintf("优-未知, ")
		}
		if len(recipe.Guests) > 1 && recipe.Guests[1] != "" {
			guests += fmt.Sprintf("特-%s, ", recipe.Guests[1])
		} else {
			guests += fmt.Sprintf("特-未知, ")
		}
		if len(recipe.Guests) > 2 && recipe.Guests[2] != "" {
			guests += fmt.Sprintf("神-%s", recipe.Guests[2])
		} else {
			guests += fmt.Sprintf("神-未知")
		}
		msg += fmt.Sprintf("%s %s %s\n", recipe.GalleryId, recipe.Name, recipe.FormatRarity())
		msg += fmt.Sprintf("💰: %d(%d) --- %d/h\n", recipe.Price, recipe.Price+recipe.ExPrice, recipe.GoldEfficiency)
		msg += fmt.Sprintf("来源: %s\n", recipe.Origin)
		msg += fmt.Sprintf("单时间: %s\n", util.FormatSecondToString(recipe.Time))
		msg += fmt.Sprintf("总时间: %s (%d份)\n", util.FormatSecondToString(recipe.Time*recipe.Limit), recipe.Limit)
		msg += fmt.Sprintf("技法: %s\n", recipeSkill)
		msg += fmt.Sprintf("食材: %s\n", strings.Join(materials, ","))
		msg += fmt.Sprintf("耗材效率: %d/h\n", recipe.MaterialEfficiency)
		msg += fmt.Sprintf("可解锁: %s\n", recipe.Unlock)
		msg += fmt.Sprintf("可合成: %s\n", strings.Join(recipe.Combo, ","))
		msg += fmt.Sprintf("神级符文: %s\n", recipe.Gift)
		msg += fmt.Sprintf("贵客礼物: %s\n", giftInfo)
		msg += fmt.Sprintf("升阶贵客: %s", guests)
	}
	return msg
}

// 根据排序规则与分页参数，返回菜谱列表消息数据
func echoRecipesMessage(recipes []database.Recipe, order string, page int, private bool) string {
	if len(recipes) == 0 {
		return "本店没有相关的菜呢!"
	} else if len(recipes) == 1 {
		return echoRecipeMessage(recipes[0])
	} else {
		var msg string
		listLength := config.AppConfig.Bot.GroupMsgMaxLen
		if private {
			listLength = config.AppConfig.Bot.PrivateMsgMaxLen
		}
		maxPage := (len(recipes)-1)/listLength + 1
		if page > maxPage {
			page = maxPage
		}
		if len(recipes) > listLength {
			msg += fmt.Sprintf("这里有你想点的菜吗 (%d/%d)\n", page, maxPage)
		} else {
			msg += "这里有你想点的菜吗\n"
		}
		for i := (page - 1) * listLength; i < page*listLength && i < len(recipes); i++ {
			orderInfo := getRecipeInfoWithOrder(recipes[i], order)
			msg += fmt.Sprintf("%s %s %s", recipes[i].GalleryId, recipes[i].Name, orderInfo)
			if i < page*listLength-1 && i < len(recipes)-1 {
				msg += "\n"
			}
		}
		if page < maxPage {
			msg += "\n......"
		}
		return msg
	}
}

// 根据排序参数获取菜谱需要输出的信息
func getRecipeInfoWithOrder(recipe database.Recipe, order string) string {
	switch order {
	case "单时间":
		return util.FormatSecondToString(recipe.Time)
	case "总时间":
		return util.FormatSecondToString(recipe.Time * recipe.Limit)
	case "单价", "售价":
		return fmt.Sprintf("💰%d", recipe.Price)
	case "金币效率":
		return fmt.Sprintf("💰%d/h", recipe.GoldEfficiency)
	case "耗材效率":
		return fmt.Sprintf("🥗%d/h", recipe.MaterialEfficiency)
	case "稀有度":
		return recipe.FormatRarity()
	default:
		return ""
	}
}

func GenerateRecipeImage(recipe database.RecipeData, font *truetype.Font, bgImg image.Image, rarityImg image.Image, condimentImg image.Image) (image.Image, error) {
	titleSize := 48 // 标题字体尺寸
	fontSize := 32  // 内容字体尺寸

	img := image.NewRGBA(image.Rect(0, 0, 800, 800))
	draw.Draw(img, img.Bounds(), bgImg, bgImg.Bounds().Min, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	fontColor := color.RGBA{A: 255}
	c.SetSrc(image.NewUniform(fontColor))

	// 输出图鉴ID与菜谱名
	c.SetFontSize(float64(titleSize))
	_, err := c.DrawString(fmt.Sprintf("%s %s", recipe.GalleryId, recipe.Name), freetype.Pt(25, 30+titleSize))
	if err != nil {
		return nil, err
	}

	// 输出菜谱图鉴图片
	width := recipe.Avatar.Bounds().Dx()
	height := recipe.Avatar.Bounds().Dy()
	draw.Draw(img,
		image.Rect(70+200/2-width/2, 100+200/2-height/2, 70+200/2+width/2, 100+200/2+height/2),
		recipe.Avatar,
		image.Point{},
		draw.Over)

	// 输出稀有度
	draw.Draw(img,
		image.Rect(50, 310, 50+240, 310+44),
		rarityImg,
		image.Point{},
		draw.Over)

	// 输出单价信息
	fontColor = color.RGBA{R: 45, G: 45, B: 45, A: 255}
	c.SetSrc(image.NewUniform(fontColor))
	c.SetFontSize(float64(fontSize))
	_, err = c.DrawString(fmt.Sprintf("%d", recipe.Price), freetype.Pt(435, 105+fontSize))
	if err != nil {
		return nil, err
	}
	fontColor = color.RGBA{R: 120, G: 120, B: 120, A: 255}
	c.SetSrc(image.NewUniform(fontColor))
	_, err = c.DrawString(fmt.Sprintf("+%d", recipe.ExPrice), freetype.Pt(515, 105+fontSize))
	if err != nil {
		return nil, err
	}
	fontColor = color.RGBA{R: 45, G: 45, B: 45, A: 255}
	c.SetSrc(image.NewUniform(fontColor))
	// 输出金币效率
	_, err = c.DrawString(fmt.Sprintf("%d / h", recipe.GoldEfficiency), freetype.Pt(626, 105+fontSize))
	if err != nil {
		return nil, err
	}
	// 输出份数
	_, err = c.DrawString(fmt.Sprintf("%d 份 / 组", recipe.Limit), freetype.Pt(627, 175+fontSize))
	if err != nil {
		return nil, err
	}
	// 输出单份制作时间
	_, err = c.DrawString(fmt.Sprintf("%s", util.FormatSecondToString(recipe.Time)), freetype.Pt(435, 175+fontSize))
	if err != nil {
		return nil, err
	}
	// 输出整组制作总时间
	_, err = c.DrawString(fmt.Sprintf("%s", util.FormatSecondToString(recipe.TotalTime)), freetype.Pt(435, 245+fontSize))
	if err != nil {
		return nil, err
	}
	// 输出调料
	draw.Draw(img,
		image.Rect(370, 310, 370+61, 310+53),
		condimentImg,
		image.Point{},
		draw.Over)

	// 输出技法
	for i, skill := range recipe.Skills {
		draw.Draw(img,
			image.Rect(460+i*170, 310, 460+i*170+140, 310+53),
			skill.Image,
			image.Point{},
			draw.Over)
		_, err = c.DrawString(fmt.Sprintf("%d", skill.Value), freetype.Pt(525+i*170, 315+fontSize))
		if err != nil {
			return nil, err
		}
	}
	// 输出食材
	materials := make([]string, len(recipe.Materials))
	for i, material := range recipe.Materials {
		materials[i] = fmt.Sprintf("%s*%d", material.Material.Name, material.Quantity)
	}
	_, err = c.DrawString(fmt.Sprintf("%s", strings.Join(materials, " ")), freetype.Pt(170, 388+fontSize))
	if err != nil {
		return nil, err
	}
	// 输出贵客礼物
	giftInfo := "无"
	if len(recipe.GuestGifts) != 0 {
		gifts := make([]string, len(recipe.GuestGifts))
		for i, gift := range recipe.GuestGifts {
			gifts[i] = fmt.Sprintf("%s-%s", gift.GuestName, gift.Antique)
		}
		giftInfo = strings.Join(gifts, ", ")
	}
	_, err = c.DrawString(fmt.Sprintf("%s", giftInfo), freetype.Pt(170, 448+fontSize))
	if err != nil {
		return nil, err
	}

	// 输出来源
	_, err = c.DrawString(fmt.Sprintf("%s", recipe.Origin), freetype.Pt(170, 508+fontSize))
	if err != nil {
		return nil, err
	}

	// 输出升阶贵客
	for p, guest := range recipe.Guests {
		if guest == "" {
			guest = "未知"
		}
		_, err = c.DrawString(fmt.Sprintf("%s", guest), freetype.Pt(85, 620+p*54+fontSize))
		if err != nil {
			return nil, err
		}
	}

	// 输出耗材效率
	_, err = c.DrawString(fmt.Sprintf("%d / h", recipe.MaterialEfficiency), freetype.Pt(525, 576+fontSize))
	if err != nil {
		return nil, err
	}

	//输出神级奖励
	reward := recipe.Gift
	if recipe.Gift == "-" {
		reward = recipe.Unlock
	}
	_, err = c.DrawString(fmt.Sprintf("%s", reward), freetype.Pt(525, 655+fontSize))
	if err != nil {
		return nil, err
	}

	// 输出可合成的后厨菜数据
	var combo string
	if len(recipe.Combo) == 0 {
		combo = "无"
	} else {
		combo = strings.Join(recipe.Combo, ",")
	}
	_, err = c.DrawString(fmt.Sprintf("%s", combo), freetype.Pt(490, 734+fontSize))
	if err != nil {
		return nil, err
	}

	return img, nil
}

func GenerateAllRecipesImages(recipes []database.Recipe, galleryImg image.Image, imgCSS *gamedata.ImgCSS) error {
	magnification := 5 // 截取的图像相比图鉴网原始图片的放大倍数

	// 载入字体文件
	font, err := util.LoadFont(config.AppConfig.Resource.Font)
	if err != nil {
		return fmt.Errorf("载入字体文件失败 %v", err)
	}

	resourceImgDir := config.AppConfig.Resource.Image
	commonImgPath := resourceImgDir + "/common"
	recipeImgPath := resourceImgDir + "/recipe"

	// 放大菜谱图鉴图像
	logger.Debugf("菜谱图片原始尺寸:%d*%d", galleryImg.Bounds().Dx(), galleryImg.Bounds().Dy())
	galleryImg = resize.Resize(
		uint(galleryImg.Bounds().Dx()*magnification/2.0),
		uint(galleryImg.Bounds().Dy()*magnification/2.0),
		galleryImg, resize.MitchellNetravali)

	// 加载背景图片
	bgImg, err := util.LoadPngImageFile(fmt.Sprintf("%s/recipe_bg.png", recipeImgPath))
	if err != nil {
		return fmt.Errorf("载入菜谱背景图片失败 %v", err)
	}

	// 载入稀有度图片
	mRarityImages := make(map[int]image.Image)
	for _, rarity := range []int{1, 2, 3, 4, 5} {
		img, err := util.LoadPngImageFile(fmt.Sprintf("%s/rarity_%d.png", commonImgPath, rarity))
		if err != nil {
			return fmt.Errorf("载入稀有度图标失败 %v", err)
		}
		mRarityImages[rarity] = img
	}

	// 载入技法数值图片
	mSkillImages := make(map[string]image.Image)
	for _, skill := range []string{"stirfry", "bake", "boil", "steam", "fry", "cut"} {
		img, err := util.LoadPngImageFile(fmt.Sprintf("%s/icon_%s_value.png", commonImgPath, skill))
		if err != nil {
			return fmt.Errorf("载入技法图标失败 %v", err)
		}
		mSkillImages[skill] = img
	}

	// 载入调料属性图片
	mCondimentImages := make(map[string]image.Image)
	for _, condiment := range []string{"sweet", "sour", "spicy", "salty", "bitter", "tasty"} {
		img, err := util.LoadPngImageFile(fmt.Sprintf("%s/icon_%s.png", commonImgPath, condiment))
		if err != nil {
			return fmt.Errorf("载入调料图标失败 %v", err)
		}
		mCondimentImages[condiment] = img
	}

	for _, recipe := range recipes {
		// 载入与计算菜谱信息
		recipeImgInfo := imgCSS.RecipeImg[recipe.RecipeId]
		avatarStartX := recipeImgInfo.X * magnification
		avatarStartY := recipeImgInfo.Y * magnification
		avatarWidth := recipeImgInfo.Width * magnification
		avatarHeight := recipeImgInfo.Height * magnification
		avatar := image.NewRGBA(image.Rect(0, 0, avatarWidth, avatarHeight))
		draw.Draw(avatar,
			image.Rect(0, 0, avatarWidth, avatarHeight),
			galleryImg,
			image.Point{X: avatarStartX, Y: avatarStartY},
			draw.Over)

		skills := make([]database.RecipeSkillData, 0)
		for skill, value := range recipe.GetSkillValueMap() {
			if value != 0 {
				skills = append(skills, database.RecipeSkillData{
					Type:  skill,
					Value: value,
					Image: mSkillImages[skill],
				})
			}
		}

		recipeData := database.RecipeData{
			Recipe: recipe,
			Avatar: avatar,
			Skills: skills,
		}

		img, err := GenerateRecipeImage(recipeData, font, bgImg, mRarityImages[recipe.Rarity], mCondimentImages[strings.ToLower(recipe.Condiment)])
		if err != nil {
			return fmt.Errorf("绘制菜谱 %s 的数据出错 %v", recipe.GalleryId, err)
		}

		// 以PNG格式保存文件
		err = util.SavePngImage(fmt.Sprintf("%s/recipe_%s_%s.png", recipeImgPath, recipe.GalleryId, strings.ReplaceAll(recipe.Name, " ", "_")), img)
		if err != nil {
			return fmt.Errorf("保存菜谱 %s 图鉴图片出错 %v", recipe.GalleryId, err)
		}
	}
	return nil
}
