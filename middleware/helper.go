package middleware

import "bcjh-bot/scheduler"

// Helper 功能帮助说明，在匹配到功能但没有写参数的时候提供一个简要说明
func Helper(fn func() string) func(c *scheduler.Context) {
	return func(c *scheduler.Context) {
		if c.PretreatedMessage == "" {
			_, _ = c.Reply(fn())
			c.Abort()
			return
		}
	}
}
