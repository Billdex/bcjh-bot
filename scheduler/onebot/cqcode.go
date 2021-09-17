package onebot

func GetCQImage(path string, pathType string) string {
	switch pathType {
	case "file":
		return "[CQ:image,file=file:///" + path + "]"
	case "url":
		return "[CQ:image,file=" + path + "]"
	case "base64":
		return "[CQ:image,file=base64://" + path + "]"
	default:
		return ""
	}
}
