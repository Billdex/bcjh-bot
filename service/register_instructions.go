package service

type InstructionHandlerFunc func([]string)

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

func RegisterInstructions() {
	Ins = NewInstructions()
	Ins.Bind("更新数据", Update)
	Ins.Bind("厨师", ChefQuery)
	Ins.Bind("厨具", EquipmentQuery)
	Ins.Bind("菜谱", RecipeQuery)
}
