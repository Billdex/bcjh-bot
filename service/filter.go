package service

import (
	"bcjh-bot/util"
	"strings"
)

func PrefixFilter(str string, prefix string) (string, bool) {
	hasPrefix := strings.HasPrefix(str, prefix)
	if !hasPrefix {
		return "", false
	}
	return str[len(prefix):], true
}

func InstructionFilter(str string, instructions map[string]InstructionHandlerFunc) (InstructionHandlerFunc, []string) {
	for instruction, handler := range instructions {
		if strings.HasPrefix(str, instruction) {
			strArgs := strings.TrimSpace(str[len(instruction):])
			args := strings.Split(strArgs, util.ArgsSplitCharacter)
			return handler, args
		}
	}
	return nil, make([]string, 0)
}
