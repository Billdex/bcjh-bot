package util

import (
	"fmt"
	"os"
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
			time += fmt.Sprintf("%d小时", hour)
		}
		if minute > 0 {
			time += fmt.Sprintf("%d分", minute)
		}
		if second > 0 {
			time += fmt.Sprintf("%d秒", second)
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
