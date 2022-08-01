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
	"bytes"
	"fmt"
	"github.com/golang/freetype"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func ChefQuery(c *scheduler.Context) {
	if strings.TrimSpace(c.PretreatedMessage) == "" {
		_, _ = c.Reply(chefHelp())
		return
	}

	order := "稀有度"
	page := 1
	var note string
	chefs := make([]database.Chef, 0)
	err := dao.DB.Find(&chefs)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_, _ = c.Reply(e.SystemErrorNote)
	}
	args := strings.Split(strings.TrimSpace(c.PretreatedMessage), " ")
	argCount := 0
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
				pageNum, err := strconv.Atoi(arg[1:])
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
		argCount++
	}

	if argCount == 0 {
		_, _ = c.Reply(recipeHelp())
		return
	}

	// 对厨师查询结果排序
	chefs, note = orderChefs(chefs, order)
	if note != "" {
		logger.Info("厨师查询失败:", note)
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
	pattern := ".*" + strings.ReplaceAll(origin, "%", ".*") + ".*"
	// 单独增加在售礼包查询
	if origin == "仅礼包" || origin == "在售礼包" {
		for i := range chefs {
			if chefs[i].Origin == "限时礼包" {
				result = append(result, chefs[i])
			}
		}
	} else {
		for i := range chefs {
			re := regexp.MustCompile(pattern)
			if re.MatchString(chefs[i].Origin) {
				result = append(result, chefs[i])
			}
		}
	}
	return result, ""
}

// 根据厨师技能筛选厨师
func filterChefsBySkill(chefs []database.Chef, skill string) ([]database.Chef, string) {
	// 处理某些技能关键词
	if s, has := util.WhatPrefixIn(skill, "炒光环", "烤光环", "煮光环", "蒸光环", "炸光环", "切光环", "光环"); has {
		skill = "场上所有厨师" + strings.ReplaceAll(s, "光环", "") + "%" + strings.ReplaceAll(skill, s, "")
	}
	if s, has := util.WhatPrefixIn(skill, "贵客", "贵宾", "客人", "宾客", "稀客"); has {
		skill = "稀有客人" + "%" + strings.ReplaceAll(skill, s, "")
	}
	if strings.HasPrefix(skill, "采集") {
		skill = "探索" + "%" + strings.ReplaceAll(skill, "采集", "")
	}
	result := make([]database.Chef, 0)
	skills := make(map[int]database.Skill)
	err := dao.DB.Where("description like ?", "%"+skill+"%").Find(&skills)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		return result, e.SystemErrorNote
	}
	for i := range chefs {
		if _, ok := skills[chefs[i].SkillId]; ok {
			result = append(result, chefs[i])
			continue
		}
		if _, ok := skills[chefs[i].UltimateSkill]; ok {
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
		pattern := ".*" + strings.ReplaceAll(name, "%", ".*") + ".*"
		for i := range chefs {
			re := regexp.MustCompile(pattern)
			if re.MatchString(chefs[i].Name) {
				result = append(result, chefs[i])
			}
		}
	} else {
		if numId%3 != 0 {
			numId = numId + (3 - numId%3)
		}
		galleryId := fmt.Sprintf("%03d", numId)
		for i := range chefs {
			if chefs[i].GalleryId == galleryId {
				result = append(result, chefs[i])
			}
		}
	}
	return result, ""
}

type chefWrapper struct {
	chef     []database.Chef
	chefLess func(p *database.Chef, q *database.Chef) bool
}

func (w chefWrapper) Len() int {
	return len(w.chef)
}

func (w chefWrapper) Swap(i int, j int) {
	w.chef[i], w.chef[j] = w.chef[j], w.chef[i]
}

func (w chefWrapper) Less(i int, j int) bool {
	return w.chefLess(&w.chef[i], &w.chef[j])
}

