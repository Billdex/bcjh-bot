package messageservice

import (
	"bcjh-bot/dao"
	"bcjh-bot/model/database"
	"bcjh-bot/model/userdata"
	"bcjh-bot/scheduler"
	"bcjh-bot/util/logger"
	"encoding/json"
	"fmt"
	"math"
	"modernc.org/mathutil"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 指定厨师修炼结果
type userUltimateResult struct {
	HasBoundID bool           // 是否绑定「白菜菊花 ID」
	Chef       database.Chef  // 指定厨师
	ChefGot    bool           // 是否已有（未绑定「白菜菊花 ID」则默认未 false）
	Equip      database.Equip // 指定厨具
	Recipes    []resultRecipe

	UtlimateQuests   []database.Quest // 厨师多个修炼
	UtlimateMessages []string         // 厨师多个修炼任务对应的回复文本
	Page             int
}

func (ur userUltimateResult) String() string {
	gotchef := map[bool]string{true: "[已有]", false: map[bool]string{true: "[未拥有]", false: "[公开]"}[ur.HasBoundID]}
	gotreci := map[bool]string{true: "✅ ", false: ""}
	ranks := []string{"难", "可", "优", "特", "神", "传"}
	pagesize := 8

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s【%s】%s", gotchef[ur.ChefGot], ur.Chef.Name, strings.Repeat("🔥", ur.Chef.Rarity)))
	if ur.Equip.EquipId != 0 {
		sb.WriteString("\n" + ur.Equip.Name + "：" + strings.Join(ur.Equip.SkillDescs, "，"))
	}
	sb.WriteString("\n技法：" + cookstr([]int{ur.Chef.Stirfry, ur.Chef.Boil, ur.Chef.Cut, ur.Chef.Fry, ur.Chef.Bake, ur.Chef.Steam}))
	sb.WriteString("\n修炼技能：" + ur.Chef.UltimateSkillDesc)
	sb.WriteString("\n修炼任务：")
	for i, quest := range ur.UtlimateQuests {
		sb.WriteString(fmt.Sprintf("\n[%d] %s", i+1, quest.Goal))
		if len(quest.Conditions) > 0 {
			c := quest.Conditions[0]
			if !(len(c.Skill) > 0 && c.Rarity > 0 && c.Rank > 0 && c.Num > 0) {
				continue
			}
			sb.WriteString("\n菜谱推荐（未拥有显示来源）：")
			for ii := (ur.Page - 1) * pagesize; ii < len(ur.Recipes) && ii < ur.Page*pagesize; ii++ {
				rr := ur.Recipes[ii]
				// 已有显示时间，未拥有显示池子
				sb.WriteString(fmt.Sprintf("\n[%s] %s%s", ranks[rr.Rank], gotreci[rr.RecipeGot], rr.Recipe.Name))
				t := (time.Duration(rr.Recipe.TotalTime) * time.Second).String()

				if rr.CanDo {
					sb.WriteString(fmt.Sprintf(" %s", t))
				} else {
					sb.WriteString(" " + cookstr([]int{rr.Recipe.Stirfry, rr.Recipe.Boil, rr.Recipe.Cut, rr.Recipe.Fry, rr.Recipe.Bake, rr.Recipe.Steam}))
				}
				if !rr.RecipeGot {
					sb.WriteString(fmt.Sprintf("（%s）", rr.Recipe.Origin))
				}
			}
			if ur.Page*10 < len(ur.Recipes) {
				sb.WriteString("\n......")
			}
		}
	}
	sb.WriteString(fmt.Sprintf("\n每页 %d 条，共 %d 条（p%d/p%d）", pagesize, len(ur.Recipes), ur.Page, len(ur.Recipes)/pagesize+1-map[bool]int{true: 1, false: 0}[len(ur.Recipes)%pagesize == 0]))
	return sb.String()
}

type resultRecipe struct {
	Recipe    database.Recipe // 若未达到修炼条件（CanDo == false），其中的技法数值会修改为神差值
	RecipeGot bool            // 是否已有
	CanDo     bool            // 是否满足条件
	Rank      int             // 做到的等级
	RankDiff  int             // 目标等级差值
}

func getUserData(userId int64) (userdata.UserData, bool) {
	var userData userdata.UserData
	ud, err := dao.FindUserDataWithUserId(userId)
	if err != nil {
		//_, _ = c.Reply("用户未导入数据，请使用「导入数据 <白菜菊花个人数据ID>」来导入数据，导入后长期有效")
		return userData, false
	}
	if ud.Data == "" {
		logger.Warnf("用户数据为空")
		return userData, false
	}
	if err = json.Unmarshal([]byte(ud.Data), &userData); err != nil {
		logger.Errorf("读取用户数据错误：%s", err)
		return userData, false
	}
	return userData, true
}

func UltimateQuery(c *scheduler.Context) {
	// 默认参数：厨师名，厨具关键词（可以是厨具 id），页码，已有
	chefName, equipName, page, onlyHave := "", "", 1, false

	// 处理参数
	args := strings.Split(c.PretreatedMessage, " ")
	for _, arg := range args {
		if arg == "已有" {
			onlyHave = true
		} else if r := regexp.MustCompile("[pP]-?([0-9]+)"); r.MatchString(arg) {
			match := r.FindAllStringSubmatch(arg, -1)
			if len(match) < 1 || len(match[0]) < 2 {
				continue
			}
			if p, err := strconv.Atoi(match[0][1]); err == nil && p > 1 {
				page = p
			}
		} else if chefName == "" {
			chefName = arg
		} else if equipName == "" {
			equipName = arg
		}
	}
	if chefName == "" {
		_, _ = c.Reply("参数不足")
		return
	}

	ultResult := userUltimateResult{Page: page}

	var userData userdata.UserData
	userData, ultResult.HasBoundID = getUserData(c.GetSenderId()) // 获取用户个人数据

	if onlyHave && !ultResult.HasBoundID {
		_, _ = c.Reply(fmt.Sprintf("查看已有菜谱信息需要先绑定白菜菊花ID，请使用「%s个人数据导入 <ID>」绑定", prefixCharacters[0]))
		return
	}

	// 获取数据库中所有厨师数据
	allChefs, err := dao.FindAllChefs()
	if err != nil {
		logger.Errorw("获取所有厨师信息失败", "err", err.Error())
		_, _ = c.Reply("获取所有厨师信息失败")
		return
	}

	// 查找用户指定的厨师
	chefs, note := filterChefsByName(allChefs, chefName)
	if note != "" {
		logger.Error("厨师查询失败:", note)
		_, _ = c.Reply(note)
		return
	}

	if len(chefs) == 0 {
		_, _ = c.Reply("没有找到名为 " + chefName + " 的厨师")
		return
	} else if len(chefs) > 15 {
		_, _ = c.Reply(fmt.Sprintf("关键词 [%s] 一共匹配到了 %d 个厨师，请具体一点", chefName, len(chefs)))
		return
	} else if len(chefs) > 1 {
		chefsName := make([]string, 0, len(chefs))
		for _, chef := range chefs {
			chefsName = append(chefsName, chef.Name)
		}
		_, _ = c.Reply("识别到多位厨师：" + strings.Join(chefsName, " "))
		return
	}

	ultResult.Chef = chefs[0] // 确定厨师

	if equipName != "" {
		equips, _ := dao.SearchEquipsWithName(equipName)
		if equips != nil && len(equips) > 0 {
			if len(equips) > 15 {
				_, _ = c.Reply(fmt.Sprintf("关键词 [%s] 一共匹配到了 %d 个厨具，请具体一点", equipName, len(equips)))
				return
			} else if len(equips) > 1 {
				equipsName := make([]string, 0, len(equips))
				for _, chef := range equips {
					equipsName = append(equipsName, chef.Name)
				}
				_, _ = c.Reply("识别到多个厨具：" + strings.Join(equipsName, " "))
				return
			}
			ultResult.Equip = equips[0] // 确定厨具
		}
	}

	if ultResult.HasBoundID {
		var gotMap = make(map[int]bool, len(allChefs))
		bs, _ := userData.ChefGot.MarshalJSON()
		if err = json.Unmarshal(bs, &gotMap); err != nil {
			logger.Errorf("解析用户获取厨师数据失败：%s", err)
			_, _ = c.Reply("用户厨师数据异常")
			return
		}
		ultResult.ChefGot = gotMap[ultResult.Chef.ChefId] // 确定厨师是否已拥有
	}

	// 解析菜谱
	// 如果绑定 ID 且有“已有”参数，则使用已有菜谱
	recipes, err := dao.FindAllRecipes()
	if err != nil {
		logger.Error("获取所有菜谱失败", err)
		_, _ = c.Reply("糟糕，菜…菜谱消失了！")
		return
	}
	gotRecipesMap := make(map[int]bool, len(recipes))
	if ultResult.HasBoundID {
		bs, _ := userData.RepGot.MarshalJSON()
		if err = json.Unmarshal(bs, &gotRecipesMap); err != nil {
			logger.Errorf("解析用户获取菜谱数据失败：%s", err)
			_, _ = c.Reply("用户菜谱数据异常")
			return
		}
		if onlyHave {
			for i := 0; i < len(recipes); i++ {
				if got, ok := gotRecipesMap[recipes[i].RecipeId]; !ok || !got {
					recipes = append(recipes[:i], recipes[i+1:]...)
					i--
				}
			}
		}
	}

	if len(recipes) == 0 {
		_, _ = c.Reply("找不到菜谱了，呜呜呜~")
		return
	}

	// 给厨师叠 buff
	addBuff(&ultResult.Chef, ultResult.Equip, userData)

	// 获取厨师修炼任务
	ultResult.UtlimateQuests, err = dao.FindQuestsWithIds(ultResult.Chef.UltimateGoals)
	if err != nil {
		logger.Error("获取厨师修炼任务失败", err)
		_, _ = c.Reply("获取厨师修炼任务失败了")
		return
	}
	if len(ultResult.UtlimateQuests) == 0 {
		_, _ = c.Reply(fmt.Sprintf("天呐！莫非 %s 就是传说中无需修炼的绝世奇才！", ultResult.Chef.Name))
		return
	}

	// 分析修炼任务的条件，遍历查找满足条件的菜谱
	for _, quest := range ultResult.UtlimateQuests {

		// 分析任务条件，不同条件不同处理方案
		if len(quest.Conditions) == 0 {
			// 无条件：收集符文
			continue
		}

		cond := quest.Conditions[0]
		if len(cond.Skill)*cond.Rarity*cond.Rank*cond.Num != 0 {
			// 都不为零：技法挑战
			_recipes := recipes
			_recipes, _ = filterRecipesByRarity(_recipes, cond.Rarity, false)
			_recipes, _ = filterRecipesBySkill(_recipes, cond.Skill)

			ultResult.Recipes = make([]resultRecipe, 0, len(_recipes))
			for _, recipe := range _recipes {
				rank := chefDoLevel(ultResult.Chef, recipe)
				got, ok := gotRecipesMap[recipe.RecipeId]
				if !ok {
					got = false
				}
				rr := resultRecipe{Recipe: recipe, RecipeGot: got, CanDo: rank >= cond.Rank, Rank: rank, RankDiff: 0}
				if !rr.CanDo {
					rr.Recipe, rr.RankDiff = chefDoDiff(ultResult.Chef, recipe, cond.Rank)
				}
				ultResult.Recipes = append(ultResult.Recipes, rr)
			}
			ultResult.Recipes = sortRecipe(ultResult.Recipes)
			continue
		}
	}
	_, _ = c.Reply(ultResult.String())
}

// addBuff 计算添加后的 buff
func addBuff(chef *database.Chef, equip database.Equip, data userdata.UserData) {
	male := data.UserUltimate.Male
	if chef.Gender == 0 || chef.Gender == 1 {
		// 男
		chef.Stirfry += int(male)
		chef.Boil += int(male)
		chef.Cut += int(male)
		chef.Fry += int(male)
		chef.Bake += int(male)
		chef.Steam += int(male)
	}
	female := data.UserUltimate.Female
	if chef.Gender == 0 || chef.Gender == 2 {
		// 女
		chef.Stirfry += int(female)
		chef.Boil += int(female)
		chef.Cut += int(female)
		chef.Fry += int(female)
		chef.Bake += int(female)
		chef.Steam += int(female)
	}
	all := data.UserUltimate.All
	chef.Stirfry += int(data.UserUltimate.Stirfry) + int(all)
	chef.Boil += int(data.UserUltimate.Boil) + int(all)
	chef.Cut += int(data.UserUltimate.Knife) + int(all)
	chef.Fry += int(data.UserUltimate.Fry) + int(all)
	chef.Bake += int(data.UserUltimate.Bake) + int(all)
	chef.Steam += int(data.UserUltimate.Steam) + int(all)
	// 计算厨具的 buff
	if equip.EquipId != 0 {
		skillsMap, _ := dao.GetSkillsMap()
		if skillsMap != nil {
			for _, skillId := range equip.Skills {
				skill, ok := skillsMap[skillId]
				if !ok {
					continue
				}
				for _, effect := range skill.Effects {
					// 检查生效条件
					if effect.Tag != 0 && chef.Gender != effect.Tag {
						// 性别不一致
						logger.Warnf("性别不一致")
						continue
					}
					switch effect.Condition {
					case "Partial": // 场上所有厨师
					case "Global": // 全体厨师
					case "Self": // 自身
					default:
						continue
					}
					var adder func(old int) int
					switch effect.Calculation {
					case "Abs":
						adder = func(old int) int { return old + int(math.Ceil(effect.Value)) }
					case "Percent":
						adder = func(old int) int { return old + int(math.Ceil(effect.Value*float64(old)/100)) }
					default:
						continue
					}
					switch effect.Type {
					case "Stirfry":
						chef.Stirfry = adder(chef.Stirfry)
					case "Knife":
						chef.Cut = adder(chef.Cut)
					case "Bake":
						chef.Bake = adder(chef.Bake)
					case "Fry":
						chef.Fry = adder(chef.Fry)
					case "Boil":
						chef.Boil = adder(chef.Boil)
					case "Steam":
						chef.Steam = adder(chef.Steam)
					}
				}
			}
		}
	}
}

// chefDoLevel 厨师做这道菜的等级是多少
// 0: 做不了
// 1: 可
// 2: 优
// 3: 特
// 4: 神
// 5: 传
func chefDoLevel(chef database.Chef, recipe database.Recipe) int {
	ranks := make([]int, 0, 6)
	if recipe.Stirfry > 0 {
		ranks = append(ranks, chef.Stirfry/recipe.Stirfry)
	}
	if recipe.Boil > 0 {
		ranks = append(ranks, chef.Boil/recipe.Boil)
	}
	if recipe.Cut > 0 {
		ranks = append(ranks, chef.Cut/recipe.Cut)
	}
	if recipe.Fry > 0 {
		ranks = append(ranks, chef.Fry/recipe.Fry)
	}
	if recipe.Bake > 0 {
		ranks = append(ranks, chef.Bake/recipe.Bake)
	}
	if recipe.Steam > 0 {
		ranks = append(ranks, chef.Steam/recipe.Steam)
	}
	return mathutil.MinVal(5, ranks...)
}

// chefDoDiff 计算厨师做菜的神查值
func chefDoDiff(chef database.Chef, recipe database.Recipe, rank int) (database.Recipe, int) {
	recipe.Stirfry = ((chef.Stirfry - rank*recipe.Stirfry) >> (strconv.IntSize - 1)) & 1 * (rank*recipe.Stirfry - chef.Stirfry)
	recipe.Boil = ((chef.Boil - rank*recipe.Boil) >> (strconv.IntSize - 1)) & 1 * (rank*recipe.Boil - chef.Boil)
	recipe.Cut = ((chef.Cut - rank*recipe.Cut) >> (strconv.IntSize - 1)) & 1 * (rank*recipe.Cut - chef.Cut)
	recipe.Fry = ((chef.Fry - rank*recipe.Fry) >> (strconv.IntSize - 1)) & 1 * (rank*recipe.Fry - chef.Fry)
	recipe.Bake = ((chef.Bake - rank*recipe.Bake) >> (strconv.IntSize - 1)) & 1 * (rank*recipe.Bake - chef.Bake)
	recipe.Steam = ((chef.Steam - rank*recipe.Steam) >> (strconv.IntSize - 1)) & 1 * (rank*recipe.Steam - chef.Steam)
	return recipe, recipe.Stirfry + recipe.Boil + recipe.Cut + recipe.Fry + recipe.Bake + recipe.Steam
}

// sortRecipe 菜谱结果排序
func sortRecipe(recipes []resultRecipe) []resultRecipe {
	// 单时间
	for i := 0; i < len(recipes); i++ {
		for j := 0; j < len(recipes)-i-1; j++ {
			if recipes[j].Recipe.Time > recipes[j+1].Recipe.Time {
				recipes[j], recipes[j+1] = recipes[j+1], recipes[j]
			}
		}
	}
	// 是否可做
	for i := 0; i < len(recipes); i++ {
		for j := 0; j < len(recipes)-i-1; j++ {
			if recipes[j].CanDo == false && recipes[j+1].CanDo == true {
				recipes[j], recipes[j+1] = recipes[j+1], recipes[j]
			}
		}
	}
	// 神差值
	for i := 0; i < len(recipes); i++ {
		for j := 0; j < len(recipes)-i-1; j++ {
			if recipes[j].RankDiff > recipes[j+1].RankDiff {
				recipes[j], recipes[j+1] = recipes[j+1], recipes[j]
			}
		}
	}
	return recipes
}

func cookstr(points []int) string {
	cooks := []string{"炒", "煮", "切", "炸", "烤", "蒸"}
	items := make([]string, 0, len(cooks))
	for i := 0; i < len(points); i++ {
		maxIndex := 0
		for j := 1; j < len(points); j++ {
			if points[j] > points[maxIndex] {
				maxIndex = j
			}
		}
		if points[maxIndex] > 0 {
			items = append(items, fmt.Sprintf("%s:%d", cooks[maxIndex], points[maxIndex]))
		}
		cooks = append(cooks[:maxIndex], cooks[maxIndex+1:]...)
		points = append(points[:maxIndex], points[maxIndex+1:]...)
		i--
	}
	return strings.Join(items, " ")
}
