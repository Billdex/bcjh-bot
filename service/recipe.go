package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/config"
	"bcjh-bot/model/database"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"fmt"
	"github.com/golang/freetype"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func RecipeQuery(c *onebot.Context, args []string) {
	logger.Info("菜谱查询, 参数:", args)

	if len(args) == 0 {
		err := bot.SendMessage(c, recipeHelp())
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
		return
	}

	order := "图鉴序"
	page := 1
	var note string
	recipes := make([]database.Recipe, 0)
	err := database.DB.Find(&recipes)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_ = bot.SendMessage(c, util.SystemErrorNote)
	}
	argCount := 0
	for _, arg := range args {
		if arg == "" {
			continue
		}
		switch arg {
		case "图鉴序", "单时间", "总时间", "单价", "售价", "金币效率", "耗材效率":
			order = arg
		case "1火", "1星", "一火", "一星":
			recipes, note = filterRecipesByRarity(recipes, 1)
		case "2火", "2星", "二火", "二星", "两火", "两星":
			recipes, note = filterRecipesByRarity(recipes, 2)
		case "3火", "3星", "三火", "三星":
			recipes, note = filterRecipesByRarity(recipes, 3)
		case "4火", "4星", "四火", "四星":
			recipes, note = filterRecipesByRarity(recipes, 4)
		case "5火", "5星", "五火", "五星":
			recipes, note = filterRecipesByRarity(recipes, 5)
		case "炒技法", "烤技法", "煮技法", "蒸技法", "炸技法", "切技法":
			recipes, note = filterRecipesBySkill(recipes, strings.TrimSuffix(arg, "技法"))
		case "甜味", "酸味", "辣味", "咸味", "苦味", "鲜味":
			recipes, note = filterRecipesByCondiment(recipes, strings.TrimSuffix(arg, "味"))
		default:
			if util.HasPrefixIn(arg, "食材", "材料") {
				materials := strings.Split(arg, util.ArgsConnectCharacter)
				recipes, note = filterRecipesByMaterials(recipes, materials[1:])
			} else if util.HasPrefixIn(arg, "技法") {
				skills := strings.Split(arg, util.ArgsConnectCharacter)
				recipes, note = filterRecipesBySkills(recipes, skills[1:])
			} else if util.HasPrefixIn(arg, "贵客", "稀有客人") {
				guests := strings.Split(arg, util.ArgsConnectCharacter)
				recipes, note = filterRecipesByGuests(recipes, guests[1:])
			} else if util.HasPrefixIn(arg, "符文", "礼物") {
				antiques := strings.Split(arg, util.ArgsConnectCharacter)
				if len(antiques) > 1 {
					recipes, note = filterRecipesByAntique(recipes, antiques[1])
				}
			} else if util.HasPrefixIn(arg, "来源") {
				origins := strings.Split(arg, util.ArgsConnectCharacter)
				if len(origins) > 1 {
					recipes, note = filterRecipesByOrigin(recipes, origins[1])
				}
			} else if util.HasPrefixIn(arg, "调料", "调味", "味道") {
				condiments := strings.Split(arg, util.ArgsConnectCharacter)
				if len(condiments) > 1 {
					recipes, note = filterRecipesByCondiment(recipes, condiments[1])
				}
			} else if util.HasPrefixIn(arg, "$") {
				num, err := strconv.Atoi(arg[1:])
				if err != nil {
					note = "单价筛选参数有误"
				} else {
					recipes, note = filterRecipesByPrice(recipes, num)
				}
			} else if util.HasPrefixIn(arg, "p", "P") {
				pageNum, err := strconv.Atoi(arg[1:])
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
			_ = bot.SendMessage(c, note)
			return
		}
		argCount++
	}

	if argCount == 0 {
		_ = bot.SendMessage(c, recipeHelp())
		return
	}
	// 对菜谱查询结果排序
	recipes, note = orderRecipes(recipes, order)
	if note != "" {
		logger.Info("菜谱查询失败:", note)
		_ = bot.SendMessage(c, note)
		return
	}
	// 根据结果翻页并发送消息
	msg := echoRecipesMessage(recipes, order, page, c.MessageType == util.OneBotMessagePrivate)
	logger.Info("发送菜谱查询结果:", msg)
	err = bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}

// 根据稀有度筛选菜谱
func filterRecipesByRarity(recipes []database.Recipe, rarity int) ([]database.Recipe, string) {
	if len(recipes) == 0 {
		return recipes, ""
	}
	result := make([]database.Recipe, 0)
	for i, _ := range recipes {
		if recipes[i].Rarity >= rarity {
			result = append(result, recipes[i])
		}
	}
	return result, ""
}

