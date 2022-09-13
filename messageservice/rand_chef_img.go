package messageservice

import (
	"bcjh-bot/config"
	"bcjh-bot/dao"
	"bcjh-bot/model/database"
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math/rand"
	"net/http"
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
	func(rarity int) string { return fmt.Sprintf("一天内获得%d个群的龙王", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内发出%d个红包", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("一次发出总金额大于%d的红包", 5*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内发送%d个表情包", 10*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内被管理员禁言%d次", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内群内开车%d次", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("单次被禁言时长达到%d分钟", 10*rarity) },
	func(rarity int) string { return fmt.Sprintf("连续潜水%d天", 10*rarity) },
	func(rarity int) string {
		return fmt.Sprintf("营业结束忘记结算单次时长超过%d小时", 3*rarity)
	},
	func(rarity int) string {
		return fmt.Sprintf("结算后忘记开业单次时长超过%d小时", 3*rarity)
	},
	func(rarity int) string { return fmt.Sprintf("一天内实验室研发成功%d次", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内实验室研发失败%d次", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("连续%d天运势指数在中吉以上", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内璧池抽奖%d发", 3*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内新手池抽奖%d发", 10*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内中级池抽奖%d发", 5*rarity) },
	func(rarity int) string { return fmt.Sprintf("一天内高级池抽奖%d发", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("连续%d周厨神分享保", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("连续%d周厨神高保", 2*rarity) },
	func(rarity int) string { return fmt.Sprintf("连续%d天开车未被管理员发现", 2*rarity) },
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
	func(rarity int) string { return fmt.Sprintf("累计因调戏管理员被禁言%d次", 2*rarity) },
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
	func(rarity int) string { return fmt.Sprintf("开车被管理员抓住的概率减少%d%%", 2*rarity) },
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
	weekday := now.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	timeSeed -= int64(weekday-1) * 3600 * 24
	selfRand := rand.New(rand.NewSource(c.GetSenderId() + timeSeed))

	event := c.GetGroupEvent()
	skills, err := dao.FindAllSkills()
	if err != nil || len(skills) == 0 {
		logger.Error("数据库查询出错!", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}
	rarity := selfRand.Intn(5) + 1
	avatar, err := DownloadImage(fmt.Sprintf("https://q1.qlogo.cn/g?b=qq&nk=%d&s=100", event.Sender.UserId))
	if err != nil {
		logger.Error("下载图片出错!", err)
		_, _ = c.Reply("获取用户头像失败啦")
		return
	}
	avatar = resize.Resize(200, 200, avatar, resize.Bilinear)
	groupInfo, err := c.GetBot().GetGroupInfo(c.GetGroupEvent().GroupId)
	if err != nil {
		logger.Error("获取群信息失败", err)
		_, _ = c.Reply("生成图鉴失败")
		return
	}
	var name = c.GetSenderNickname()
	if len([]rune(name)) > 10 {
		name = string([]rune(name)[:10])
	}
	chef := database.ChefData{
		Chef: database.Chef{
			Name:      name,
			Rarity:    rarity,
			Origin:    groupInfo.GroupName,
			Meat:      selfRand.Intn(rarity + 4),
			Flour:     selfRand.Intn(rarity + 4),
			Fish:      selfRand.Intn(rarity + 4),
			Vegetable: selfRand.Intn(rarity + 4),
		},
		Avatar: avatar,
		Skill:  skills[selfRand.Intn(len(skills))].Description,
		UltimateGoals: []string{
			practiceTaskList1[selfRand.Intn(len(practiceTaskList1))](rarity),
			practiceTaskList2[selfRand.Intn(len(practiceTaskList2))](rarity),
			practiceTaskList3[selfRand.Intn(len(practiceTaskList3))](rarity),
		},
		UltimateSkill: practiceSkillList[selfRand.Intn(len(practiceSkillList))](rarity),
	}

	for chef.Stirfry < 15+60*rarity && chef.Bake < 15+60*rarity && chef.Boil < 15+60*rarity &&
		chef.Steam < 15+60*rarity && chef.Fry < 15+60*rarity && chef.Cut < 15+60*rarity {
		chef.Stirfry = selfRand.Intn(15*rarity*rarity - 2*rarity + 120)
		chef.Bake = selfRand.Intn(15*rarity*rarity - 2*rarity + 120)
		chef.Boil = selfRand.Intn(15*rarity*rarity - 2*rarity + 120)
		chef.Steam = selfRand.Intn(15*rarity*rarity - 2*rarity + 120)
		chef.Fry = selfRand.Intn(15*rarity*rarity - 2*rarity + 120)
		chef.Cut = selfRand.Intn(15*rarity*rarity - 2*rarity + 120)
	}
	condiment := 20*rarity + (rarity-1)/2*10 + rarity/5*10 + selfRand.Intn(10+10*(rarity/2))
	switch selfRand.Intn(6) {
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

	// 数值调整
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

	// 载入一些资源文件
	// 载入字体文件
	font, err := util.LoadFontFile(fmt.Sprintf("%s/%s", config.AppConfig.Resource.Font, "yuan500W.ttf"))
	if err != nil {
		logger.Error("加载字体失败", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}
	// 载入背景图片
	bgImg, err := util.LoadPngImageFile(fmt.Sprintf("%s/chef/chef_%s.png", config.AppConfig.Resource.Image, chef.GetCondimentType()))
	if err != nil {
		logger.Error("加载背景图片失败", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}
	// 载入稀有度图片
	rarityImg, err := util.LoadPngImageFile(fmt.Sprintf("%s/chef/rarity_%d.png", config.AppConfig.Resource.Image, chef.Rarity))
	if err != nil {
		logger.Error("加载稀有度图片失败", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}
	// 绘制结果
	img, err := GenerateChefImage(chef, font, bgImg, nil, rarityImg)
	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		logger.Error("图片转buffer失败", err)
		_, _ = c.Reply(e.SystemErrorNote)
		return
	}

	base64Img := base64.StdEncoding.EncodeToString(buf.Bytes())
	cqImg := onebot.GetCQImage(base64Img, "base64")
	_, _ = c.Reply(cqImg)
}

// DownloadImage 下载图片并导出image.Image 对象
func DownloadImage(url string) (image.Image, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// 导出image实例
	img, err := jpeg.Decode(bytes.NewReader(body))
	if err != nil {
		// 读取失败则再尝试用png读取
		img, err = png.Decode(bytes.NewReader(body))
		if err != nil {
			// 再失败就试试gif
			img, err = gif.Decode(bytes.NewReader(body))
			if err != nil {
				return nil, err
			}
		}
	}
	return img, nil
}
