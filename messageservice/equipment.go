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
	"image/png"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func EquipmentQuery(c *scheduler.Context) {
	if strings.TrimSpace(c.PretreatedMessage) == "" {
		_, _ = c.Reply(equipmentHelp())
		return
	}

	order := "稀有度"
	page := 1
	var note string
	equips := make([]database.Equip, 0)
	err := dao.DB.Find(&equips)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_, _ = c.Reply(e.SystemErrorNote)
	}
	args := strings.Split(strings.TrimSpace(c.PretreatedMessage), " ")
	for _, arg := range args {
		if arg == "" {
			continue
		}
		switch arg {
		case "图鉴序", "稀有度":
			order = arg
		default:
			if util.HasPrefixIn(arg, "来源") {
				origin := strings.Split(arg, "-")
				if len(origin) > 1 {
					equips, note = filterEquipsByOrigin(equips, origin[1])
				}
			} else if util.HasPrefixIn(arg, "技能", "效果", "功能") {
				skill := strings.Split(arg, "-")
				if len(skill) > 1 {
					equips, note = filterEquipsBySkill(equips, strings.Join(skill[1:], "-"))
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
				equips, note = filterEquipsByName(equips, arg)
			}
		}

		if note != "" {
			logger.Info("厨具查询失败:", note)
			_, _ = c.Reply(note)
			return
		}
	}

	// 对厨具查询结果排序
	equips, note = orderEquips(equips, order)
	if note != "" {
		logger.Info("厨具查询失败:", note)
		_, _ = c.Reply(note)
		return
	}
	// 根据结果翻页并发送消息
	msg := echoEquipsMessage(equips, order, page, c.GetMessageType() == onebot.MessageTypePrivate)
	logger.Info("发送菜谱查询结果:", msg)
	_, _ = c.Reply(msg)
}

// 根据厨具名或厨具ID筛选厨具
func filterEquipsByName(equips []database.Equip, name string) ([]database.Equip, string) {
	result := make([]database.Equip, 0)
	numId, err := strconv.Atoi(name)
	if err != nil {
		pattern := ".*" + strings.ReplaceAll(name, "%", ".*") + ".*"
		re, err := regexp.Compile(pattern)
		if err != nil {
			logger.Error("查询正则格式有误", err)
			return nil, "查询格式有误"
		}
		for i := range equips {
			if equips[i].Name == name {
				return []database.Equip{equips[i]}, ""
			}
			if re.MatchString(equips[i].Name) {
				result = append(result, equips[i])
			}
		}
	} else {
		galleryId := fmt.Sprintf("%03d", numId)
		for i := range equips {
			if equips[i].GalleryId == galleryId {
				result = append(result, equips[i])
			}
		}
	}
	return result, ""
}

// 根据来源筛选厨具
func filterEquipsByOrigin(equips []database.Equip, origin string) ([]database.Equip, string) {
	if len(equips) == 0 {
		return equips, ""
	}
	result := make([]database.Equip, 0)
	pattern := ".*" + strings.ReplaceAll(origin, "%", ".*") + ".*"
	re, err := regexp.Compile(pattern)
	if err != nil {
		logger.Error("查询正则格式有误", err)
		return nil, "查询格式有误"
	}
	for i := range equips {
		if re.MatchString(equips[i].Origin) {
			result = append(result, equips[i])
		}
	}
	return result, ""
}

// 根据厨具技能筛选厨具
func filterEquipsBySkill(equips []database.Equip, skill string) ([]database.Equip, string) {
	// 处理某些技能关键词
	if s, has := util.WhatPrefixIn(skill, "贵客", "贵宾", "客人", "宾客", "稀客"); has {
		skill = "稀有客人" + "%" + strings.ReplaceAll(skill, s, "")
	}
	result := make([]database.Equip, 0)
	skills := make(map[int]database.Skill)
	err := dao.DB.Where("description like ?", "%"+skill+"%").Find(&skills)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		return result, e.SystemErrorNote
	}
	for i := range equips {
		for _, skillId := range equips[i].Skills {
			if _, ok := skills[skillId]; ok {
				result = append(result, equips[i])
				break
			}
		}
	}
	return result, ""
}

