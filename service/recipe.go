package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/model/database"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"fmt"
	"strconv"
	"strings"
)

// 处理菜谱查询请求
func RecipeQuery(c *onebot.Context, args []string) {
	logger.Info("菜谱查询, 参数:", args)
	if len(args) == 0 {
		err := bot.SendMessage(c, recipeHelp())
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
		return
	}

	recipes := make([]database.Recipe, 0)
	note := ""
	order := ""
	rarity := 1
	price := 1
	page := 1
	if len(args) == 1 {
		// 处理简单查询
		recipes, note = getRecipesWithName(args[0])
	} else {
		// 处理组合查询
		for i := 1; i < len(args); i++ {
			updateQueryArgs(args[i], &order, &rarity, &price, &page)
		}
		switch args[0] {
		case "任意", "%":
			recipes, note = getAllRecipes(order)
		case "食材", "材料":
			recipes, note = getRecipesWithMaterial(args[1], order)
		case "技法":
			recipes, note = getRecipesWithSkill(args[1], order)
		case "贵客":
			recipes, note = getRecipesWithGuest(args[1], order)
		case "符文", "礼物":
			recipes, note = getRecipesWithAntique(args[1], order)
		case "来源":
			recipes, note = getRecipesWithOrigin(args[1], order)
		default:
			note = util.QueryParamWrongNote
		}
	}
	if note != "" {
		logger.Info("菜谱查询失败结果:", note)
		_ = bot.SendMessage(c, note)
		return
	}

	msg := getRecipesMessage(recipes, order, rarity, price, page)
	logger.Info("发送菜谱查询结果:", msg)
	err := bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}

// 更新查询参数信息
func updateQueryArgs(arg string, order *string, rarity *int, price *int, page *int) {
	switch arg {
	// 判断是否是排序参数
	case "图鉴序", "单时间", "总时间", "单价", "金币效率", "耗材效率", "食材效率":
		*order = arg
	// 判断是否是稀有度筛选参数
	case "1火", "1星", "一火", "一星":
		*rarity = 1
	case "2火", "2星", "二火", "二星", "两火", "两星":
		*rarity = 2
	case "3火", "3星", "三火", "三星":
		*rarity = 3
	case "4火", "4星", "四火", "四星":
		*rarity = 4
	case "5火", "5星", "五火", "五星":
		*rarity = 5
	default:
		// 判断是否是单价筛选参数
		if strings.HasPrefix(arg, "$") {
			num, err := strconv.Atoi(arg[1:])
			if err != nil {
				return
			} else {
				*price = num
				return
			}
		}
		// 判断是否是分页参数
		if strings.HasPrefix(arg, "p") || strings.HasPrefix(arg, "P") {
			num, err := strconv.Atoi(arg[1:])
			if err != nil {
				return
			} else {
				if num < 1 {
					num = 1
				}
				*page = num
				return
			}
		}
	}
}

// 根据排序参数获取order by的sql语句
func getRecipeOrderString(order string) (string, bool) {
	switch order {
	case "单时间":
		return "`time` ASC", true
	case "总时间":
		return "`total_time` ASC", true
	case "单价":
		return "`price` DESC", true
	case "金币效率":
		return "`gold_efficiency` DESC", true
	case "耗材效率":
		return "`material_efficiency` DESC", true
	case "":
		return "`gallery_id` ASC", true
	default:
		return "", false
	}
}

// 根据排序参数获取菜谱需要输出的信息
func getRecipeInfoWithOrder(recipe database.Recipe, order string) string {
	switch order {
	case "单时间":
		return util.FormatSecondToString(recipe.Time)
	case "总时间":
		return util.FormatSecondToString(recipe.Time * recipe.Limit)
	case "单价":
		return fmt.Sprintf("💰%d", recipe.Price)
	case "金币效率":
		return fmt.Sprintf("💰%d/h", recipe.GoldEfficiency)
	case "耗材效率":
		return fmt.Sprintf("🥗%d/h", recipe.MaterialEfficiency)
	case "食材效率":
		return fmt.Sprintf("🥗%d/h", recipe.MaterialEfficiency)
	case "":
		return ""
	default:
		return ""
	}
}

// 输出单菜谱消息数据
func getRecipeMessage(recipe database.Recipe) string {
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
	// 组合消息信息
	var msg string
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
	return msg
}

