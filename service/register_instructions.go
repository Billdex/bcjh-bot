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

func (i *Instructions) Bind(instruction string, handler InstructionHandlerFunc) {
	i.instructions[instruction] = handler
}

func (i *Instructions) GetInstructions() map[string]InstructionHandlerFunc {
	return i.instructions
}

var Ins Instructions

//注册指令，绑定文本指令对应的处理方法
func RegisterInstructions() {
	Ins = NewInstructions()
	Ins.Bind("更新数据", UpdateData)
	Ins.Bind("帮助", HelpGuide)
	Ins.Bind("图鉴网", GalleryWebsite)
	Ins.Bind("厨师", ChefQuery)
	Ins.Bind("厨具", EquipmentQuery)
	Ins.Bind("菜谱", RecipeQuery)
	Ins.Bind("贵客", GuestQuery)
}