// 根据食材筛选菜谱
func filterRecipesByMaterial(recipes []database.Recipe, material string) ([]database.Recipe, string) {
	if len(recipes) == 0 {
		return recipes, ""
	}
	result := make([]database.Recipe, 0)
	// 将所有菜谱信息存入recipeMap
	recipeMap := make(map[string]database.Recipe)
	for _, recipe := range recipes {
		recipeMap[recipe.GalleryId] = recipe
	}
	// 根据食材名或食材类型找出对应的菜谱
	dbMaterials := make([]database.Material, 0)
	var materialOrigin []string
	switch material {
	case "鱼类", "水产", "水产类", "海鲜", "海鲜类", "池塘":
		materialOrigin = []string{"池塘"}
	case "蔬菜", "蔬菜类", "菜类":
		materialOrigin = []string{"菜棚", "菜地", "森林"}
	case "肉类":
		materialOrigin = []string{"牧场", "鸡舍", "猪圈"}
	case "面类", "加工类", "作坊":
		materialOrigin = []string{"作坊"}
	default:
		materialOrigin = []string{}
	}
	if len(materialOrigin) > 0 {
		err := database.DB.In("origin", materialOrigin).Find(&dbMaterials)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return nil, util.SystemErrorNote
		}
	} else {
		err := database.DB.Where("name like ?", "%"+material+"%").Find(&dbMaterials)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return nil, util.SystemErrorNote
		}
		if len(dbMaterials) == 0 {
			return nil, fmt.Sprintf("厨师长说没有用%s做过菜", material)
		}
	}
	// 找出符合食材要求的菜谱图鉴id
	materialsId := make([]int, 0)
	for _, dbMaterial := range dbMaterials {
		materialsId = append(materialsId, dbMaterial.MaterialId)
	}
	recipeMaterials := make([]database.RecipeMaterial, 0)
	err := database.DB.In("material_id", materialsId).Find(&recipeMaterials)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		return nil, util.SystemErrorNote
	}
	// 从recipeMap中选出符合要求的菜
	newRecipeMap := make(map[string]database.Recipe)
	for _, recipeMaterial := range recipeMaterials {
		if _, has := recipeMap[recipeMaterial.RecipeGalleryId]; has {
			newRecipeMap[recipeMaterial.RecipeGalleryId] = recipeMap[recipeMaterial.RecipeGalleryId]
		}
	}
	for k, _ := range newRecipeMap {
		result = append(result, newRecipeMap[k])
	}
	return result, ""
}

// 根据食材列表筛选菜谱
func filterRecipesByMaterials(recipes []database.Recipe, materials []string) ([]database.Recipe, string) {
	if len(materials) == 0 {
		return nil, "你想查什么食材呀"
	}
	if len(recipes) == 0 {
		return recipes, ""
	}
	result := recipes
	var note string
	materialCount := 0
	for _, material := range materials {
		if material == "" {
			continue
		} else {
			result, note = filterRecipesByMaterial(result, material)
			if note != "" {
				return nil, note
			}
			materialCount++
		}
	}
	if materialCount == 0 {
		return nil, "你想查什么食材呀"
	}

	return result, ""
}

// 根据技法筛选菜谱
func filterRecipesBySkill(recipes []database.Recipe, skill string) ([]database.Recipe, string) {
	if len(recipes) == 0 {
		return recipes, ""
	}
	result := make([]database.Recipe, 0)
	for _, recipe := range recipes {
		switch skill {
		case "炒":
			if recipe.Stirfry > 0 {
				result = append(result, recipe)
			}
		case "烤":
			if recipe.Bake > 0 {
				result = append(result, recipe)
			}
		case "煮":
			if recipe.Boil > 0 {
				result = append(result, recipe)
			}
		case "蒸":
			if recipe.Steam > 0 {
				result = append(result, recipe)
			}
		case "炸":
			if recipe.Fry > 0 {
				result = append(result, recipe)
			}
		case "切":
			if recipe.Cut > 0 {
				result = append(result, recipe)
			}
		default:
			return nil, fmt.Sprintf("%s是什么技法呀", skill)
		}
	}
	return result, ""
}

