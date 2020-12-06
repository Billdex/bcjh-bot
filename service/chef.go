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
)

func ChefQuery(c *onebot.Context, args []string) {
	logger.Info("厨师查询，参数:", args)

	if len(args) == 0 {
		err := bot.SendMessage(c, chefHelp())
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

	chefs := make([]database.Chef, 0)
	err := database.DB.Where("gallery_id = ?", args[0]).Asc("gallery_id").Find(&chefs)
	if err != nil {
		logger.Error("查询数据库出错!", err)
		_ = bot.SendMessage(c, "查询数据失败!")
		return
	}
	if len(chefs) == 0 {
		err = database.DB.Where("name like ?", "%"+args[0]+"%").Asc("gallery_id").Find(&chefs)
		if err != nil {
			logger.Error("查询数据库出错!", err)
			_ = bot.SendMessage(c, "查询数据失败!")
			return
		}
	}

	var msg string
	if len(chefs) == 0 {
		msg = "哎呀，好像找不到呢!"
	} else if len(chefs) == 1 {
		chef := chefs[0]
		// 尝试寻找图片文件，未找到则按照文字格式发送
		resourceImageDir := config.AppConfig.Resource.Image + "/chef"
		imagePath := fmt.Sprintf("%s/chef_%s.png", resourceImageDir, chef.GalleryId)
		logger.Debug("imagePath:", imagePath)
		if has, err := util.PathExists(imagePath); has {
			msg = bot.GetCQImage(imagePath, "file")
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
			_, err = database.DB.Where("skill_id = ?", chef.SkillId).Get(skill)
			if err != nil {
				logger.Error("查询数据库出错!", err)
				_ = bot.SendMessage(c, util.SystemErrorNote)
				return
			}
			ultimateSkill := new(database.Skill)
			_, err = database.DB.Where("skill_id = ?", chef.UltimateSkill).Get(ultimateSkill)
			if err != nil {
				logger.Error("查询数据库出错!", err)
				_ = bot.SendMessage(c, util.SystemErrorNote)
				return
			}
			ultimateGoals := make([]database.Quest, 0)
			err = database.DB.In("quest_id", chef.UltimateGoals).Find(&ultimateGoals)
			if err != nil {
				logger.Error("查询数据库出错!", err)
				_ = bot.SendMessage(c, util.SystemErrorNote)
				return
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
			//msg += fmt.Sprintf("修炼任务:%s", goals)
		}

	} else {
		msg = "查询到以下厨师:\n"
		for p, chef := range chefs {
			msg += fmt.Sprintf("%s %s", chef.GalleryId, chef.Name)
			if p != len(chefs)-1 {
				msg += "\n"
				if p == util.MaxQueryListLength-1 {
					msg += "......"
					break
				}
			}
		}
	}

	err = bot.SendMessage(c, msg)
	if err != nil {
		logger.Error("发送信息失败!", err)
	}
}

func ChefInfoToImage(chefs []database.Chef) error {
	dx := 800           // 图鉴背景图片的宽度
	dy := 800           // 图鉴背景图片的高度
	avatarWidth := 200  // 厨师头像图片宽度
	avatarHeight := 200 // 厨师头像图片高度
	titleSize := 52     // 标题字体尺寸
	fontSize := 36      // 内容字体尺寸
	fontDPI := 72.0     // dpi

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
	fontColor := color.RGBA{0, 0, 0, 255}
	// 获取头像图鉴总图
	resourceImgDir := config.AppConfig.Resource.Image
	chefImgPath := resourceImgDir + "/chef"
	galleryFilePath := chefImgPath + "/chef_gallery.png"
	galleryFile, err := os.Open(galleryFilePath)
	if err != nil {
		return err
	}
	defer galleryFile.Close()
	galleryImg, err := png.Decode(galleryFile)
	if err != nil {
		return err
	}
	for p, chef := range chefs {
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
		defer bgFile.Close()
		img := image.NewRGBA(image.Rect(0, 0, dx, dy))
		bg, _ := png.Decode(bgFile)

		draw.Draw(img, img.Bounds(), bg, bg.Bounds().Min, draw.Src)

		c := freetype.NewContext()
		c.SetDPI(fontDPI)
		c.SetFont(font)
		c.SetClip(img.Bounds())
		c.SetDst(img)
		c.SetSrc(image.NewUniform(fontColor))
		c.SetFontSize(float64(titleSize))

		// 输出厨师头像
		avatarStartX := (p % 18) * avatarWidth
		avatarStartY := (p / 18) * avatarHeight
		draw.Draw(img,
			image.Rect(50, 118, 50+avatarWidth, 118+avatarHeight),
			galleryImg,
			image.Point{avatarStartX, avatarStartY},
			draw.Over)

		// 输出图鉴ID与厨师名
		pt := freetype.Pt(35, 20+titleSize)
		_, err = c.DrawString(fmt.Sprintf("%s %s", chef.GalleryId, chef.Name), pt)
		if err != nil {
			return err
		}

		// 输出性别
		genderFile, err := os.Open(fmt.Sprintf("%s/gender_%d.png", chefImgPath, chef.Gender))
		if err != nil {
			return err
		}
		defer genderFile.Close()
		genderImg, _ := png.Decode(genderFile)
		draw.Draw(img,
			image.Rect(480, 30, 480+44, 30+44),
			genderImg,
			image.ZP,
			draw.Over)

		// 输出稀有度
		rarityFile, err := os.Open(fmt.Sprintf("%s/rarity_%d.png", chefImgPath, chef.Rarity))
		if err != nil {
			return err
		}
		defer genderFile.Close()
		rarityImg, _ := png.Decode(rarityFile)
		draw.Draw(img,
			image.Rect(545, 30, 545+240, 30+44),
			rarityImg,
			image.ZP,
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
		pt = freetype.Pt(149, 360+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%s", chef.Origin), pt)
		if err != nil {
			return err
		}

		// 输出技法数据
		skill := new(database.Skill)
		_, err = database.DB.Where("skill_id = ?", chef.SkillId).Get(skill)
		if err != nil {
			return err
		}
		pt = freetype.Pt(149, 430+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%s", skill.Description), pt)
		if err != nil {
			return err
		}

		// 输出修炼效果数据
		ultimateSkill := new(database.Skill)
		_, err = database.DB.Where("skill_id = ?", chef.UltimateSkill).Get(ultimateSkill)
		if err != nil {
			return err
		}
		pt = freetype.Pt(149, 500+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%s", ultimateSkill.Description), pt)
		if err != nil {
			return err
		}

		// 输出修炼任务数据
		ultimateGoals := make([]database.Quest, 0)
		err = database.DB.In("quest_id", chef.UltimateGoals).Find(&ultimateGoals)
		if err != nil {
			return err
		}
		for p, goal := range ultimateGoals {
			pt = freetype.Pt(120, 620+p*50+fontSize)
			_, err = c.DrawString(fmt.Sprintf("%s", goal.Goal), pt)
			if err != nil {
				return err
			}
		}

		// 以PNG格式保存文件
		dst, err := os.Create(fmt.Sprintf("%s/chef_%s.png", chefImgPath, chef.GalleryId))
		if err != nil {
			return err
		}
		defer dst.Close()

		err = png.Encode(dst, img)
		if err != nil {
			return err
		}
	}
	return nil
}
