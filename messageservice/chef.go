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

func ChefQuery(c *scheduler.Context) {
	order := "稀有度"
	page := 1
	var note string
	chefs, err := dao.FindAllChefs()
	if err != nil {
		logger.Error("查询厨师数据出错!", err)
		_, _ = c.Reply(e.SystemErrorNote)
	}
	args := strings.Split(c.PretreatedMessage, " ")
	for _, arg := range args {
		switch arg {
		case "图鉴序", "稀有度":
			order = arg
		case "1火", "1星", "一火", "一星":
			chefs, note = filterChefsByRarity(chefs, 1)
		case "2火", "2星", "二火", "二星", "两火", "两星":
			chefs, note = filterChefsByRarity(chefs, 2)
		case "3火", "3星", "三火", "三星":
			chefs, note = filterChefsByRarity(chefs, 3)
		case "4火", "4星", "四火", "四星":
			chefs, note = filterChefsByRarity(chefs, 4)
		case "5火", "5星", "五火", "五星":
			chefs, note = filterChefsByRarity(chefs, 5)
		default:
			if util.HasPrefixIn(arg, "来源") {
				origin := strings.Split(arg, "-")
				if len(origin) > 1 {
					chefs, note = filterChefsByOrigin(chefs, origin[1])
				}
			} else if util.HasPrefixIn(arg, "技能") {
				skill := strings.Split(arg, "-")
				if len(skill) > 1 {
					chefs, note = filterChefsBySkill(chefs, strings.Join(skill[1:], "-"))
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
				chefs, note = filterChefsByName(chefs, arg)
			}
		}
		if note != "" {
			logger.Info("厨师查询失败:", note)
			_, _ = c.Reply(note)
			return
		}
	}

	// 对菜谱查询结果排序
	chefs, note = orderChefs(chefs, order)
	if note != "" {
		logger.Info("厨师排序失败:", note)
		_, _ = c.Reply(note)
		return
	}

	// 根据查询结果分页并发送消息
	msg := echoChefsMessage(chefs, order, page, c.GetMessageType() == onebot.MessageTypePrivate)
	logger.Info("发送厨师查询结果:", msg)
	_, _ = c.Reply(msg)
}

// 根据厨师稀有度筛选厨师
func filterChefsByRarity(chefs []database.Chef, rarity int) ([]database.Chef, string) {
	if len(chefs) == 0 {
		return chefs, ""
	}
	result := make([]database.Chef, 0)
	for i := range chefs {
		if chefs[i].Rarity == rarity {
			result = append(result, chefs[i])
		}
	}
	return result, ""
}

// 根据来源筛选厨师
func filterChefsByOrigin(chefs []database.Chef, origin string) ([]database.Chef, string) {
	if len(chefs) == 0 {
		return chefs, ""
	}
	result := make([]database.Chef, 0)
	pattern := strings.ReplaceAll(origin, "%", ".*")
	// 单独增加未入坑礼包查询
	if origin == "仅礼包" || origin == "在售礼包" || origin == "未入坑礼包" {
		pattern = "^限时礼包$"
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		logger.Error("查询正则格式有误", err)
		return nil, "来源筛选格式有误"
	}
	for i := range chefs {
		if re.MatchString(chefs[i].Origin) {
			result = append(result, chefs[i])
		}
	}

	return result, ""
}

// 根据厨师技能筛选厨师
func filterChefsBySkill(chefs []database.Chef, skill string) ([]database.Chef, string) {
	if skill == "" {
		return nil, "你想筛选什么技能呀? 贵客参数格式为「贵客-贵客名」"
	}
	// 处理某些技能关键词
	if s, has := util.WhatPrefixIn(skill, "炒光环", "烤光环", "煮光环", "蒸光环", "炸光环", "切光环", "光环"); has {
		skill = "场上所有厨师" + strings.ReplaceAll(s, "光环", "")
	}
	if s, has := util.WhatPrefixIn(skill, "贵客", "贵宾", "客人", "宾客", "稀客"); has {
		skill = "稀有客人" + strings.TrimLeft(skill, s)
	}
	if strings.HasPrefix(skill, "采集") {
		skill = "探索" + strings.TrimLeft(skill, "采集")
	}
	pattern := strings.ReplaceAll(skill, "%", ".*")
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Sprintf("技能描述格式有误 %v", err)
	}
	result := make([]database.Chef, 0)
	for i := range chefs {
		if re.MatchString(chefs[i].SkillDesc) || re.MatchString(chefs[i].UltimateSkillDesc) {
			result = append(result, chefs[i])
		}
	}
	return result, ""
}

