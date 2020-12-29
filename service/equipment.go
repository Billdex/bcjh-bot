package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/config"
	"bcjh-bot/model/database"
	"bcjh-bot/model/gamedata"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
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
)

func EquipmentQuery(c *onebot.Context, args []string) {
	logger.Info("厨具查询，参数:", args)
	if len(args) == 0 {
		err := bot.SendMessage(c, equipmentHelp())
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
		return
	}
	if args[0] == "%" {
		err := bot.SendMessage(c, "参数有误!")
		if err != nil {
			logger.Error("发送信息失败!", err)
		}
		return
	}

	equips := make([]database.Equip, 0)
	err := database.DB.Where("gallery_id = ?", args[0]).Asc("gallery_id").Find(&equips)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_ = bot.SendMessage(c, "查询数据失败!")
		return
	}
	if len(equips) == 0 {
		err = database.DB.Where("name like ?", "%"+args[0]+"%").Asc("gallery_id").Find(&equips)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			_ = bot.SendMessage(c, "查询数据失败!")
			return
		}
	}

	var msg string
	if len(equips) == 0 {
		msg = "哎呀，好像找不到呢!"
	} else if len(equips) == 1 {
		resourceImageDir := config.AppConfig.Resource.Image + "/equip"
		imagePath := fmt.Sprintf("%s/equip_%s.png", resourceImageDir, equips[0].GalleryId)
		logger.Debug("imagePath:", imagePath)
		if has, err := util.PathExists(imagePath); has {
			logger.Debugf("存在厨具图片文件, 返回图片数据")
			msg = bot.GetCQImage(imagePath, "file")
		} else {
			if err != nil {
				logger.Debugf("无法确定文件是否存在, 返回文字数据", err)
			}
			equip := equips[0]
			rarity := ""
			for i := 0; i < equip.Rarity; i++ {
				rarity += "🔥"
			}
			skills := ""
			for p, skillId := range equip.Skills {
				skill := new(database.Skill)
				has, err := database.DB.Where("skill_id = ?", skillId).Get(skill)
				if err != nil {
					logger.Error("查询数据库出错!", err)
					_ = bot.SendMessage(c, "查询数据失败!")
					return
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
	} else {
		msg = "查询到以下厨具:\n"
		for p, equip := range equips {
			msg += fmt.Sprintf("%s %s", equip.GalleryId, equip.Name)
			if p != len(equips)-1 {
				msg += "\n"
				if p == util.MaxQueryListLength-1 {
					msg += "......"
					break
				}
			}
		}
	}

	logger.Debug("msg:", msg)
	err = bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}

func EquipmentInfoToImage(equips []database.Equip, imgCSS *gamedata.ImgCSS) error {
	dx := 800          // 图鉴背景图片的宽度
	dy := 300          // 图鉴背景图片的高度
	magnification := 4 // 截取的图像相比图鉴网原始图片的放大倍数
	titleSize := 48    // 标题字体尺寸
	fontSize := 32     // 内容字体尺寸
	fontDPI := 72.0    // dpi
	// 需要使用的字体文件
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

	// 从图鉴网下载头像图鉴总图
	resourceImgDir := config.AppConfig.Resource.Image
	equipImgPath := resourceImgDir + "/equip"
	galleryImagePath := equipImgPath + "/equip_gallery.png"
	r, err := http.Get(util.EquipImageRetinaURL)
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
	defer out.Close()
	_, err = io.Copy(out, bytes.NewReader(body))
	if err != nil {
		return err
	}

	galleryImg, err := png.Decode(bytes.NewReader(body))
	if err != nil {
		return err
	}

	// 放大厨具图鉴图像
	logger.Debugf("厨具图片尺寸:%d*%d", galleryImg.Bounds().Dx(), galleryImg.Bounds().Dy())
	galleryImg = resize.Resize(
		uint(galleryImg.Bounds().Dx()*magnification/2),
		uint(galleryImg.Bounds().Dy()*magnification/2),
		galleryImg, resize.Bilinear)

	for _, equip := range equips {
		// 绘制背景色
		bgFile, err := os.Open(fmt.Sprintf("%s/equip_bg.png", equipImgPath))
		if err != nil {
			return err
		}
		defer bgFile.Close()
		img := image.NewRGBA(image.Rect(0, 0, dx, dy))
		bg, _ := png.Decode(bgFile)

		draw.Draw(img, img.Bounds(), bg, bg.Bounds().Min, draw.Src)

		c := freetype.NewContext()
		c.SetDPI(fontDPI)
		c.SetFont(font)
		c.SetClip(img.Bounds())
		c.SetDst(img)
		fontColor := color.RGBA{0, 0, 0, 255}
		c.SetSrc(image.NewUniform(fontColor))

		//	绘制ID与厨具名
		c.SetFontSize(float64(titleSize))
		pt := freetype.Pt(30, 10+titleSize)
		_, err = c.DrawString(fmt.Sprintf("%s %s", equip.GalleryId, equip.Name), pt)

		// 绘制稀有度
		rarityFile, err := os.Open(fmt.Sprintf("%s/rarity_%d.png", equipImgPath, equip.Rarity))
		if err != nil {
			return err
		}
		defer rarityFile.Close()
		rarityImg, _ := png.Decode(rarityFile)
		draw.Draw(img,
			image.Rect(545, 10, 545+240, 30+44),
			rarityImg,
			image.ZP,
			draw.Over)

		//	绘制厨具图鉴图片
		equipImgInfo := imgCSS.EquipImg[equip.EquipId]
		avatarStartX := equipImgInfo.X * magnification
		avatarStartY := equipImgInfo.Y * magnification
		avatarWidth := equipImgInfo.Width * magnification
		avatarHeight := equipImgInfo.Height * magnification
		draw.Draw(img,
			image.Rect(30+2+210/2-avatarWidth/2, 70+2+210/2-avatarHeight/2, 30+210/2+avatarWidth/2-2, 70+210/2+avatarHeight/2-2),
			galleryImg,
			image.Point{X: avatarStartX + 2, Y: avatarStartY + 2},
			draw.Over)

		c.SetFontSize(float64(fontSize))
		//	输出来源数据
		pt = freetype.Pt(270, 75+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%s", equip.Origin), pt)
		if err != nil {
			return err
		}
		// 输出技法效果数据
		skills := make([]database.Skill, 0)
		err = database.DB.In("skill_id", equip.Skills).Find(&skills)
		if err != nil {
			return err
		}
		for i, skill := range skills {
			pt = freetype.Pt(270, 140+i*50+fontSize)
			_, err = c.DrawString(fmt.Sprintf("%s", skill.Description), pt)
			if err != nil {
				return err
			}
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
		dst.Close()
	}
	return nil
}