type equipWrapper struct {
	equip     []database.Equip
	equipLess func(p *database.Equip, q *database.Equip) bool
}

func (w equipWrapper) Len() int {
	return len(w.equip)
}

func (w equipWrapper) Swap(i int, j int) {
	w.equip[i], w.equip[j] = w.equip[j], w.equip[i]
}

func (w equipWrapper) Less(i int, j int) bool {
	return w.equipLess(&w.equip[i], &w.equip[j])
}

// 根据排序参数排序厨具
func orderEquips(equips []database.Equip, order string) ([]database.Equip, string) {
	if len(equips) == 0 {
		return equips, ""
	}
	switch order {
	case "图鉴序":
		sort.Sort(equipWrapper{equips, func(m, n *database.Equip) bool {
			return m.EquipId < n.EquipId
		}})
	case "稀有度":
		sort.Sort(equipWrapper{equips, func(m, n *database.Equip) bool {
			if m.Rarity == n.Rarity {
				return m.EquipId < n.EquipId
			} else {
				return m.Rarity > n.Rarity
			}
		}})
	default:
		return nil, "排序参数有误"
	}
	return equips, ""
}

// 根据排序参数获取厨具需要输出的信息
func getEquipInfoWithOrder(equip database.Equip, order string) string {
	switch order {
	case "图鉴序":
		msg := ""
		for i := 0; i < equip.Rarity; i++ {
			msg += "🔥"
		}
		return msg
	case "稀有度":
		msg := ""
		for i := 0; i < equip.Rarity; i++ {
			msg += "🔥"
		}
		return msg
	default:
		return ""
	}
}

// 根据来源和排序参数，输出厨具消息列表
func echoEquipsMessage(equips []database.Equip, order string, page int, private bool) string {
	if len(equips) == 0 {
		return "哎呀，好像找不到呢!"
	} else if len(equips) == 1 {
		return echoEquipMessage(equips[0])
	} else {
		logger.Debug("查询到多个厨具")
		var msg string
		listLength := config.AppConfig.Bot.GroupMsgMaxLen
		if private {
			listLength = config.AppConfig.Bot.PrivateMsgMaxLen
		}
		maxPage := (len(equips)-1)/listLength + 1
		if page > maxPage {
			page = maxPage
		}
		if len(equips) > listLength {
			msg += fmt.Sprintf("查询到以下厨具: (%d/%d)\n", page, maxPage)
		} else {
			msg += "查询到以下厨具:\n"
		}
		for i := (page - 1) * listLength; i < page*listLength && i < len(equips); i++ {
			orderInfo := getEquipInfoWithOrder(equips[i], order)
			msg += fmt.Sprintf("%s %s %s", equips[i].GalleryId, equips[i].Name, orderInfo)
			if i < page*listLength-1 && i < len(equips)-1 {
				msg += "\n"
			}
		}
		if page < maxPage {
			msg += "\n......"
		}
		return msg
	}
}

