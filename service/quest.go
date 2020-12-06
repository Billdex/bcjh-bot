package service

import (
	"bcjh-bot/bot"
	"bcjh-bot/model/database"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"fmt"
	"strconv"
	"strings"
)

func QuestQuery(c *onebot.Context, args []string) {
	logger.Info("任务查询，参数:", args)
	maxLen := util.MaxQueryListLength
	// 无参数的情况
	if len(args) == 0 {
		if err := bot.SendMessage(c, questHelp()); err != nil {
			logger.Error("发送信息失败!", err)
		}
		return
	}

	// 建立会话
	Session := database.DB.Select("*")

	// 1. 第一个参数
	if args[0] != "" {
		if arg, ok := StringContainsAny(args[0], []string{"主线", "支线"}); ok { // 主线或支线任务
			Session.Where("type like ?", arg+"%")
		} else if args[0] == "限时" { // 限时或所有
			Session.Where("type <> '主线任务' and type <> '支线任务'")
		} else {
			Session.Where("type like ?", "%"+arg+"%")
		}
	}

	// 2. 确定任务 id
	if len(args) > 1 && args[1] != "" {
		if _, ok := StringContainsAny(args[0], []string{"主线"}); ok { // 主线任务
			id, _ := strconv.Atoi(args[1])
			if id > 700 {
				_ = bot.SendMessage(c, "主线任务目前只有 700 个哦")
				return
			}
			// 如果是查询主线区间
			if len(args) > 2 && args[2] != "" {
				left := id
				right, _ := strconv.Atoi(args[2])

				if right > 700 {
					_ = bot.SendMessage(c, "主线任务目前只有 700 个哦")
					return
				}
				if left > right {
					left, right = right, left
				}
				// 区间不能过大，不然消息太长
				if right-left > maxLen-1 { // 可以查 5 条
					_ = bot.SendMessage(c, "区间跨度不能太大哦，不然消息会很长")
					right = left + maxLen - 1
				}
				Session.Where("quest_id >= ? and quest_id <= ?", left, right)
			} else {
				// 限制查询的 id 在 700 以内
				Session.Where("quest_id = ?", id)
			}
		} else { // 支线或限时
			Session.Where("quest_id_disp = ?", args[1])
		}
	} else {
		_ = bot.SendMessage(c, "要指定一下任务 id 哦")
		return
	}

	// 查询得到结果
	quests := make([]database.Quest, 0)
	err := Session.Find(&quests)

	if err != nil {
		logger.Error("查询数据库出错!", err)
		_ = bot.SendMessage(c, "查询数据失败!")
		return
	}

	var msg string
	switch {
	case len(quests) == 0:
		msg = "哎呀，好像找不到呢!"
	default:
		msg = makeQuestsString(quests)
	}
	_ = bot.SendMessage(c, msg)
}

func makeQuestsString(quests []database.Quest) string {
	sb := strings.Builder{}

	for count, quest := range quests {
		sb.WriteString(fmt.Sprintf("[%s] ", quest.Type))
		if quest.Type == "主线任务" {
			sb.WriteString(fmt.Sprintf("%v", quest.QuestId))
		} else {
			sb.WriteString(fmt.Sprintf("%v", quest.QuestIdDisp))
		}
		sb.WriteString(fmt.Sprintf("\n要求：%s", quest.Goal))
		sb.WriteString("\n奖励：")
		if len(quest.Rewards) == 0 {
			sb.WriteString("无")
		} else {
			for i, reward := range quest.Rewards {
				if reward.Quantity == "" {
					sb.WriteString(fmt.Sprintf("%s", reward.Name))
				} else {
					sb.WriteString(fmt.Sprintf("%s*%v", reward.Name, reward.Quantity))
				}
				if i != len(quest.Rewards)-1 {
					sb.WriteString(", ")
				}
			}
		}
		if count != len(quests)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