// 根据排序参数排序厨师
func orderChefs(chefs []database.Chef, order string) ([]database.Chef, string) {
	if len(chefs) == 0 {
		return chefs, ""
	}
	switch order {
	case "图鉴序":
		sort.Sort(chefWrapper{chefs, func(m, n *database.Chef) bool {
			return m.ChefId < n.ChefId
		}})
	case "稀有度":
		sort.Sort(chefWrapper{chefs, func(m, n *database.Chef) bool {
			if m.Rarity == n.Rarity {
				return m.ChefId < n.ChefId
			} else {
				return m.Rarity > n.Rarity
			}
		}})
	default:
		return nil, "排序参数有误"
	}
	return chefs, ""
}

// 输出单厨师消息数据
func echoChefMessage(chef database.Chef) string {
	// 尝试寻找图片文件，未找到则按照文字格式发送
	resourceImageDir := config.AppConfig.Resource.Image + "/chef"
	imagePath := fmt.Sprintf("%s/chef_%s.png", resourceImageDir, chef.GalleryId)
	logger.Debug("imagePath:", imagePath)
	var msg string
	if has, err := util.PathExists(imagePath); has {
		msg = onebot.GetCQImage(imagePath, "file")
	} else {
		if err != nil {
			logger.Debugf("无法确定文件是否存在!", err)
		}
		logger.Info("未找到厨师图鉴图片, 以文字格式发送数据")
		var gender string
		if chef.Gender == 1 {
			gender = "♂️"
		} else if chef.Gender == 2 {
			gender = "♀️"
		}
		rarity := ""
		for i := 0; i < chef.Rarity; i++ {
			rarity += "🔥"
		}
		skill := new(database.Skill)
		_, err = dao.DB.Where("skill_id = ?", chef.SkillId).Get(skill)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return e.SystemErrorNote
		}
		ultimateSkill := new(database.Skill)
		_, err = dao.DB.Where("skill_id = ?", chef.UltimateSkill).Get(ultimateSkill)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return e.SystemErrorNote
		}
		ultimateGoals := make([]database.Quest, 0)
		err = dao.DB.In("quest_id", chef.UltimateGoals).Find(&ultimateGoals)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			return e.SystemErrorNote
		}
		goals := ""
		for p, ultimateGoal := range ultimateGoals {
			goals += fmt.Sprintf("\n%d.%s", p+1, ultimateGoal.Goal)
		}
		msg += fmt.Sprintf("%s %s %s\n", chef.GalleryId, chef.Name, gender)
		msg += fmt.Sprintf("%s\n", rarity)
		msg += fmt.Sprintf("来源: %s\n", chef.Origin)
		msg += fmt.Sprintf("炒:%d 烤:%d 煮:%d\n", chef.Stirfry, chef.Bake, chef.Boil)
		msg += fmt.Sprintf("蒸:%d 炸:%d 切:%d\n", chef.Steam, chef.Fry, chef.Cut)
		msg += fmt.Sprintf("🍖:%d 🍞:%d 🥕:%d 🐟:%d\n", chef.Meat, chef.Flour, chef.Vegetable, chef.Fish)
		msg += fmt.Sprintf("技能:%s\n", skill.Description)
		msg += fmt.Sprintf("修炼效果:%s\n", ultimateSkill.Description)
		msg += fmt.Sprintf("修炼任务:%s", goals)
	}
	return msg
}

// 根据来源和排序参数，输出厨师列表消息数据
func echoChefsMessage(chefs []database.Chef, order string, page int, private bool) string {
	if len(chefs) == 0 {
		return "哎呀，好像找不到呢!"
	} else if len(chefs) == 1 {
		return echoChefMessage(chefs[0])
	} else {
		logger.Debug("查询到多个厨师")
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
			msg += fmt.Sprintf("查询到以下厨师: (%d/%d)\n", page, maxPage)
		} else {
			msg += "查询到以下厨师:\n"
		}
		for i := (page - 1) * listLength; i < page*listLength && i < len(chefs); i++ {
			orderInfo := getChefInfoWithOrder(chefs[i], order)
			msg += fmt.Sprintf("%s %s %s", chefs[i].GalleryId, chefs[i].Name, orderInfo)
			if i < page*listLength-1 && i < len(chefs)-1 {
				msg += "\n"
			}
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
	case "图鉴序":
		msg := ""
		for i := 0; i < chef.Rarity; i++ {
			msg += "🔥"
		}
		return msg
	case "稀有度":
		msg := ""
		for i := 0; i < chef.Rarity; i++ {
			msg += "🔥"
		}
		return msg
	default:
		return ""
	}
}

