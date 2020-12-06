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
	queryType := ""
	arg := ""
	getArg := false
	order := ""
	rarity := 1
	price := 1
	page := 1
	condition := 0
	// 处理组合查询
	for i := 0; i < len(args); i++ {
		c := updateQueryArgs(args[i], &queryType, &order, &rarity, &price, &page)
		if queryType != "" && i+1 < len(args) && !getArg {
			arg = args[i+1]
			getArg = true
		}
		if c != 0 {
			condition = 1
		}
	}
	if queryType != "" && queryType != "任意" && arg == "" {
		_ = bot.SendMessage(c, "请填一下查询参数哦")
		return
	}
	if queryType == "" && condition == 1 {
		queryType = "任意"
	}
	switch queryType {
	case "任意", "%":
		recipes, note = getAllRecipes(order)
	case "食材", "材料":
		recipes, note = getRecipesWithMaterial(arg, order)
	case "技法":
		recipes, note = getRecipesWithSkill(arg, order)
	case "贵客":
		recipes, note = getRecipesWithGuest(arg, order)
	case "符文", "礼物":
		recipes, note = getRecipesWithAntique(arg, order)
	case "来源":
		recipes, note = getRecipesWithOrigin(arg, order)
	default:
		if len(args) == 1 && condition == 0 {
			// 处理简单查询
			recipes, note = getRecipesWithName(args[0])
		} else {
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

// 更新查询参数信息, 返回值1表示有修改, 0表示无修改
func updateQueryArgs(arg string, queryType *string, order *string, rarity *int, price *int, page *int) int {
	switch arg {
	// 判断是否是查询类型参数
	case "任意", "食材", "材料", "技法", "贵客", "符文", "礼物", "来源":
		*queryType = arg
		return 1
	// 判断是否是排序参数
	case "图鉴序", "单时间", "总时间", "单价", "金币效率", "耗材效率", "食材效率":
		*order = arg
		return 1
	// 判断是否是稀有度筛选参数
	case "1火", "1星", "一火", "一星":
		*rarity = 1
		return 1
	case "2火", "2星", "二火", "二星", "两火", "两星":
		*rarity = 2
		return 1
	case "3火", "3星", "三火", "三星":
		*rarity = 3
		return 1
	case "4火", "4星", "四火", "四星":
		*rarity = 4
		return 1
	case "5火", "5星", "五火", "五星":
		*rarity = 5
		return 1
	default:
		// 判断是否是单价筛选参数
		if strings.HasPrefix(arg, "$") {
			num, err := strconv.Atoi(arg[1:])
			if err != nil {
				return 0
			} else {
				*price = num
				return 1
			}
		}
		// 判断是否是分页参数
		if strings.HasPrefix(arg, "p") || strings.HasPrefix(arg, "P") {
			num, err := strconv.Atoi(arg[1:])
			if err != nil {
				return 0
			} else {
				if num < 1 {
					num = 1
				}
				*page = num
				return 1
			}
		}
	}
	return 0
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

	gallery := recipeGallery{
		GalleryId:          recipe.GalleryId,
		Name:               recipe.Name,
		Rarity:             recipe.Rarity,
		Price:              recipe.Price,
		ExPrice:            recipe.ExPrice,
		Limit:              recipe.Limit,
		GoldEfficiency:     recipe.GoldEfficiency,
		Origin:             recipe.Origin,
		SingleTime:         util.FormatSecondToString(recipe.Time),
		TotalTime:          util.FormatSecondToString(recipe.Time * recipe.Limit),
		Skill:              recipeSkill,
		Condiment:          recipe.Condiment,
		Materials:          materials,
		MaterialEfficiency: recipe.MaterialEfficiency,
		Unlock:             recipe.Unlock,
		Combo:              recipe.Combo,
		Gift:               recipe.Gift,
		GuestGift:          giftInfo,
		UpgradeGuests:      recipe.Guests,
	}

	resourceImageDir := config.AppConfig.Resource.Image + "recipe"
	imagePath := fmt.Sprintf("%s/recipe_%s.png", resourceImageDir, recipe.GalleryId)
	logger.Debug("imagePath:", imagePath)
	if has, err := util.PathExists(imagePath); !has {
		if err != nil {
			logger.Debugf("无法确定文件是否存在!", err)
		}
		logger.Info("未找到菜谱图鉴，重新生成")
		dst, _ := os.Create(imagePath)
		defer dst.Close()
		err = RecipeInfoToImage(gallery, dst)
		if err != nil {
			logger.Error("菜谱数据转图鉴出错!", err)
			return util.SystemErrorNote
		}
	}
	msg := bot.GetCQImage(imagePath, "file")

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

type recipeGallery struct {
	GalleryId          string
	Name               string
	Rarity             int
	Price              int
	ExPrice            int
	GoldEfficiency     int
	Limit              int
	Origin             string
	SingleTime         string
	TotalTime          string
	Condiment          string
	Skill              string
	Materials          string
	MaterialEfficiency int
	Unlock             string
	Combo              string
	Gift               string
	GuestGift          string
	UpgradeGuests      []string
}

func RecipeInfoToImage(recipe recipeGallery, dst *os.File) error {
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
	_, err = c.DrawString(fmt.Sprintf("%s", recipe.SingleTime), pt)
	if err != nil {
		return err
	}
	// 输出总时间
	pt = freetype.Pt(500, 184+fontSize)
	_, err = c.DrawString(fmt.Sprintf("%s", recipe.TotalTime), pt)
	if err != nil {
		return err
	}
	// 输出技法
	pt = freetype.Pt(110, 262+fontSize)
	_, err = c.DrawString(fmt.Sprintf("%s", recipe.Skill), pt)
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
	pt = freetype.Pt(110, 342+fontSize)
	_, err = c.DrawString(fmt.Sprintf("%s", recipe.Materials), pt)
	if err != nil {
		return err
	}
	// 输出贵客礼物?/
	pt = freetype.Pt(110, 420+fontSize)
	_, err = c.DrawString(fmt.Sprintf("%s", recipe.GuestGift), pt)
	if err != nil {
		return err
	}
	// 输出升阶贵客
	for p, guest := range recipe.UpgradeGuests {
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
	err = png.Encode(dst, img)
	if err != nil {
		return err
	}
	return nil
}
