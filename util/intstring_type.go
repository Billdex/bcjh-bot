package util

import (
	"encoding/json"
	"strconv"
)

// IntString 一个基准类型为 int 的类型，可以由空字符串、null 反序列化得到
type IntString int

func (is *IntString) UnmarshalJSON(data []byte) error {
	var i int
	if err := json.Unmarshal(data, &i); err != nil {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		if s == "" || s == "null" || s == "0" {
			*is = 0
			return nil
		}
		i, err = strconv.Atoi(s)
		if err != nil {
			return err
		}
	}
	*is = IntString(i)
	return nil
}

func (is *IntString) Marshal() ([]byte, error) {
	return json.Marshal(*is)
}