func ChefInfoToImage(chefs []database.Chef, imgURL string, imgCSS *gamedata.ImgCSS) error {
	dx := 800          // 图鉴背景图片的宽度
	dy := 800          // 图鉴背景图片的高度
	magnification := 4 // 截取的图像相比图鉴网原始图片的放大倍数
	titleSize := 50    // 标题字体尺寸
	fontSize := 36     // 内容字体尺寸
	fontDPI := 72.0    // dpi

	// 获取字体文件
	resourceFontDir := config.AppConfig.Resource.Font
	fontFile := resourceFontDir + "/yuan500W.ttf"
	//读字体数据
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}
	fontColor := color.RGBA{A: 255}
	// 从图鉴网下载头像图鉴总图
	resourceImgDir := config.AppConfig.Resource.Image
	chefImgPath := resourceImgDir + "/chef"
	galleryImagePath := chefImgPath + "/chef_gallery.png"
	r, err := http.Get(imgURL)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	out, err := os.Create(galleryImagePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, bytes.NewReader(body))
	if err != nil {
		return err
	}

	galleryImg, err := png.Decode(bytes.NewReader(body))
	if err != nil {
		return err
	}

	// 放大厨师图鉴图像
	galleryImg = resize.Resize(
		uint(galleryImg.Bounds().Dx()*magnification/2.0),
		uint(galleryImg.Bounds().Dy()*magnification/2.0),
		galleryImg, resize.Bilinear)

	for _, chef := range chefs {
		condiment := 0
		condimentType := "Sweet"
		if chef.Sweet > 0 {
			condiment = chef.Sweet
			condimentType = "Sweet"
		} else if chef.Sour > 0 {
			condiment = chef.Sour
			condimentType = "Sour"
		} else if chef.Spicy > 0 {
			condiment = chef.Spicy
			condimentType = "Spicy"
		} else if chef.Salty > 0 {
			condiment = chef.Salty
			condimentType = "Salty"
		} else if chef.Bitter > 0 {
			condiment = chef.Bitter
			condimentType = "Bitter"
		} else if chef.Tasty > 0 {
			condiment = chef.Tasty
			condimentType = "Tasty"
		}
		bgFile, err := os.Open(fmt.Sprintf("%s/chef_%s.png", chefImgPath, condimentType))
		if err != nil {
			return err
		}
		img := image.NewRGBA(image.Rect(0, 0, dx, dy))
		bg, _ := png.Decode(bgFile)
		bgFile.Close()

		draw.Draw(img, img.Bounds(), bg, bg.Bounds().Min, draw.Src)

		c := freetype.NewContext()
		c.SetDPI(fontDPI)
		c.SetFont(font)
		c.SetClip(img.Bounds())
		c.SetDst(img)
		c.SetSrc(image.NewUniform(fontColor))
		c.SetFontSize(float64(titleSize))

		// 输出厨师头像, 双线性插值算法会对边缘造成影响，去除一点边框
		chefImgInfo := imgCSS.ChefImg[chef.ChefId]
		avatarStartX := chefImgInfo.X * magnification
		avatarStartY := chefImgInfo.Y * magnification
		draw.Draw(img,
			image.Rect(50+2, 118+2, 50+200-2, 118+200-2),
			galleryImg,
			image.Point{X: avatarStartX + 2, Y: avatarStartY + 2},
			draw.Over)

		// 输出图鉴ID与厨师名
		pt := freetype.Pt(165, 22+titleSize)
		_, err = c.DrawString(fmt.Sprintf("%s", chef.Name), pt)
		if err != nil {
			return err
		}

		pt = freetype.Pt(45, 18+titleSize)
		_, err = c.DrawString(fmt.Sprintf("%03d", chef.ChefId), pt)
		if err != nil {
			return err
		}
		c.SetFontSize(float64(25))
		pt = freetype.Pt(30, 70+25)
		_, err = c.DrawString(fmt.Sprintf("(%03d,%03d)", chef.ChefId-2, chef.ChefId-1), pt)
		if err != nil {
			return err
		}

		// 输出性别
		genderFile, err := os.Open(fmt.Sprintf("%s/gender_%d.png", chefImgPath, chef.Gender))
		if err != nil {
			return err
		}
		genderImg, _ := png.Decode(genderFile)
		genderFile.Close()
		draw.Draw(img,
			image.Rect(490, 30, 490+44, 30+44),
			genderImg,
			image.Point{},
			draw.Over)

		// 输出稀有度
		rarityFile, err := os.Open(fmt.Sprintf("%s/rarity_%d.png", chefImgPath, chef.Rarity))
		if err != nil {
			return err
		}
		rarityImg, _ := png.Decode(rarityFile)
		rarityFile.Close()
		draw.Draw(img,
			image.Rect(545, 30, 545+240, 30+44),
			rarityImg,
			image.Point{},
			draw.Over)

		c.SetFontSize(float64(fontSize))
		// 输出技法数据
		pt = freetype.Pt(365, 104+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%d", chef.Stirfry), pt)
		if err != nil {
			return err
		}
		pt = freetype.Pt(536, 104+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%d", chef.Bake), pt)
		if err != nil {
			return err
		}
		pt = freetype.Pt(705, 104+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%d", chef.Boil), pt)
		if err != nil {
			return err
		}
		pt = freetype.Pt(365, 164+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%d", chef.Steam), pt)
		if err != nil {
			return err
		}
		pt = freetype.Pt(536, 164+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%d", chef.Fry), pt)
		if err != nil {
			return err
		}
		pt = freetype.Pt(705, 164+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%d", chef.Cut), pt)
		if err != nil {
			return err
		}

		// 输出采集数据
		pt = freetype.Pt(365, 230+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%d", chef.Meat), pt)
		if err != nil {
			return err
		}
		pt = freetype.Pt(536, 230+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%d", chef.Flour), pt)
		if err != nil {
			return err
		}
		pt = freetype.Pt(365, 290+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%d", chef.Vegetable), pt)
		if err != nil {
			return err
		}
		pt = freetype.Pt(536, 290+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%d", chef.Fish), pt)
		if err != nil {
			return err
		}

		// 输出调料数据
		pt = freetype.Pt(705, 290+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%d", condiment), pt)
		if err != nil {
			return err
		}

		// 输出来源数据
		pt = freetype.Pt(150, 365+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%s", chef.Origin), pt)
		if err != nil {
			return err
		}

		// 输出技法数据
		skill := new(database.Skill)
		_, err = dao.DB.Where("skill_id = ?", chef.SkillId).Get(skill)
		if err != nil {
			return err
		}
		pt = freetype.Pt(150, 435+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%s", skill.Description), pt)
		if err != nil {
			return err
		}

		// 输出修炼效果数据
		ultimateSkill := new(database.Skill)
		_, err = dao.DB.Where("skill_id = ?", chef.UltimateSkill).Get(ultimateSkill)
		if err != nil {
			return err
		}
		pt = freetype.Pt(150, 505+fontSize)
		if ultimateSkill.Description == "" {
			ultimateSkill.Description = "暂无"
		}
		_, err = c.DrawString(fmt.Sprintf("%s", ultimateSkill.Description), pt)
		if err != nil {
			return err
		}

		// 输出修炼任务数据
		ultimateGoals := make([]database.Quest, 0)
		err = dao.DB.In("quest_id", chef.UltimateGoals).Find(&ultimateGoals)
		if err != nil {
			return err
		}
		for i := 0; i < 3; i++ {
			pt = freetype.Pt(120, 625+i*50+fontSize)
			if len(ultimateGoals)-1 < i {
				_, err = c.DrawString(fmt.Sprintf("暂无"), pt)
				if err != nil {
					return err
				}
			} else {
				_, err = c.DrawString(fmt.Sprintf("%s", ultimateGoals[i].Goal), pt)
				if err != nil {
					return err
				}

			}
		}

		// 以PNG格式保存文件
		dst, err := os.Create(fmt.Sprintf("%s/chef_%s.png", chefImgPath, chef.GalleryId))
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