// 根据技法列表筛选菜谱
func filterRecipesBySkills(recipes []database.Recipe, skills []string) ([]database.Recipe, string) {
	if len(skills) == 0 {
		return nil, "你想查什么技法呀"
	}
	if len(recipes) == 0 {
		return recipes, ""
	}
	result := recipes
	var note string
	skillCount := 0
	for _, skill := range skills {
		if skill == "" {
			continue
		} else {
			result, note = filterRecipesBySkill(result, skill)
			if note != "" {
				return nil, note
			}
			skillCount++
		}
	}
	if skillCount == 0 {
		return nil, "你想查什么技法呀"
	}
	return result, ""
}

// 根据贵客筛选菜谱
func filterRecipeByGuest(recipes []database.Recipe, guest string) ([]database.Recipe, string) {
	if len(recipes) == 0 {
		return recipes, ""
	}
	result := make([]database.Recipe, 0)
	// 将所有recipe存入map
	recipeMap := make(map[string]database.Recipe)
	for _, recipe := range recipes {
		recipeMap[recipe.Name] = recipe
	}
	// 根据贵客名找出对应的菜谱
	guestGifts := make([]database.GuestGift, 0)
	err := database.DB.Where("guest_name like ?", "%"+guest+"%").Find(&guestGifts)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		return nil, util.SystemErrorNote
	}
	if len(guestGifts) == 0 {
		return nil, fmt.Sprintf("%s是什么神秘符文呀", guest)
	}
	// 将符合条件的菜谱存入新map
	newRecipeMap := make(map[string]database.Recipe)
	for _, guestGift := range guestGifts {
		if _, has := recipeMap[guestGift.Recipe]; has {
			newRecipeMap[guestGift.Recipe] = recipeMap[guestGift.Recipe]
		}
	}
	for k, _ := range newRecipeMap {
		result = append(result, newRecipeMap[k])
	}
	return result, ""
}

// 根据贵客列表查询菜谱
func filterRecipesByGuests(recipes []database.Recipe, guests []string) ([]database.Recipe, string) {
	if len(guests) == 0 {
		return nil, "你想查询哪位贵客呀"
	}
	result := recipes
	var note string
	guestCount := 0
	for _, guest := range guests {
		if guest == "" {
			continue
		} else {
			result, note = filterRecipeByGuest(result, guest)
			if note != "" {
				return nil, note
			}
			guestCount++
		}
	}
	if guestCount == 0 {
		return nil, "你想查询哪位贵客呀"
	}

	return result, ""
}

// 根据符文礼物查询菜谱
func filterRecipesByAntique(recipes []database.Recipe, antique string) ([]database.Recipe, string) {
	if len(recipes) == 0 {
		return recipes, ""
	}
	result := make([]database.Recipe, 0)
	// 将所有recipe存入map
	recipeMap := make(map[string]database.Recipe)
	for _, recipe := range recipes {
		recipeMap[recipe.Name] = recipe
	}
	// 根据符文礼物名找出对应的菜谱
	guestGifts := make([]database.GuestGift, 0)
	err := database.DB.Where("antique like ?", "%"+antique+"%").Find(&guestGifts)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		return nil, util.SystemErrorNote
	}
	if len(guestGifts) == 0 {
		return nil, fmt.Sprintf("%s是什么神秘符文呀", antique)
	}
	// 将符合条件的recipe存入新map
	newRecipeMap := make(map[string]database.Recipe)
	for _, guestGift := range guestGifts {
		if _, has := recipeMap[guestGift.Recipe]; has {
			newRecipeMap[guestGift.Recipe] = recipeMap[guestGift.Recipe]
		}
	}
	for k, _ := range newRecipeMap {
		result = append(result, newRecipeMap[k])
	}
	return result, ""
}

// 根据来源筛选菜谱
func filterRecipesByOrigin(recipes []database.Recipe, origin string) ([]database.Recipe, string) {
	if len(recipes) == 0 {
		return recipes, ""
	}
	result := make([]database.Recipe, 0)
	pattern := ".*" + strings.ReplaceAll(origin, "%", ".*") + ".*"
	for i, _ := range recipes {
		re := regexp.MustCompile(pattern)
		if re.MatchString(recipes[i].Origin) {
			result = append(result, recipes[i])
		}
	}
	return result, ""
}

