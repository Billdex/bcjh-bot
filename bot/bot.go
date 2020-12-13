package bot

import (
	"bcjh-bot/config"
	"bcjh-bot/model/onebot"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func OneBotPost(byteMsg []byte, url string) error {
	request, err := http.NewRequest("POST", url, bytes.NewReader(byteMsg))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case util.OneBotTokenEmpty:
		return errors.New("未提供access token")
	case util.OneBotTokenWrong:
		return errors.New("access token有误")
	case util.OneBotContentTypeError:
		return errors.New("不支持的Content-Type")
	case util.OneBotTextFormatError:
		return errors.New("请求的正文格式不正确")
	case util.OneBotAPINotFound:
		return errors.New("请求的API不存在")
	case util.OneBotStatusOK:
		return nil

	default:
		return err
	}
}

func SendMessage(c *onebot.Context, msg string) error {
	switch c.MessageType {
	case util.OneBotMessagePrivate:
		privateMsg := onebot.PrivateMsg{
			UserId:     c.UserId,
			Message:    msg,
			AutoEscape: false,
		}
		return SendPrivateMsg(privateMsg)
	case util.OneBotMessageGroup:
		groupMsg := onebot.GroupMsg{
			GroupId:    c.GroupId,
			Message:    msg,
			AutoEscape: false,
		}
		return SendGroupMsg(groupMsg)
	default:
		return errors.New("未知类型")
	}
}

func SendPrivateMsg(msg onebot.PrivateMsg) error {
	byteMsg, err := json.Marshal(&msg)
	if err != nil {
		return err
	}
	baseUrl := "http://" + config.AppConfig.OneBot.Host + ":" + strconv.Itoa(config.AppConfig.OneBot.Port)
	url := baseUrl + "/send_private_msg"
	logger.Debug("尝试发送一条私聊消息:", msg)
	err = OneBotPost(byteMsg, url)
	return err
}

func SendGroupMsg(msg onebot.GroupMsg) error {
	byteMsg, err := json.Marshal(&msg)
	if err != nil {
		return err
	}
	baseUrl := "http://" + config.AppConfig.OneBot.Host + ":" + strconv.Itoa(config.AppConfig.OneBot.Port)
	url := baseUrl + "/send_group_msg"
	logger.Debug("尝试发送一条群聊消息:", msg)
	err = OneBotPost(byteMsg, url)
	return err
}

func GetCQImage(path string, pathType string) string {
	switch pathType {
	case "file":
		return "[CQ:image,file=file:///" + path + "]"
	case "url":
		return "[CQ:image,url=" + path + "]"
	case "base64":
		return "[CQ:image,base64=base64://" + path + "]"
	default:
		return ""
	}
}
