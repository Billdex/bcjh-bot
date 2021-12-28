package middleware

import (
	"bcjh-bot/scheduler"
	"bcjh-bot/scheduler/onebot"
	"bcjh-bot/util/logger"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func PretreatedImage(c *scheduler.Context) {
	msg := c.PretreatedMessage
	imgList := c.GetImageList()
	for _, img := range imgList {
		imgInfo, err := c.GetBot().GetImageInfo(img)
		if err != nil {
			logger.Errorf("获取消息列表图片出错", err)
			continue
		}
		r, err := http.Get(imgInfo.Url)
		if err != nil {
			logger.Error("下载图片出错", err)
			continue
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			continue
		}
		base64Img := base64.StdEncoding.EncodeToString(body)
		r.Body.Close()
		pattern, err := regexp.Compile(fmt.Sprintf(`\[CQ:image,file=%s.*?\]`, img))
		if err != nil {
			continue
		}
		msg = pattern.ReplaceAllString(msg, onebot.GetCQImage(base64Img, "base64"))
	}
	c.PretreatedMessage = msg
	c.Next()
}
