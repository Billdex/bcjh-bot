package util

import (
	"fmt"
	"os"
	"strings"
)

func FormatSecondToString(t int) string {
	if t < 0 {
		return ""
	} else if t == 0 {
		return "0秒"
	} else {
		var time string
		hour := t / 3600
		minute := t % 3600 / 60
		second := t % 3600 % 60
		if hour > 0 {
			time += fmt.Sprintf("%d 小时 ", hour)
		}
		if minute > 0 {
			time += fmt.Sprintf("%d 分 ", minute)
		}
		if second > 0 {
			time += fmt.Sprintf("%d 秒", second)
		}
		return time
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func HasPrefixIn(s string, prefix ...string) bool {
	for _, p := range prefix {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
