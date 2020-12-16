package server

import (
	"bcjh-bot/service"
	"bcjh-bot/util"
	"strings"
)

// 功能：前缀过滤器
// 入参: 消息内容, 前缀
// 返回值: 正文内容, 前缀是否符合
func PrefixFilter(str string, prefix string) (string, bool) {
	str = strings.TrimSpace(str)
	hasPrefix := strings.HasPrefix(str, prefix)
	if !hasPrefix {
		return "", false
	}
	return str[len(prefix):], true
}

// 功能: 指令过滤器
// 入参: 文本内容, map(指令-处理方法)
// 返回值: 具体处理方法的函数指针, 参数列表
func InstructionFilter(str string, instructions map[string]service.InstructionHandlerFunc) (service.InstructionHandlerFunc, []string) {
	str = strings.TrimSpace(str)
	for instruction, handler := range instructions {
		if strings.HasPrefix(str, instruction) {
			args := strings.Split(str[len(instruction):], util.ArgsSplitCharacter)
			if args[0] == "" {
				args = make([]string, 0)
			}
			return handler, args
		}
	}
	return nil, make([]string, 0)
}
