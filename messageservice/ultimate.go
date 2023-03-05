package messageservice

import (
	"bcjh-bot/dao"
	"bcjh-bot/model/database"
	"bcjh-bot/model/userdata"
	"bcjh-bot/scheduler"
	"bcjh-bot/util/logger"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Ultimate(c *scheduler.Context) {

	var response = struct {
		Chef     database.Chef
		ChefGot  bool
		CanDo    []database.Recipe
		CanNotDo []database.Recipe
	}{}

	args := strings.Split(c.PretreatedMessage, " ")
	if len(args) < 2 {
		_, _ = c.Reply("参数不足")
		return
	}

	bcjhId, err := strconv.Atoi(args[0])
	if err != nil {
		_, _ = c.Reply("白菜菊花 id 错误")
	}

	chefName := args[1]

	userData, err := userdata.LoadUserData(bcjhId)
	if err != nil {
		logger.Error(err)
		_, _ = c.Reply("获取用户数据错误，可能是因为 id 已失效，请检查白菜菊花 id 是否正确")
	}

	allChefs, err := dao.FindAllChefs()
	if err != nil {
		logger.Errorw("获取所有厨师信息失败", "err", err.Error())
		_, _ = c.Reply("获取所有厨师信息失败")
		return
	}
	chefs, note := filterChefsByName(allChefs, chefName)
	if note != "" {
		logger.Info("厨师查询失败:", note)
		_, _ = c.Reply(note)
		return
	}
	logger.Debugw("识别到厨师", "chefs", chefs)
	if len(chefs) > 1 {
		chefsName := ""
		for _, chef := range chefs {
			chefsName += " " + chef.Name
		}
		_, _ = c.Reply("识别到多位厨师：" + chefsName)
		return
	}

	chef := chefs[0]
	chefGotMap, err := userGotChefData(userData)
	if err != nil {
		logger.Errorw("解析用户获取厨师数据失败：" + err.Error())
		return
	}
	response.ChefGot = chefGotMap[chef.ChefId]

	reciptGotMap, err := userGotReciptData(userData)
	if err != nil {
		logger.Errorw("解析用户获取菜谱数据失败：" + err.Error())
		return
	}

	allRecipes, err := dao.FindAllRecipes()
	if err != nil {
		logger.Error("获取所有菜谱失败")
		return
	}

	var gotRecipes []database.Recipe
	for _, recipt := range allRecipes {
		if reciptGotMap[recipt.RecipeId] {
			gotRecipes = append(gotRecipes, recipt)
		}
	}

	// 给厨师叠 buff
	addBuff(&chef, userData)

	// 分析修炼条件

	// 获取厨师修炼任务
	quests, err := dao.FindQuestsWithIds(chef.UltimateGoals)
	if err != nil {
		logger.Error("获取厨师修炼任务失败")
		return
	}

	for _, quest := range quests {
		logger.Debugf("\t%-2d. %s\n", quest.QuestId, quest.Goal)
		if len(quest.Conditions) == 0 {
			continue
		}
		c := quest.Conditions[0]
		if !(len(c.Skill) > 0 && c.Rarity > 0 && c.Rank > 0 && c.Num > 0) {
			continue
		}

		// 找菜
		var bucket = make(map[int][]database.Recipe) // 神差值
		var cando []database.Recipe
		for _, recipe := range gotRecipes {
			if recipe.Rarity != c.Rarity {
				continue
			}
			switch {
			case c.Skill == "stirfry" && recipe.Stirfry == 0,
				c.Skill == "boil" && recipe.Boil == 0,
				c.Skill == "knife" && recipe.Cut == 0,
				c.Skill == "fry" && recipe.Fry == 0,
				c.Skill == "bake" && recipe.Bake == 0,
				c.Skill == "steam" && recipe.Steam == 0:
				continue
			}

			if chef.Stirfry >= c.Rank*recipe.Stirfry &&
				chef.Boil >= c.Rank*recipe.Boil &&
				chef.Cut >= c.Rank*recipe.Cut &&
				chef.Fry >= c.Rank*recipe.Fry &&
				chef.Bake >= c.Rank*recipe.Bake &&
				chef.Steam >= c.Rank*recipe.Steam {
				cando = append(cando, recipe)
			} else {
				diffVal := 0
				r := recipe
				r.Stirfry = ((chef.Stirfry - c.Rank*recipe.Stirfry) >> (strconv.IntSize - 1)) & 1 * (c.Rank*recipe.Stirfry - chef.Stirfry)
				r.Boil = ((chef.Boil - c.Rank*recipe.Boil) >> (strconv.IntSize - 1)) & 1 * (c.Rank*recipe.Boil - chef.Boil)
				r.Cut = ((chef.Cut - c.Rank*recipe.Cut) >> (strconv.IntSize - 1)) & 1 * (c.Rank*recipe.Cut - chef.Cut)
				r.Fry = ((chef.Fry - c.Rank*recipe.Fry) >> (strconv.IntSize - 1)) & 1 * (c.Rank*recipe.Fry - chef.Fry)
				r.Bake = ((chef.Bake - c.Rank*recipe.Bake) >> (strconv.IntSize - 1)) & 1 * (c.Rank*recipe.Bake - chef.Bake)
				r.Steam = ((chef.Steam - c.Rank*recipe.Steam) >> (strconv.IntSize - 1)) & 1 * (c.Rank*recipe.Steam - chef.Steam)

				diffVal = r.Stirfry + r.Boil + r.Cut + r.Fry + r.Bake + r.Steam
				if diffVal > 250 {
					continue
				}
				if arr, ok := bucket[diffVal]; !ok {
					bucket[diffVal] = []database.Recipe{r}
				} else {
					arr = append(arr, r)
				}
			}
		}
		response.CanDo = sortRecipeByTime(cando)

		keys := make([]int, 0, len(bucket))
		for k := range bucket {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for _, k := range keys {
			for _, recipe := range bucket[k] {
				response.CanNotDo = append(response.CanNotDo, recipe)
			}
		}
	}

	got := "未拥有"
	if response.ChefGot {
		got = "已有"
	}
	msg := fmt.Sprintf("【%s】%s", chef.Name, got)
	msg += fmt.Sprintf("\n修炼技能：%s", chef.UltimateSkillDesc)
	msg += fmt.Sprintf("\n修炼任务：")
	for i, quest := range quests {
		msg += fmt.Sprintf("\n[%d] %s", i+1, quest.Goal)
		if len(quest.Conditions) > 0 {
			c := quest.Conditions[0]
			if !(len(c.Skill) > 0 && c.Rarity > 0 && c.Rank > 0 && c.Num > 0) {
				continue
			}
			msg += fmt.Sprintf("\n菜谱推荐：")
			for i, recipe := range response.CanDo {
				if i >= 5 {
					break
				}
				msg += fmt.Sprintf("\n%s（%s） ✔️", recipe.Name, (time.Duration(recipe.TotalTime) * time.Second).String())
			}
			for i, recipe := range response.CanNotDo {
				if i >= 10 {
					msg += "\n......"
					break
				}
				var items []string
				if recipe.Stirfry > 0 {
					items = append(items, fmt.Sprintf("炒:%d", recipe.Stirfry))
				}
				if recipe.Boil > 0 {
					items = append(items, fmt.Sprintf("煮:%d", recipe.Boil))
				}
				if recipe.Cut > 0 {
					items = append(items, fmt.Sprintf("切:%d", recipe.Cut))
				}
				if recipe.Fry > 0 {
					items = append(items, fmt.Sprintf("炸:%d", recipe.Fry))
				}
				if recipe.Bake > 0 {
					items = append(items, fmt.Sprintf("烤:%d", recipe.Bake))
				}
				if recipe.Steam > 0 {
					items = append(items, fmt.Sprintf("蒸:%d", recipe.Steam))
				}
				msg += fmt.Sprintf("\n%s 神差值：%s", recipe.Name, strings.Join(items, " "))
			}
		}
	}
	c.Reply(msg)
}

func userGotChefData(userData userdata.UserData) (map[int]bool, error) {
	var gotMap = make(map[int]bool, 0)
	bs, err := userData.ChefGot.MarshalJSON()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bs, &gotMap)
	return gotMap, err
}

func userGotReciptData(userData userdata.UserData) (map[int]bool, error) {
	var gotMap = make(map[int]bool, 0)
	bs, err := userData.RepGot.MarshalJSON()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bs, &gotMap)
	return gotMap, err
}

func StringMustInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func addBuff(chef *database.Chef, data userdata.UserData) {
	male := StringMustInt(data.UserUltimate.Male)
	if chef.Gender == 0 || chef.Gender == 1 {
		// 男
		chef.Stirfry += male
		chef.Boil += male
		chef.Cut += male
		chef.Fry += male
		chef.Bake += male
		chef.Steam += male
	}
	female := StringMustInt(data.UserUltimate.Female)
	if chef.Gender == 0 || chef.Gender == 2 {
		// 女
		chef.Stirfry += female
		chef.Boil += female
		chef.Cut += female
		chef.Fry += female
		chef.Bake += female
		chef.Steam += female
	}
	all := StringMustInt(data.UserUltimate.All)
	chef.Stirfry += StringMustInt(data.UserUltimate.Stirfry) + all
	chef.Boil += StringMustInt(data.UserUltimate.Boil) + all
	chef.Cut += StringMustInt(data.UserUltimate.Knife) + all
	chef.Fry += StringMustInt(data.UserUltimate.Fry) + all
	chef.Bake += StringMustInt(data.UserUltimate.Bake) + all
	chef.Steam += StringMustInt(data.UserUltimate.Steam) + all
}

func sortRecipeByTime(recipes []database.Recipe) []database.Recipe {
	for i := 0; i < len(recipes); i++ {
		for j := 0; j < len(recipes)-i-1; j++ {
			if recipes[j].Time > recipes[j+1].Time {
				recipes[j], recipes[j+1] = recipes[j+1], recipes[j]
			}
		}
	}
	return recipes
}
