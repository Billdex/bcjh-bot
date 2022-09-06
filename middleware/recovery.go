package middleware

import (
	"bcjh-bot/scheduler"
	"bcjh-bot/util/e"
	"bcjh-bot/util/logger"
	"runtime"
)

// Recovery 防住某些 panic
func Recovery(c *scheduler.Context) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			logger.Errorf("some query panic! err: %+v, raw message: %s, stack: %s", err, c.GetRawMessage(), string(buf[:n]))
			_, _ = c.Reply(e.SystemErrorNote)
			c.Abort()
		}
	}()
	c.Next()
}