// 根据排序规则、稀有度、售价与分页参数，返回消息数据
func getRecipesMessage(recipes []database.Recipe, order string, rarity int, price int, page int) string {
	if len(recipes) == 0 {
		logger.Debug("未查询到菜谱")
		return "本店没有这道菜呢!"
	} else if len(recipes) == 1 {
		logger.Debug("查询到一个菜谱")
		return getRecipeMessage(recipes[0])
	} else {
		logger.Debug("查询到多个菜谱")
		results := make([]database.Recipe, 0)
		for _, recipe := range recipes {
			if recipe.Rarity >= rarity && recipe.Price >= price {
				results = append(results, recipe)
			}
		}
		var msg string
		listLength := util.MaxQueryListLength
		maxPage := (len(results)-1)/listLength + 1
		if len(results) > listLength {
			if page > maxPage {
				page = maxPage
			}
			msg += fmt.Sprintf("这里有你想点的菜吗: (%d/%d)\n", page, maxPage)
		} else {
			msg += "这里有你想点的菜吗:\n"
		}
		for i := (page - 1) * listLength; i < page*listLength && i < len(results); i++ {
			orderInfo := getRecipeInfoWithOrder(results[i], order)
			msg += fmt.Sprintf("%s %s %s", results[i].GalleryId, results[i].Name, orderInfo)
			if i < page*listLength-1 && i < len(results)-1 {
				msg += "\n"
			}
		}
		if page < maxPage {
			msg += "\n......"
		}

		return msg
	}
}

// 根据菜谱名字或ID查询菜谱
func getRecipesWithName(arg string) ([]database.Recipe, string) {
	recipes := make([]database.Recipe, 0)
	err := database.DB.Where("gallery_id = ?", arg).Asc("gallery_id").Find(&recipes)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		return nil, util.SystemErrorNote
	}
	if len(recipes) == 0 {
		err = database.DB.Where("name like ?", "%"+arg+"%").Asc("gallery_id").Find(&recipes)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return nil, util.SystemErrorNote
		}
	}
	return recipes, ""
}

// 参数"任意", 查询出所有菜谱
func getAllRecipes(order string) ([]database.Recipe, string) {
	recipes := make([]database.Recipe, 0)
	orderStr, success := getRecipeOrderString(order)
	if !(success) {
		return nil, util.QueryParamWrongNote
	}
	err := database.DB.OrderBy(orderStr).Find(&recipes)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		return nil, util.SystemErrorNote
	}
	return recipes, ""
}

// 根据食材名字查询菜谱
func getRecipesWithMaterial(arg string, order string) ([]database.Recipe, string) {
	// 根据食材名查询食材信息
	materials := make([]database.Material, 0)
	materialsId := make([]int, 0)
	var materialOrigin []string
	switch arg {
	case "鱼类", "水产", "海鲜":
		materialOrigin = []string{"池塘"}
	case "蔬菜", "菜类":
		materialOrigin = []string{"菜棚", "菜地", "森林"}
	case "肉类":
		materialOrigin = []string{"牧场", "鸡舍", "猪圈"}
	case "面类":
		materialOrigin = []string{"作坊"}
	default:
		materialOrigin = []string{}
	}
	if len(materialOrigin) != 0 {
		err := database.DB.In("origin", materialOrigin).Find(&materials)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return nil, util.SystemErrorNote
		}
	} else {
		err := database.DB.Where("name = ?", arg).Find(&materials)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return nil, util.SystemErrorNote
		}
		if len(materials) == 0 {
			return nil, fmt.Sprintf("厨师长说没有用%s做过菜", arg)
		}
	}
	for _, material := range materials {
		materialsId = append(materialsId, material.MaterialId)
	}
	recipes := make([]database.Recipe, 0)
	recipeMaterials := make([]database.RecipeMaterial, 0)
	if order == "食材效率" {
		// 根据食材id查菜谱-食材表并根据食材效率排序
		err := database.DB.In("material_id", materialsId).Desc("efficiency").Find(&recipeMaterials)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return nil, util.SystemErrorNote
		}
		// 根据查出的信息查询菜谱信息
		for _, recipeMaterial := range recipeMaterials {
			var recipe database.Recipe
			has, err := database.DB.Where("gallery_id = ?", recipeMaterial.RecipeGalleryId).Get(&recipe)
			if err != nil {
				logger.Error("查询数据库出错!", err)
				return nil, util.SystemErrorNote
			}
			if !has {
				logger.Warnf("菜谱%s的食材信息可能有误!", recipeMaterial.RecipeGalleryId)
				continue
			}
			recipe.MaterialEfficiency = recipeMaterial.Efficiency
			recipes = append(recipes, recipe)
		}
	} else {
		// 根据食材id查菜谱-食材表
		err := database.DB.In("material_id", materialsId).Find(&recipeMaterials)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return nil, util.SystemErrorNote
		}
		// 根据菜谱id查询菜谱信息并根据order参数排序
		recipeIds := make([]string, 0)
		for _, recipeMaterial := range recipeMaterials {
			recipeIds = append(recipeIds, recipeMaterial.RecipeGalleryId)
		}
		orderStr, success := getRecipeOrderString(order)
		if !(success) {
			return nil, util.QueryParamWrongNote
		}
		err = database.DB.In("gallery_id", recipeIds).OrderBy(orderStr).Find(&recipes)

		if err != nil {
			logger.Error("查询数据库出错!", err)
			return nil, util.SystemErrorNote
		}
	}
	return recipes, ""
}

