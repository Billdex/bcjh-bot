package util

import (
	"fmt"
)

// PaginationOutput 对数据结果进行分页输出, 超出总页数的按照最后一页输出，页数小于 1 的按照第一页输出
func PaginationOutput[T any](list []T, page int, perPage int, title string, output func(item T) string) string {
	var msg string = title
	maxPage := (len(list)-1)/perPage + 1
	if len(list) > perPage {
		if page > maxPage {
			page = maxPage
		}
		msg = fmt.Sprintf("%s (%d/%d)", title, page, maxPage)
	}
	if page <= 0 {
		page = 1
	}
	for i := (page - 1) * perPage; i < page*perPage && i < len(list); i++ {
		msg += fmt.Sprintf("\n%s", output(list[i]))
	}
	if page < maxPage {
		msg += "\n......"
	}
	return msg
}