// 根据调料筛选菜谱
func filterRecipesByCondiment(recipes []database.Recipe, condiment string) ([]database.Recipe, string) {
	if len(recipes) == 0 {
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
		return nil, fmt.Sprintf("%s是啥味道呀", condiment)
	}
	for i, _ := range recipes {
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
		pattern := ".*" + strings.ReplaceAll(name, "%", ".*") + ".*"
		for i, _ := range recipes {
			re := regexp.MustCompile(pattern)
			if re.MatchString(recipes[i].Name) {
				result = append(result, recipes[i])
			}
		}
	} else {
		galleryId := fmt.Sprintf("%03d", numId)
		for i, _ := range recipes {
			if recipes[i].GalleryId == galleryId {
				result = append(result, recipes[i])
			}
		}

	}
	return result, ""
}

// 根据菜谱单价筛选菜谱
func filterRecipesByPrice(recipes []database.Recipe, price int) ([]database.Recipe, string) {
	result := make([]database.Recipe, 0)
	for i, _ := range recipes {
		if recipes[i].Price >= price {
			result = append(result, recipes[i])
		}
	}
	return result, ""
}

type recipeWrapper struct {
	recipe     []database.Recipe
	recipeLess func(p *database.Recipe, q *database.Recipe) bool
}

func (w recipeWrapper) Len() int {
	return len(w.recipe)
}

func (w recipeWrapper) Swap(i int, j int) {
	w.recipe[i], w.recipe[j] = w.recipe[j], w.recipe[i]
}

func (w recipeWrapper) Less(i int, j int) bool {
	return w.recipeLess(&w.recipe[i], &w.recipe[j])
}

// 根据排序参数排序菜谱
func orderRecipes(recipes []database.Recipe, order string) ([]database.Recipe, string) {
	if len(recipes) == 0 {
		return recipes, ""
	}
	switch order {
	case "图鉴序":
		sort.Sort(recipeWrapper{recipes, func(m, n *database.Recipe) bool {
			return m.RecipeId < n.RecipeId
		}})
	case "单时间":
		sort.Sort(recipeWrapper{recipes, func(m, n *database.Recipe) bool {
			return m.Time < n.Time
		}})
	case "总时间":
		sort.Sort(recipeWrapper{recipes, func(m, n *database.Recipe) bool {
			return m.TotalTime < n.TotalTime
		}})
	case "单价", "售价":
		sort.Sort(recipeWrapper{recipes, func(m, n *database.Recipe) bool {
			return m.Price > n.Price
		}})
	case "金币效率":
		sort.Sort(recipeWrapper{recipes, func(m, n *database.Recipe) bool {
			return m.GoldEfficiency > n.GoldEfficiency
		}})
	case "耗材效率":
		sort.Sort(recipeWrapper{recipes, func(m, n *database.Recipe) bool {
			return m.MaterialEfficiency > n.MaterialEfficiency
		}})
	default:
		return nil, "排序参数有误"
	}
	return recipes, ""
}

