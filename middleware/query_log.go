package middleware

import (
	"bcjh-bot/scheduler"
	"bcjh-bot/util/logger"
	"fmt"
	"time"
)

// QueryLog 查询日志中间件，记录每条查询的内容与耗时
func QueryLog(c *scheduler.Context) {
	start := time.Now()

	c.Next()

	latency := time.Now().Sub(start)
	bot := c.GetBot()
	group := c.GetGroupId()
	sender := c.GetSenderId()
	logMsg := fmt.Sprintf("[Query] bot:%12d | %13v | group:%12d | sender:%12d | %s",
		bot.BotId, latency, group, sender, c.GetRawMessage())
	if c.WarnMessage != "" {
		logMsg = fmt.Sprintf("%s | %s", logMsg, c.WarnMessage)
	}
	logger.Infof(logMsg)
}
