package middleware

import (
	"bcjh-bot/scheduler"
	"regexp"
	"strings"
)

// MergeRepeatSpace 合并重复的空白字符
func MergeRepeatSpace(c *scheduler.Context) {
	reg := regexp.MustCompile("\\s+")
	c.PretreatedMessage = reg.ReplaceAllString(strings.TrimSpace(c.PretreatedMessage), " ")
	c.Next()
}