// 根据排序规则与分页参数，返回菜谱列表消息数据
func echoRecipesMessage(recipes []database.Recipe, order string, page int, private bool) string {
	if len(recipes) == 0 {
		logger.Debug("未查询到菜谱")
		return "本店没有这道菜呢!"
	} else if len(recipes) == 1 {
		logger.Debug("查询到一个菜谱")
		return getRecipeMessage(recipes[0])
	} else {
		logger.Debug("查询到多个菜谱")
		var msg string
		listLength := util.MaxQueryListLength
		if private {
			listLength = listLength * 2
		}
		maxPage := (len(recipes)-1)/listLength + 1
		if page > maxPage {
			page = maxPage
		}
		if len(recipes) > listLength {
			msg += fmt.Sprintf("这里有你想点的菜吗: (%d/%d)\n", page, maxPage)
		} else {
			msg += "这里有你想点的菜吗:\n"
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
	default:
		return ""
	}
}

// 输出单菜谱消息数据
func getRecipeMessage(recipe database.Recipe) string {
	// 尝试寻找图片文件，未找到则按照文字格式发送
	resourceImageDir := config.AppConfig.Resource.Image + "/recipe"
	imagePath := fmt.Sprintf("%s/recipe_%s.png", resourceImageDir, recipe.GalleryId)
	logger.Debug("imagePath:", imagePath)
	var msg string
	if has, err := util.PathExists(imagePath); has {
		msg = bot.GetCQImage(imagePath, "file")
	} else {
		if err != nil {
			logger.Debugf("无法确定文件是否存在!", err)
		}
		logger.Info("未找到菜谱图鉴图片, 以文字格式发送数据")
		// 稀有度数据
		rarity := ""
		for i := 0; i < recipe.Rarity; i++ {
			rarity += "🔥"
		}
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
		materials := ""
		recipeMaterials := make([]database.RecipeMaterial, 0)
		err := database.DB.Where("recipe_id = ?", recipe.GalleryId).Find(&recipeMaterials)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return util.SystemErrorNote
		}
		for _, recipeMaterial := range recipeMaterials {
			material := new(database.Material)
			has, err := database.DB.Where("material_id = ?", recipeMaterial.MaterialId).Get(material)
			if err != nil {
				logger.Error("查询数据库出错!", err)
				return util.SystemErrorNote
			}
			if !has {
				logger.Warnf("菜谱%d数据缺失", recipeMaterial.MaterialId)
			} else {
				materials += fmt.Sprintf("%s*%d ", material.Name, recipeMaterial.Quantity)
			}
		}
		// 贵客礼物数据
		giftInfo := ""
		guestGifts := make([]database.GuestGift, 0)
		err = database.DB.Where("recipe = ?", recipe.Name).Find(&guestGifts)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return util.SystemErrorNote
		}
		for _, gift := range guestGifts {
			if giftInfo != "" {
				giftInfo += ", "
			}
			giftInfo += fmt.Sprintf("%s-%s", gift.GuestName, gift.Antique)
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
		msg += fmt.Sprintf("%s %s %s\n", recipe.GalleryId, recipe.Name, rarity)
		msg += fmt.Sprintf("💰: %d(%d) --- %d/h\n", recipe.Price, recipe.Price+recipe.ExPrice, recipe.GoldEfficiency)
		msg += fmt.Sprintf("来源: %s\n", recipe.Origin)
		msg += fmt.Sprintf("单时间: %s\n", util.FormatSecondToString(recipe.Time))
		msg += fmt.Sprintf("总时间: %s (%d份)\n", util.FormatSecondToString(recipe.Time*recipe.Limit), recipe.Limit)
		msg += fmt.Sprintf("技法: %s\n", recipeSkill)
		msg += fmt.Sprintf("食材: %s\n", materials)
		msg += fmt.Sprintf("耗材效率: %d/h\n", recipe.MaterialEfficiency)
		msg += fmt.Sprintf("可解锁: %s\n", recipe.Unlock)
		msg += fmt.Sprintf("可合成: %s\n", recipe.Combo)
		msg += fmt.Sprintf("神级符文: %s\n", recipe.Gift)
		msg += fmt.Sprintf("贵客礼物: %s\n", giftInfo)
		msg += fmt.Sprintf("升阶贵客: %s", guests)
	}
	return msg
}

