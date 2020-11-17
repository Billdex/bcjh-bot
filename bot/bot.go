package bot

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func PostMsg(byteMsg []byte, url string) error {
	request, err := http.NewRequest("POST", url, bytes.NewReader(byteMsg))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("CoolQ Response Info:", string(body))
	return nil
}

//func SendProvateMsg(msg PrivateMsg) error {
//
//}
//
//func SendGroupMsg(msg GroupMsg) (err error) {
//	byteMsg, err := json.Marshal(&msg)
//	if err != nil {
//		return err
//	}
//	url := "http://"+config.AppConfig.CQHTTP.Host+":"+strconv.Itoa(config.AppConfig.CQHTTP.Port) + "/send_group_msg"
//	err = SendMsg(byteMsg, url)
//	if err != nil {
//		return err
//	}
//	log.Println("Send Group Msg Success!", msg)
//	return nil
//}
