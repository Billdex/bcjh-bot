package util

//入参：文本str, 前缀prefix
//返回值：去除前缀后的内容, 是否包含前缀
func PrefixFilter(str string, prefix string) (string, bool) {
	//len(s) >= len(prefix) && s[0:len(prefix)] == prefix
	hasPrefix := len(str) > len(prefix) && str[0:len(prefix)] == prefix
	if !hasPrefix {
		return "", false
	}
	return str[len(prefix):len(str)], true
}