// 根据厨师名或厨师ID筛选厨师
func filterChefsByName(chefs []database.Chef, name string) ([]database.Chef, string) {
	result := make([]database.Chef, 0)
	numId, err := strconv.Atoi(name)
	if err != nil {
		re, err := regexp.Compile(strings.ReplaceAll(name, "%", ".*"))
		if err != nil {
			logger.Error("查询正则格式有误", err)
			return nil, "查询格式有误"
		}
		for i := range chefs {
			if chefs[i].Name == name {
				return []database.Chef{chefs[i]}, ""
			}
			if re.MatchString(chefs[i].Name) {
				result = append(result, chefs[i])
			}
		}
	} else {
		for i := range chefs {
			if chefs[i].ChefId == (numId+2)/3*3 {
				result = append(result, chefs[i])
			}
		}
	}
	return result, ""
}

// 根据排序参数排序厨师
func orderChefs(chefs []database.Chef, order string) ([]database.Chef, string) {
	if len(chefs) == 0 {
		return chefs, ""
	}
	switch order {
	case "图鉴序":
		sort.Slice(chefs, func(i, j int) bool {
			return chefs[i].ChefId < chefs[j].ChefId
		})
	case "稀有度":
		sort.Slice(chefs, func(i, j int) bool {
			return chefs[i].Rarity == chefs[j].Rarity && chefs[i].ChefId < chefs[j].ChefId ||
				chefs[i].Rarity > chefs[j].Rarity
		})
	default:
		return nil, "排序参数有误"
	}
	return chefs, ""
}

// 输出单厨师消息数据
func echoChefMessage(chef database.Chef) string {
	// 尝试寻找图片文件，未找到则按照文字格式发送
	imagePath := fmt.Sprintf("%s/chef/chef_%s_%s.png", config.AppConfig.Resource.Image, chef.GalleryId, strings.ReplaceAll(chef.Name, " ", "_"))
	var msg string
	if has, err := util.PathExists(imagePath); has {
		msg = onebot.GetCQImage(imagePath, "file")
	} else {
		if err != nil {
			logger.Warnf("厨师 %d %s 图片文件可能不存在!", chef.GalleryId, chef.Name, err)
		}
		logger.Infof("未找到厨师 %d %s 图鉴图片, 以文字格式发送数据", chef.GalleryId, chef.Name)
		var gender string
		if chef.Gender == 1 {
			gender = "♂️"
		} else if chef.Gender == 2 {
			gender = "♀️"
		}
		mSkills, err := dao.GetSkillsMap()
		if err != nil {
			logger.Error("查询技能数据出错!", err)
			return e.SystemErrorNote
		}
		ultimateGoals, err := dao.FindQuestsWithIds(chef.UltimateGoals)
		if err != nil {
			logger.Error("查询厨师修炼效果数据出错!", err)
			return e.SystemErrorNote
		}
		goals := ""
		for p, ultimateGoal := range ultimateGoals {
			goals += fmt.Sprintf("\n%d.%s", p+1, ultimateGoal.Goal)
		}
		msg += fmt.Sprintf("%s %s %s\n", chef.GalleryId, chef.Name, gender)
		msg += fmt.Sprintf("%s\n", strings.Repeat("🔥", chef.Rarity))
		msg += fmt.Sprintf("来源: %s\n", chef.Origin)
		msg += fmt.Sprintf("炒:%d 烤:%d 煮:%d\n", chef.Stirfry, chef.Bake, chef.Boil)
		msg += fmt.Sprintf("蒸:%d 炸:%d 切:%d\n", chef.Steam, chef.Fry, chef.Cut)
		msg += fmt.Sprintf("🍖:%d 🍞:%d 🥕:%d 🐟:%d\n", chef.Meat, chef.Flour, chef.Vegetable, chef.Fish)
		msg += fmt.Sprintf("技能:%s\n", mSkills[chef.SkillId].Description)
		msg += fmt.Sprintf("修炼效果:%s\n", mSkills[chef.UltimateSkill].Description)
		msg += fmt.Sprintf("修炼任务:%s", goals)
	}
	return msg
}

