package util

import "fmt"

func FormatSecondToString(t int) string {
	if t < 0 {
		return ""
	}
	if t == 0 {
		return "0秒"
	}
	var time string
	hour := t / 3600
	minute := t % 3600 / 60
	second := t % 3600 % 60
	if hour > 0 {
		time += fmt.Sprintf("%d小时", hour)
	}
	if minute > 0 {
		time += fmt.Sprintf("%d分钟", minute)
	}
	if second > 0 {
		time += fmt.Sprintf("%d秒", second)
	}
	return time
}
