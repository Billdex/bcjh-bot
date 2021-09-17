package messageservice

import (
	"bcjh-bot/config"
	"bcjh-bot/global"
	"bcjh-bot/model/database"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/golang/freetype"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type userChefInfo struct {
	Name          string
	Avatar        string
	Rarity        int
	Origin        string
	Stirfry       int
	Bake          int
	Boil          int
	Steam         int
	Fry           int
	Cut           int
	Meat          int
	Flour         int
	Fish          int
	Vegetable     int
	Sweet         int
	Sour          int
	Spicy         int
	Salty         int
	Bitter        int
	Tasty         int
	Skill         string
	UltimateGoals []string
	UltimateSkill string
}

var practiceTaskList1 = []func(rarity int) string{
	func(rarity int) string { return fmt.Sprintf("获得群内龙王%d次", 4*rarity) },
	func(rarity int) string { return fmt.Sprintf("获得群内财神荣誉%d次", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("成为屠龙勇士%d次", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("获得群聊之火%d天", 10*rarity) },
	func(rarity int) string { return fmt.Sprintf("成为快乐源泉%d天", 10*rarity) },
}

var practiceTaskList2 = []func(rarity int) string{
	func(rarity int) string { return fmt.Sprintf("一天内水群%d条消息", 50*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内获得%d个群的龙王", 1*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内发出%d个红包", 1*rarity) },
	func(rarity int) string { return fmt.Sprintf("一次发出总金额大于%d的红包", 5*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内发送%d个表情包", 10*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内被管理员禁言%d次", 1*rarity) },
	func(rarity int) string { return fmt.Sprintf("单次被禁言时长达到%d分钟", 10*rarity) },
	func(rarity int) string { return fmt.Sprintf("连续潜水%d天", 10*rarity) },
	func(rarity int) string {
		return fmt.Sprintf("营业结束忘记结算单次时长超过%d小时", 3*rarity)
	},
	func(rarity int) string {
		return fmt.Sprintf("结算后忘记开业单次时长超过%d小时", 3*rarity)
	},
	func(rarity int) string { return fmt.Sprintf("一天内实验室研发成功%d次", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内实验室研发失败%d次", 1*rarity) },
	func(rarity int) string { return fmt.Sprintf("连续%d天运势指数在中吉以上", 1*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内璧池抽奖%d发", 3*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内新手池抽奖%d发", 10*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内中级池抽奖%d发", 5*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内高级池抽奖%d发", 1*rarity) },
	func(rarity int) string { return fmt.Sprintf("连续%d周厨神分享保", 1*rarity) },
	func(rarity int) string { return fmt.Sprintf("连续%d周厨神高保", 1*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内完成%d个主线任务", 3*rarity) },
}

var practiceTaskList3 = []func(rarity int) string{
	func(rarity int) string { return fmt.Sprintf("累计群内活跃%d天", 10*rarity) },
	func(rarity int) string { return fmt.Sprintf("累计发送表情包%d天", 5*rarity) },
	func(rarity int) string { return fmt.Sprintf("累计发出红包%d个", 5*rarity) },
	func(rarity int) string { return fmt.Sprintf("群内累计发出红包金额%d元", 10*rarity) },
	func(rarity int) string { return fmt.Sprintf("累计潜水%d天", 5*rarity) },
	func(rarity int) string { return fmt.Sprintf("群活跃等级到达%d级", 20*rarity) },
	func(rarity int) string { return fmt.Sprintf("累计因开车被管理员禁言%d次", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("累计群内被禁言时长达到%d小时", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("累计抽到大大吉%d次", 5*rarity) },
	func(rarity int) string { return fmt.Sprintf("累计抽到小小吉%d次", 5*rarity) },
	func(rarity int) string { return fmt.Sprintf("累计忘记开业时长超过%d小时", 10*rarity) },
	func(rarity int) string {
		return fmt.Sprintf("累计因上错厨师导致厨神分享保%d次", 1*rarity)
	},
	func(rarity int) string {
		return fmt.Sprintf("累计因上错菜谱导致厨神分享保%d次", 1*rarity)
	},
	func(rarity int) string {
		return fmt.Sprintf("累计因上错厨具导致厨神分享保%d次", 1*rarity)
	},
	func(rarity int) string { return fmt.Sprintf("累计厨神御前高保%d次", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("累计集市玉璧买菜%d次", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("累计招募玉璧买菜%d次", 5*rarity) },
}

var practiceSkillList = []func(rarity int) string{
	func(rarity int) string { return fmt.Sprintf("抢红包运气王概率提升%d%%", 5*rarity) },
	func(rarity int) string { return fmt.Sprintf("游戏求签大大吉概率提升%d%%", 6*rarity) },
	func(rarity int) string { return fmt.Sprintf("获得群内龙王的概率提升%d%%", 6*rarity) },
	func(rarity int) string { return fmt.Sprintf("连任龙王被抢走的概率减少%d%%", 3*rarity) },
	func(rarity int) string { return fmt.Sprintf("周常获得玉璧奖励的概率提升%d%%", 5*rarity) },
	func(rarity int) string {
		return fmt.Sprintf("周常获得银符文奖励的概率减少%d%%", 3*rarity)
	},
	func(rarity int) string { return fmt.Sprintf("实验室研发成功概率提升%d%%", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("实验室研发炸锅概率减少%d%%", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("招募前十发出坑的概率提高%d%%", 1*rarity) },
}

func RandChefImg(c *scheduler.Context) {
	if c.GetMessageType() == onebot.MessageTypePrivate {
		return
	}

	now := time.Now()
	timeSeed := now.Unix()
	timeSeed -= int64(now.Hour() * 3600)
	timeSeed -= int64(now.Minute() * 60)
	timeSeed -= int64(now.Second())
	timeSeed -= int64(now.Weekday()-1) * 3600 * 24

	event := c.GetGroupEvent()
	skills := make([]database.Skill, 0)
	err := database.DB.Find(&skills)
	if err != nil || len(skills) == 0 {
		logger.Error("数据库查询出错!", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}
	global.RandLock.Lock()
	rand.Seed(c.GetSenderId() + timeSeed)
	rarity := rand.Intn(5) + 1
	chef := userChefInfo{
		Name:   c.GetSenderNickname(),
		Rarity: rarity,
		Avatar: fmt.Sprintf("https://q1.qlogo.cn/g?b=qq&nk=%d&s=100", event.Sender.UserId),

		Meat:      rand.Intn(rarity + 4),
		Flour:     rand.Intn(rarity + 4),
		Fish:      rand.Intn(rarity + 4),
		Vegetable: rand.Intn(rarity + 4),

		Skill: skills[rand.Intn(len(skills))].Description,
		UltimateGoals: []string{
			practiceTaskList1[rand.Intn(len(practiceTaskList1))](rarity),
			practiceTaskList2[rand.Intn(len(practiceTaskList2))](rarity),
			practiceTaskList3[rand.Intn(len(practiceTaskList3))](rarity),
		},
		UltimateSkill: practiceSkillList[rand.Intn(len(practiceSkillList))](rarity),
	}
	for chef.Stirfry < 15+60*rarity && chef.Bake < 15+60*rarity && chef.Boil < 15+60*rarity &&
		chef.Steam < 15+60*rarity && chef.Fry < 15+60*rarity && chef.Cut < 15+60*rarity {
		chef.Stirfry = rand.Intn(15*rarity*rarity - 2*rarity + 120)
		chef.Bake = rand.Intn(15*rarity*rarity - 2*rarity + 120)
		chef.Boil = rand.Intn(15*rarity*rarity - 2*rarity + 120)
		chef.Steam = rand.Intn(15*rarity*rarity - 2*rarity + 120)
		chef.Fry = rand.Intn(15*rarity*rarity - 2*rarity + 120)
		chef.Cut = rand.Intn(15*rarity*rarity - 2*rarity + 120)
	}
	condiment := 30*(rarity-1) + (rarity+1)/4*10 + rand.Intn(10+10*(rarity/2))
	switch rand.Intn(6) {
	case 0:
		chef.Sweet = condiment
	case 1:
		chef.Sour = condiment
	case 2:
		chef.Spicy = condiment
	case 3:
		chef.Salty = condiment
	case 4:
		chef.Bitter = condiment
	case 5:
		chef.Tasty = condiment
	default:
		chef.Sweet = condiment
	}
	global.RandLock.Unlock()

	// 数值调整
	groupInfo, err := c.GetBot().GetGroupInfo(c.GetGroupEvent().GroupId)
	if err != nil {
		logger.Error("获取群信息失败", err)
		_, _ = c.Reply("生成图鉴失败")
		return
	}
	chef.Origin = groupInfo.GroupName
	nameRune := []rune(c.GetSenderNickname())
	if len(nameRune) > 10 {
		chef.Name = string(nameRune[:10])
	}
	if chef.Stirfry < rarity*18 {
		chef.Stirfry = 0
	}
	if chef.Bake < rarity*18 {
		chef.Bake = 0
	}
	if chef.Boil < rarity*18 {
		chef.Boil = 0
	}
	if chef.Steam < rarity*18 {
		chef.Steam = 0
	}
	if chef.Fry < rarity*18 {
		chef.Fry = 0
	}
	if chef.Cut < rarity*18 {
		chef.Cut = 0
	}

	if chef.Meat < 2+rarity/2 {
		chef.Meat = 0
	}
	if chef.Flour < 2+rarity/2 {
		chef.Flour = 0
	}
	if chef.Fish < 2+rarity/2 {
		chef.Fish = 0
	}
	if chef.Vegetable < 2+rarity/2 {
		chef.Vegetable = 0
	}

	bytesImg, err := createRandChefImg(chef)
	if err != nil {
		logger.Error("绘制图片出错", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}

	base64Img := base64.StdEncoding.EncodeToString(bytesImg)
	cqImg := onebot.GetCQImage(base64Img, "base64")
	_, _ = c.Reply(cqImg)
}

func createRandChefImg(chef userChefInfo) ([]byte, error) {
	dx := 800       // 图鉴背景图片的宽度
	dy := 800       // 图鉴背景图片的高度
	titleSize := 50 // 标题字体尺寸
	fontSize := 36  // 内容字体尺寸
	fontDPI := 72.0 // dpi

	// 获取字体文件
	resourceFontDir := config.AppConfig.Resource.Font
	fontFile := resourceFontDir + "/yuan500W.ttf"
	//读字体数据
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	fontColor := color.RGBA{A: 255}

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
	bgFile, err := os.Open(fmt.Sprintf("%s/chef/chef_%s.png", config.AppConfig.Resource.Image, condimentType))
	if err != nil {
		return nil, err
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

	// 输出用户头像
	r, err := http.Get(chef.Avatar)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	avatarImg, err := jpeg.Decode(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	avatarImg = resize.Resize(200, 200, avatarImg, resize.Bilinear)
	draw.Draw(img,
		image.Rect(50, 118, 50+200, 118+200),
		avatarImg,
		image.Point{X: 0, Y: 0},
		draw.Over)

	// 输出名字
	pt := freetype.Pt(50, 22+titleSize)
	_, err = c.DrawString(fmt.Sprintf("%s", chef.Name), pt)
	if err != nil {
		return nil, err
	}

	// 输出稀有度
	rarityFile, err := os.Open(fmt.Sprintf("%s/chef/rarity_%d.png", config.AppConfig.Resource.Image, chef.Rarity))
	if err != nil {
		return nil, err
	}
	defer rarityFile.Close()
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
		return nil, err
	}
	pt = freetype.Pt(536, 104+fontSize)
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Bake), pt)
	if err != nil {
		return nil, err
	}
	pt = freetype.Pt(705, 104+fontSize)
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Boil), pt)
	if err != nil {
		return nil, err
	}
	pt = freetype.Pt(365, 164+fontSize)
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Steam), pt)
	if err != nil {
		return nil, err
	}
	pt = freetype.Pt(536, 164+fontSize)
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Fry), pt)
	if err != nil {
		return nil, err
	}
	pt = freetype.Pt(705, 164+fontSize)
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Cut), pt)
	if err != nil {
		return nil, err
	}

	// 输出采集数据
	pt = freetype.Pt(365, 230+fontSize)
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Meat), pt)
	if err != nil {
		return nil, err
	}
	pt = freetype.Pt(536, 230+fontSize)
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Flour), pt)
	if err != nil {
		return nil, err
	}
	pt = freetype.Pt(365, 290+fontSize)
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Vegetable), pt)
	if err != nil {
		return nil, err
	}
	pt = freetype.Pt(536, 290+fontSize)
	_, err = c.DrawString(fmt.Sprintf("%d", chef.Fish), pt)
	if err != nil {
		return nil, err
	}

	// 输出调料数据
	pt = freetype.Pt(705, 290+fontSize)
	_, err = c.DrawString(fmt.Sprintf("%d", condiment), pt)
	if err != nil {
		return nil, err
	}

	// 输出来源数据
	pt = freetype.Pt(150, 365+fontSize)
	_, err = c.DrawString(fmt.Sprintf("%s", chef.Origin), pt)
	if err != nil {
		return nil, err
	}

	// 输出技法数据
	pt = freetype.Pt(150, 435+fontSize)
	_, err = c.DrawString(fmt.Sprintf("%s", chef.Skill), pt)
	if err != nil {
		return nil, err
	}

	// 输出修炼效果数据
	pt = freetype.Pt(150, 505+fontSize)
	_, err = c.DrawString(fmt.Sprintf("%s", chef.UltimateSkill), pt)
	if err != nil {
		return nil, err
	}

	// 输出修炼任务数据
	for i, task := range chef.UltimateGoals {
		pt = freetype.Pt(120, 625+i*50+fontSize)
		_, err = c.DrawString(fmt.Sprintf("%s", task), pt)
		if err != nil {
			return nil, err
		}

	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