// 输出单厨具消息数据
func echoEquipMessage(equip database.Equip) string {
	resourceImageDir := config.AppConfig.Resource.Image + "/equip"
	imagePath := fmt.Sprintf("%s/equip_%s.png", resourceImageDir, equip.GalleryId)
	logger.Debug("imagePath:", imagePath)
	var msg string
	if has, err := util.PathExists(imagePath); has {
		logger.Debugf("存在厨具图片文件, 返回图片数据")
		msg = onebot.GetCQImage(imagePath, "file")
	} else {
		if err != nil {
			logger.Debugf("无法确定文件是否存在, 返回文字数据", err)
		}
		rarity := ""
		for i := 0; i < equip.Rarity; i++ {
			rarity += "🔥"
		}
		skills := ""
		for p, skillId := range equip.Skills {
			skill := new(database.Skill)
			has, err := dao.DB.Where("skill_id = ?", skillId).Get(skill)
			if err != nil {
				logger.Error("查询数据库出错!", err)
				return e.SystemErrorNote
			}
			if has {
				skills += skill.Description
				if p != len(equip.Skills)-1 {
					skills += ","
				}
			}
		}
		msg += fmt.Sprintf("%s %s\n", equip.GalleryId, equip.Name)
		msg += fmt.Sprintf("%s\n", rarity)
		msg += fmt.Sprintf("来源: %s\n", equip.Origin)
		msg += fmt.Sprintf("效果: %s", skills)
	}
	return msg
}

// GenerateEquipmentImage 根据厨具数据生成单个厨具图鉴图片
func GenerateEquipmentImage(equip database.EquipData, font *truetype.Font, bgImg image.Image, rarityImg image.Image, mSkillImages map[string]image.Image) (image.Image, error) {
	titleSize := 42 // 标题字体尺寸
	fontSize := 28  // 内容字体尺寸

	img := image.NewRGBA(image.Rect(0, 0, 800, 300))
	draw.Draw(img, img.Bounds(), bgImg, bgImg.Bounds().Min, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(color.RGBA{A: 255}))

	//	绘制ID与厨具名
	c.SetFontSize(float64(titleSize))
	_, err := c.DrawString(fmt.Sprintf("%s %s", equip.GalleryId, equip.Name), freetype.Pt(30, 16+titleSize))

	// 绘制稀有度
	draw.Draw(img,
		image.Rect(530, 16, 530+240, 16+44),
		rarityImg,
		image.Point{},
		draw.Over)

	//	绘制厨具图鉴图片
	width := equip.Avatar.Bounds().Dx()
	height := equip.Avatar.Bounds().Dy()
	draw.Draw(img,
		image.Rect(30+210/2-width/2, 75+210/2-height/2, 30+210/2+width/2, 75+210/2+height/2),
		equip.Avatar,
		image.Point{},
		draw.Over)

	//	输出来源数据
	c.SetFontSize(float64(32))
	_, err = c.DrawString(fmt.Sprintf("%s", equip.Origin), freetype.Pt(350, 75+32))
	if err != nil {
		return nil, err
	}

	// 输出技法效果数据
	c.SetFontSize(float64(fontSize))
	for i, skill := range equip.Skills {
		skillImg, ok := mSkillImages[skill.Effects[0].Type]
		if !ok {
			skillImg = mSkillImages["Skill"]
		}
		draw.Draw(img,
			image.Rect(270, 136+i*50, 270+60, 136+i*50+40),
			skillImg,
			image.Point{},
			draw.Over)
		_, err = c.DrawString(fmt.Sprintf("%s", skill.Description), freetype.Pt(320, 138+i*50+fontSize))
		if err != nil {
			return nil, err
		}
	}
	return img, nil
}

