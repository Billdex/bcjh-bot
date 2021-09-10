package service

import (
	"bcjh-bot/model/onebot"
)

type InstructionHandlerFunc func(*onebot.Context, []string)

type Instructions struct {
	instructions map[string]InstructionHandlerFunc
}

func NewInstructions() Instructions {
	instructions := make(map[string]InstructionHandlerFunc)
	return Instructions{
		instructions: instructions,
	}
}

func (i *Instructions) Bind(handler InstructionHandlerFunc, instructions ...string) {
	for _, instruction := range instructions {
		i.instructions[instruction] = handler
	}
}

func (i *Instructions) GetInstructions() map[string]InstructionHandlerFunc {
	return i.instructions
}

var Ins Instructions

// 注册指令，绑定文本指令对应的处理方法
func RegisterInstructions() {
	Ins = NewInstructions()
	// 主功能
	Ins.Bind(GuestQuery, "贵客", "稀有客人", "贵宾", "客人", "宾客", "稀客")
	Ins.Bind(AntiqueQuery, "符文")
	Ins.Bind(CondimentQuery, "调料")
	Ins.Bind(QuestQuery, "任务", "主线", "支线")
	Ins.Bind(TimeLimitingQuestQuery, "限时任务", "限时攻略", "限时支线", "限时任务攻略")
	Ins.Bind(UpgradeGuestQuery, "碰瓷", "升阶贵客")
	Ins.Bind(ComboQuery, "后厨", "合成")
	Ins.Bind(LaboratoryQuery, "实验室", "研究")
	Ins.Bind(StrategyQuery, "攻略")
	Ins.Bind(ExchangeQuery, "兑换码", "玉璧")

}
