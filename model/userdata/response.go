package userdata

import (
	"bcjh-bot/config"
	"bcjh-bot/util"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Result     bool      `json:"result"`
	Id         int       `json:"id"`
	User       string    `json:"user"`
	Data       string    `json:"data"`
	CreateTime time.Time `json:"create_time"`
}

func (r *Response) ParseUserData() (UserData, error) {
	var userData UserData
	err := json.Unmarshal([]byte(r.Data), &userData)
	return userData, err
}

func LoadUserData(bcjhId int) (UserData, error) {
	fname := fmt.Sprintf("%s/%d.json", config.AppConfig.Resource.UserData, bcjhId)

	// 文件不存在，则请求并下载
	if exists, _ := util.PathExists(fname); !exists {
		return DownloadUserData(bcjhId)
	}

	var userData UserData
	f, err := os.Open(fname)
	if err != nil {
		return userData, fmt.Errorf("打开 %s 文件失败，%w", fname, err)
	}

	bs, err := io.ReadAll(f)
	if err != nil {
		return userData, fmt.Errorf("读取 %s 文件失败，%w", fname, err)
	}

	var r Response
	if err = json.Unmarshal(bs, &r); err != nil {
		return userData, fmt.Errorf("解析 %s 文件失败，%w", fname, err)
	}
	return r.ParseUserData()
}

func DownloadUserData(bcjhId int) (UserData, error) {
	url := fmt.Sprintf("https://bcjh.xyz/api/download_data?id=%d", bcjhId)
	resp, err := http.Get(url)
	if err != nil {
		return UserData{}, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserData{}, err
	}

	var r Response
	err = json.Unmarshal(bs, &r)
	if err != nil {
		return UserData{}, err
	}

	fname := fmt.Sprintf("%s/%d.json", config.AppConfig.Resource.UserData, bcjhId)
	file, err := os.Create(fname)
	if err != nil {
		//fmt.Println("Error while creating the file:", err)
		return UserData{}, err
	}
	defer file.Close()

	_, err = file.Write(bs)
	if err != nil {
		return UserData{}, err
	}
	return r.ParseUserData()
}