func GenerateAllEquipmentsImages(equips []database.Equip, galleryImg image.Image, imgCSS *gamedata.ImgCSS) error {
	magnification := 4 // 截取的图像相比图鉴网原始图片的放大倍数
	// 加载字体文件
	resourceFontDir := config.AppConfig.Resource.Font
	fontPath := "yuan500W.ttf"
	fontFile := fmt.Sprintf("%s/%s", resourceFontDir, fontPath)
	//读字体数据
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	resourceImgDir := config.AppConfig.Resource.Image
	commonImgPath := resourceImgDir + "/common"
	equipImgPath := resourceImgDir + "/equip"

	// 放大厨具图鉴图像
	logger.Debugf("厨具图片原始尺寸:%d*%d", galleryImg.Bounds().Dx(), galleryImg.Bounds().Dy())
	galleryImg = resize.Resize(
		uint(galleryImg.Bounds().Dx()*magnification/2),
		uint(galleryImg.Bounds().Dy()*magnification/2),
		galleryImg, resize.Bilinear)

	// 载入背景图片
	bgFile, err := os.Open(fmt.Sprintf("%s/equip_bg.png", equipImgPath))
	if err != nil {
		return err
	}
	bgImg := image.NewRGBA(image.Rect(0, 0, 800, 300))
	bg, _ := png.Decode(bgFile)
	_ = bgFile.Close()

	draw.Draw(bgImg, bgImg.Bounds(), bg, bg.Bounds().Min, draw.Src)

	// 载入稀有度图片
	mRarityImages := make(map[int]image.Image)
	for _, rarity := range []int{1, 2, 3} {
		rarityFile, err := os.Open(fmt.Sprintf("%s/rarity_%d.png", commonImgPath, rarity))
		if err != nil {
			return err
		}
		img := image.NewRGBA(image.Rect(0, 0, 240, 44))
		bg, _ := png.Decode(rarityFile)
		_ = rarityFile.Close()
		draw.Draw(img, img.Bounds(), bg, bg.Bounds().Min, draw.Over)
		mRarityImages[rarity] = img
	}

	// 载入技能效果图标
	mSkillImages, err := loadSkillIcons(commonImgPath)
	if err != nil {
		return err
	}

	for _, equip := range equips {
		// 计算与载入厨具信息
		equipImgInfo := imgCSS.EquipImg[equip.EquipId]
		avatarStartX := equipImgInfo.X * magnification
		avatarStartY := equipImgInfo.Y * magnification
		avatarWidth := equipImgInfo.Width * magnification
		avatarHeight := equipImgInfo.Height * magnification

		avatar := image.NewRGBA(image.Rect(0, 0, avatarWidth, avatarHeight))
		draw.Draw(avatar,
			image.Rect(0, 0, avatarWidth, avatarHeight),
			galleryImg,
			image.Point{X: avatarStartX, Y: avatarStartY},
			draw.Over)

		skills, err := dao.GetSkillsByIds(equip.Skills)
		if err != nil {
			logger.Errorf("查询厨具 %s 技能数据失败, 技能id %v, err: %v", equip.Name, equip.Skills, err)
			continue
		}

		equipData := database.EquipData{
			Equip:  equip,
			Avatar: avatar,
			Skills: skills,
		}

		img, err := GenerateEquipmentImage(equipData, font, bgImg, mRarityImages[equip.Rarity], mSkillImages)
		if err != nil {
			return fmt.Errorf("绘制厨具 %s 的数据出错 %v", equip.Name, err)
		}

		// 以PNG格式保存文件
		dst, err := os.Create(fmt.Sprintf("%s/equip_%s.png", equipImgPath, equip.GalleryId))
		if err != nil {
			return err
		}
		err = png.Encode(dst, img)
		if err != nil {
			return err
		}
		_ = dst.Close()
	}
	return nil
}

// loadSkillIcons 加载技能效果图标
func loadSkillIcons(basePath string) (map[string]image.Image, error) {
	m := make(map[string]image.Image)
	for _, skill := range []string{
		"Stirfry", "Bake", "Boil", "Steam", "Fry", "Knife", "Sweet", "Sour", "Spicy", "Salty",
		"Bitter", "Tasty", "Meat", "Creation", "Vegetable", "Fish", "OpenTime", "Skill"} {
		iconFile, err := os.Open(fmt.Sprintf("%s/icon_%s.png", basePath, strings.ToLower(skill)))
		if err != nil {
			return nil, err
		}
		img, _ := png.Decode(iconFile)
		_ = iconFile.Close()
		img = resize.Resize(0, 40, img, resize.MitchellNetravali)
		m[skill] = img
		m["Use"+skill] = img
	}
	return m, nil
}