func getRecipesWithSkill(arg string, order string) ([]database.Recipe, string) {
	var skill string
	switch arg {
	case "炒":
		skill = "`stirfry` > 0"
	case "烤":
		skill = "`bake` > 0"
	case "煮":
		skill = "`boil` > 0"
	case "蒸":
		skill = "`steam` > 0"
	case "炸":
		skill = "`fry` > 0"
	case "切":
		skill = "`cut` > 0"
	default:
		return nil, util.QueryParamWrongNote
	}

	orderStr, success := getRecipeOrderString(order)
	if !(success) {
		return nil, util.QueryParamWrongNote
	}

	recipes := make([]database.Recipe, 0)
	err := database.DB.Where(skill).OrderBy(orderStr).Find(&recipes)
	if err != nil {
		logger.Error("数据库查询出错!", err)
		return nil, util.SystemErrorNote
	}
	return recipes, ""
}

func getRecipesWithGuest(arg string, order string) ([]database.Recipe, string) {
	guests := make([]database.GuestGift, 0)
	err := database.DB.Where("guest_id = ?", arg).Find(&guests)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		return nil, util.SystemErrorNote
	}

	if len(guests) == 0 {
		err = database.DB.Where("guest_name like ?", "%"+arg+"%").Find(&guests)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return nil, util.SystemErrorNote
		}
	}

	if len(guests) == 0 {
		return nil, "没有找到该贵客"
	}

	recipesName := make([]string, 0)

	for _, guest := range guests {
		recipesName = append(recipesName, guest.Recipe)
	}

	orderStr, success := getRecipeOrderString(order)
	if !(success) {
		return nil, util.QueryParamWrongNote
	}

	recipes := make([]database.Recipe, 0)
	err = database.DB.In("name", recipesName).OrderBy(orderStr).Find(&recipes)
	if err != nil {
		logger.Error("数据库查询出错!", err)
		return nil, util.SystemErrorNote
	}
	return recipes, ""
}

func getRecipesWithAntique(arg string, order string) ([]database.Recipe, string) {
	guests := make([]database.GuestGift, 0)
	err := database.DB.Where("antique like ?", "%"+arg+"%").Find(&guests)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		return nil, util.SystemErrorNote
	}

	if len(guests) == 0 {
		return nil, "没有找到该符文"
	}

	recipesName := make([]string, 0)

	for _, guest := range guests {
		recipesName = append(recipesName, guest.Recipe)
	}

	orderStr, success := getRecipeOrderString(order)
	if !(success) {
		return nil, util.QueryParamWrongNote
	}

	recipes := make([]database.Recipe, 0)
	err = database.DB.In("name", recipesName).OrderBy(orderStr).Find(&recipes)
	if err != nil {
		logger.Error("数据库查询出错!", err)
		return nil, util.SystemErrorNote
	}
	return recipes, ""
}

func getRecipesWithOrigin(arg string, order string) ([]database.Recipe, string) {
	orderStr, success := getRecipeOrderString(order)
	if !(success) {
		return nil, util.QueryParamWrongNote
	}

	recipes := make([]database.Recipe, 0)
	err := database.DB.Where("origin like ?", "%"+arg+"%").OrderBy(orderStr).Find(&recipes)
	if err != nil {
		logger.Error("数据库查询出错!", err)
		return nil, util.SystemErrorNote
	}
	return recipes, ""
}