func RecipeInfoToImage(recipes []database.Recipe) error {
	dx := 800       // 图鉴背景图片的宽度
	dy := 800       // 图鉴背景图片的高度
	titleSize := 48 // 标题字体尺寸
	fontSize := 32  // 内容字体尺寸
	fontDPI := 72.0 // dpi

	resourceFontDir := config.AppConfig.Resource.Font
	fontPath := "yuan500W.ttf"
	fontFile := fmt.Sprintf("%s/%s", resourceFontDir, fontPath) // 需要使用的字体文件
	resourceImgDir := config.AppConfig.Resource.Image
	recipeImgPath := resourceImgDir + "/recipe"

	for _, recipe := range recipes {
		bgFile, err := os.Open(fmt.Sprintf("%s/recipe_%s.png", recipeImgPath, recipe.Condiment))
		if err != nil {
			return err
		}
		defer bgFile.Close()
		img := image.NewRGBA(image.Rect(0, 0, dx, dy))
		bg, _ := png.Decode(bgFile)

		draw.Draw(img, img.Bounds(), bg, bg.Bounds().Min, draw.Src)

		//读字体数据
		fontBytes, err := ioutil.ReadFile(fontFile)
		if err != nil {
			return err
		}
		font, err := freetype.ParseFont(fontBytes)
		if err != nil {
			return err
		}

		c := freetype.NewContext()
		c.SetDPI(fontDPI)
		c.SetFont(font)
		c.SetClip(img.Bounds())
		c.SetDst(img)
		fontColor := color.RGBA{0, 0, 0, 255}
		c.SetSrc(image.NewUniform(fontColor))

		// 输出图鉴ID与菜谱名
		c.SetFontSize(float64(titleSize))
		pt := freetype.Pt(20, 20+titleSize)
		_, err = c.DrawString(fmt.Sprintf("%s %s", recipe.GalleryId, recipe.Name), pt)
		if err != nil {
			return err
		}
		// 输出稀有度
		coverRect := image.Rect(540+recipe.Rarity*48, 28, 780, 72)
		bgColor := color.RGBA{255, 242, 226, 255}
		draw.Draw(img, coverRect, image.NewUniform(bgColor), image.ZP, draw.Src)

		// 输出单价信息
		fontColor = color.RGBA{45, 45, 45, 255}
		c.SetSrc(image.NewUniform(fontColor))
		c.SetFontSize(float64(fontSize))
		pt = freetype.Pt(94, 106+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%d", recipe.Price), pt)
		if err != nil {
			return err
		}
		fontColor = color.RGBA{120, 120, 120, 255}
		c.SetSrc(image.NewUniform(fontColor))
		pt = freetype.Pt(174, 106+fontSize)
		_, err = c.DrawString(fmt.Sprintf("+%d", recipe.ExPrice), pt)
		if err != nil {
			return err
		}
		fontColor = color.RGBA{45, 45, 45, 255}
		c.SetSrc(image.NewUniform(fontColor))
		// 输出金币效率
		pt = freetype.Pt(358, 106+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%d / h", recipe.GoldEfficiency), pt)
		if err != nil {
			return err
		}
		// 输出份数
		pt = freetype.Pt(584, 106+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%d", recipe.Limit), pt)
		if err != nil {
			return err
		}
		// 输出单时间
		pt = freetype.Pt(150, 184+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%s", util.FormatSecondToString(recipe.Time)), pt)
		if err != nil {
			return err
		}
		// 输出总时间
		pt = freetype.Pt(500, 184+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%s", util.FormatSecondToString(recipe.TotalTime)), pt)
		if err != nil {
			return err
		}
		// 输出技法
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
		pt = freetype.Pt(110, 262+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%s", recipeSkill), pt)
		if err != nil {
			return err
		}
		// 输出耗材效率
		pt = freetype.Pt(530, 262+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%d / h", recipe.MaterialEfficiency), pt)
		if err != nil {
			return err
		}
		// 输出食材
		materials := ""
		recipeMaterials := make([]database.RecipeMaterial, 0)
		err = database.DB.Where("recipe_id = ?", recipe.GalleryId).Find(&recipeMaterials)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return err
		}
		for _, recipeMaterial := range recipeMaterials {
			material := new(database.Material)
			has, err := database.DB.Where("material_id = ?", recipeMaterial.MaterialId).Get(material)
			if err != nil {
				logger.Error("查询数据库出错!", err)
				return err
			}
			if !has {
				logger.Warnf("菜谱%d数据缺失", recipeMaterial.MaterialId)
			} else {
				materials += fmt.Sprintf("%s*%d ", material.Name, recipeMaterial.Quantity)
			}
		}
		pt = freetype.Pt(110, 342+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%s", materials), pt)
		if err != nil {
			return err
		}
		// 输出贵客礼物
		giftInfo := ""
		guestGifts := make([]database.GuestGift, 0)
		err = database.DB.Where("recipe = ?", recipe.Name).Find(&guestGifts)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return err
		}
		for _, gift := range guestGifts {
			if giftInfo != "" {
				giftInfo += ", "
			}
			giftInfo += fmt.Sprintf("%s-%s", gift.GuestName, gift.Antique)
		}
		pt = freetype.Pt(110, 420+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%s", giftInfo), pt)
		if err != nil {
			return err
		}
		// 输出升阶贵客
		for p, guest := range recipe.Guests {
			if guest == "" {
				guest = "未知"
			}
			pt = freetype.Pt(84, 556+p*78+fontSize)
			_, err = c.DrawString(fmt.Sprintf("%s", guest), pt)
			if err != nil {
				return err
			}
		}
		// 输出来源
		pt = freetype.Pt(460, 500+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%s", recipe.Origin), pt)
		if err != nil {
			return err
		}
		// 输出神级符文
		pt = freetype.Pt(520, 580+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%s", recipe.Gift), pt)
		if err != nil {
			return err
		}
		// 输出可解锁
		pt = freetype.Pt(490, 658+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%s", recipe.Unlock), pt)
		if err != nil {
			return err
		}
		// 输出可合成
		pt = freetype.Pt(490, 734+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%s", recipe.Combo), pt)
		if err != nil {
			return err
		}

		// 以PNG格式保存文件
		dst, err := os.Create(fmt.Sprintf("%s/recipe_%s.png", recipeImgPath, recipe.GalleryId))
		if err != nil {
			return err
		}
		err = png.Encode(dst, img)
		if err != nil {
			return err
		}
		dst.Close()
	}
	return nil
}