// 根据来源和排序参数，输出厨师列表消息数据
func echoChefsMessage(chefs []database.Chef, order string, page int, private bool) string {
	if len(chefs) == 0 {
		return "诶? 似乎查无此厨哦!"
	} else if len(chefs) == 1 {
		return echoChefMessage(chefs[0])
	} else {
		var msg string
		listLength := config.AppConfig.Bot.GroupMsgMaxLen
		if private {
			listLength = config.AppConfig.Bot.PrivateMsgMaxLen
		}
		maxPage := (len(chefs)-1)/listLength + 1
		if page > maxPage {
			page = maxPage
		}
		if len(chefs) > listLength {
			msg += fmt.Sprintf("查询到以下厨师 (%d/%d)", page, maxPage)
		} else {
			msg += "查询到以下厨师"
		}
		for i := (page - 1) * listLength; i < page*listLength && i < len(chefs); i++ {
			orderInfo := getChefInfoWithOrder(chefs[i], order)
			msg += fmt.Sprintf("\n%s %s %s", chefs[i].GalleryId, chefs[i].Name, orderInfo)
		}
		if page < maxPage {
			msg += "\n......"
		}
		return msg
	}
}

// 根据排序参数获取厨师需要输出的信息
func getChefInfoWithOrder(chef database.Chef, order string) string {
	switch order {
	case "图鉴序", "稀有度":
		return strings.Repeat("🔥", chef.Rarity)
	default:
		return ""
	}
}

// GenerateChefImage 根据厨师数据生成图鉴图片
func GenerateChefImage(chef database.ChefData, font *truetype.Font, bgImg image.Image, genderImg image.Image, rarityImg image.Image) (image.Image, error) {
	titleSize := 50 // 标题字体尺寸
	fontSize := 36  // 内容字体尺寸

	img := image.NewRGBA(image.Rect(0, 0, 800, 800))
	// 绘制背景
	draw.Draw(img, img.Bounds(), bgImg, bgImg.Bounds().Min, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(color.RGBA{A: 255}))
	c.SetFontSize(float64(titleSize))

	// 输出厨师头像, 双线性插值算法会对边缘造成影响，去除一点边框
	draw.Draw(img,
		image.Rect(50+2, 118+2, 50+200-2, 118+200-2),
		chef.Avatar,
		image.Point{X: 2, Y: 2},
		draw.Over)

	// 输出图鉴ID与厨师名
	var err error
	if chef.ChefId > 0 {
		_, err = c.DrawString(chef.Name, freetype.Pt(165, 22+titleSize))
		if err != nil {
			return nil, err
		}
		_, err = c.DrawString(fmt.Sprintf("%03d", chef.ChefId), freetype.Pt(45, 18+titleSize))
		if err != nil {
			return nil, err
		}
		c.SetFontSize(float64(25))
		_, err = c.DrawString(fmt.Sprintf("(%03d,%03d)", chef.ChefId-2, chef.ChefId-1), freetype.Pt(30, 70+25))
		if err != nil {
			return nil, err
		}
	} else {
		_, err := c.DrawString(chef.Name, freetype.Pt(50, 22+titleSize))
		if err != nil {
			return nil, err
		}
	}

	// 输出性别
	if genderImg != nil {
		draw.Draw(img,
			image.Rect(490, 30, 490+44, 30+44),
			genderImg,
			image.Point{},
			draw.Over)
	}

	// 输出稀有度
	draw.Draw(img,
		image.Rect(545, 30, 545+240, 30+44),
		rarityImg,
		image.Point{},
		draw.Over)

	c.SetFontSize(float64(fontSize))
	// 输出技法数据
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Stirfry), freetype.Pt(365, 104+fontSize))
	if err != nil {
		return nil, err
	}
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Bake), freetype.Pt(536, 104+fontSize))
	if err != nil {
		return nil, err
	}
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Boil), freetype.Pt(705, 104+fontSize))
	if err != nil {
		return nil, err
	}
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Steam), freetype.Pt(365, 164+fontSize))
	if err != nil {
		return nil, err
	}
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Fry), freetype.Pt(536, 164+fontSize))
	if err != nil {
		return nil, err
	}
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Cut), freetype.Pt(705, 164+fontSize))
	if err != nil {
		return nil, err
	}

	// 输出采集数据
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Meat), freetype.Pt(365, 230+fontSize))
	if err != nil {
		return nil, err
	}
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Flour), freetype.Pt(536, 230+fontSize))
	if err != nil {
		return nil, err
	}
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Vegetable), freetype.Pt(365, 290+fontSize))
	if err != nil {
		return nil, err
	}
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Fish), freetype.Pt(536, 290+fontSize))
	if err != nil {
		return nil, err
	}

	// 输出调料数据
	_, err = c.DrawString(fmt.Sprintf("%d", chef.GetCondimentValue()), freetype.Pt(705, 290+fontSize))
	if err != nil {
		return nil, err
	}

	// 输出来源数据
	_, err = c.DrawString(fmt.Sprintf("%s", chef.Origin), freetype.Pt(150, 365+fontSize))
	if err != nil {
		return nil, err
	}

	// 输出技能数据
	_, err = c.DrawString(fmt.Sprintf("%s", chef.SkillDesc), freetype.Pt(150, 435+fontSize))
	if err != nil {
		return nil, err
	}

	// 输出修炼效果数据
	_, err = c.DrawString(fmt.Sprintf("%s", chef.GetUltimateSkill()), freetype.Pt(150, 505+fontSize))
	if err != nil {
		return nil, err
	}

	// 输出修炼任务数据
	goals := chef.GetUltimateGoals()
	for i, goal := range goals {
		_, err = c.DrawString(goal, freetype.Pt(120, 625+i*50+fontSize))
		if err != nil {
			return nil, err
		}
	}
	return img, nil
}

