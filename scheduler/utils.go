package scheduler

import "strings"

func whatPrefixIn(s string, prefix ...string) (string, bool) {
	for _, p := range prefix {
		if strings.HasPrefix(s, p) {
			return p, true
		}
	}
	return "", false
}

func combineHandlers(baseHandlers []HandleFunc, newHandlers ...HandleFunc) []HandleFunc {
	mergedHandlers := make([]HandleFunc, 0)
	mergedHandlers = append(mergedHandlers, baseHandlers...)
	mergedHandlers = append(mergedHandlers, newHandlers...)
	return mergedHandlers
}
