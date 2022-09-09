package middleware

import "bcjh-bot/scheduler"

// Helper 功能帮助说明，在匹配到功能但没有写参数的时候提供一个简要说明
func Helper(s string) func(c *scheduler.Context) {
	return func(c *scheduler.Context) {
		if c.PretreatedMessage == "" {
			_, _ = c.Reply(s)
			c.Abort()
			return
		}
	}
}