// GenerateAllChefsImages 生成所有厨师的图鉴图片
func GenerateAllChefsImages(chefs []database.Chef, galleryImg image.Image, imgCSS *gamedata.ImgCSS) error {
	magnification := 4 // 截取的图像相比图鉴网原始图片的放大倍数，图鉴网图片imgCSS给的数据时缩小版图片记录的位置，下载的图片为高清版尺寸为两倍，因此后续计算中取不同的计算倍数

	// 载入字体文件
	font, err := util.LoadFontFile(fmt.Sprintf("%s/%s", config.AppConfig.Resource.Font, "yuan500W.ttf"))
	if err != nil {
		return fmt.Errorf("载入字体文件失败 %v", err)
	}

	resourceImgDir := config.AppConfig.Resource.Image
	chefImgPath := resourceImgDir + "/chef"
	commonImgPath := resourceImgDir + "/common"

	// 放大厨师图鉴图像
	logger.Debugf("厨师图片原始尺寸:%d*%d", galleryImg.Bounds().Dx(), galleryImg.Bounds().Dy())
	galleryImg = resize.Resize(
		uint(galleryImg.Bounds().Dx()*magnification/2.0),
		uint(galleryImg.Bounds().Dy()*magnification/2.0),
		galleryImg, resize.Bilinear)

	// 载入厨师背景图片
	mBgImages := make(map[string]image.Image)
	for _, condimentType := range []string{"Sweet", "Sour", "Spicy", "Salty", "Bitter", "Tasty"} {
		img, err := util.LoadPngImageFile(fmt.Sprintf("%s/chef_%s.png", chefImgPath, condimentType))
		if err != nil {
			return fmt.Errorf("载入厨师背景图片失败 %v", err)
		}
		mBgImages[condimentType] = img
	}

	// 载入厨师性别图片
	mGenderImages := make(map[int]image.Image)
	for _, gender := range []int{0, 1, 2} {
		img, err := util.LoadPngImageFile(fmt.Sprintf("%s/gender_%d.png", chefImgPath, gender))
		if err != nil {
			return fmt.Errorf("载入性别图标失败 %v", err)
		}
		mGenderImages[gender] = img
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

	// 载入任务数据
	mQuests, err := dao.GetQuestsMap()
	if err != nil {
		return fmt.Errorf("载入任务数据出错 %v", err)
	}

	// 逐个绘制厨师图片
	for _, chef := range chefs {
		// 载入与计算厨师信息
		chefImgInfo := imgCSS.ChefImg[chef.ChefId]
		avatarStartX := chefImgInfo.X * magnification
		avatarStartY := chefImgInfo.Y * magnification
		avatar := image.NewRGBA(image.Rect(0, 0, 200, 200))
		draw.Draw(avatar,
			image.Rect(0, 0, 200, 200),
			galleryImg,
			image.Point{X: avatarStartX, Y: avatarStartY},
			draw.Over)
		goals := make([]string, len(chef.UltimateGoals))
		for i := range chef.UltimateGoals {
			goals[i] = mQuests[chef.UltimateGoals[i]].Goal
		}
		chefData := database.ChefData{
			Chef:          chef,
			Avatar:        avatar,
			UltimateGoals: goals,
		}
		// 绘制厨师图片
		img, err := GenerateChefImage(chefData, font, mBgImages[chefData.GetCondimentType()], mGenderImages[chefData.Gender], mRarityImages[chefData.Rarity])
		if err != nil {
			return fmt.Errorf("绘制厨师 %s 的数据出错 %v", chef.Name, err)
		}

		// 以PNG格式保存文件
		err = util.SavePngImage(fmt.Sprintf("%s/chef_%s_%s.png", chefImgPath, chef.GalleryId, strings.ReplaceAll(chef.Name, " ", "_")), img)
		if err != nil {
			return fmt.Errorf("保存厨师 %s 图鉴图片出错 %v", chef.GalleryId, err)
		}
	}
	return nil
}
