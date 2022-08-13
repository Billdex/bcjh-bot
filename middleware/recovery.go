package middleware

import (
	"bcjh-bot/scheduler"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
)

// Recovery 防住某些 panic
func Recovery(c *scheduler.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("some query panic! err: %v, raw message: %s", err, c.GetRawMessage())
			_, _ = c.Reply(e.SystemErrorNote)
			c.Abort()
		}
	}()
	c.Next()
}
