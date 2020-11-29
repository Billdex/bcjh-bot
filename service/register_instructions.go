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
	Ins.Bind(UpdateData, "更新数据", "更新")
	Ins.Bind(HelpGuide, "帮助", "说明")
	Ins.Bind(Feedback, "反馈", "建议")
	Ins.Bind(GalleryWebsite, "图鉴网", "图鉴")
	Ins.Bind(TermInfo, "游戏术语", "术语", "黑话")
	Ins.Bind(ChefQuery, "厨师", "厨子")
	Ins.Bind(EquipmentQuery, "厨具", "装备")
	Ins.Bind(RecipeQuery, "菜谱")
	Ins.Bind(GuestQuery, "贵客")
	Ins.Bind(CondimentQuery, "调料")
}
