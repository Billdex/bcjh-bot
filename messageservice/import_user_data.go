package messageservice

import (
	"bcjh-bot/dao"
	"bcjh-bot/model/database"
	"bcjh-bot/model/userdata"
	"bcjh-bot/scheduler"
	"bcjh-bot/util/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func ImportUserData(c *scheduler.Context) {
	bcjhId, err := strconv.Atoi(c.PretreatedMessage)
	if err != nil {
		_, _ = c.Reply("白菜菊花个人数据 ID 格式有误")
		return
	}

	qq := c.GetSenderId()
	data, err := dao.FindUserDataWithUserId(qq)
	if err != nil {
		logger.Error(err.Error())
		_, _ = c.Reply("导入失败！")
		return
	} else if data.BcjhID == bcjhId {
		_, _ = c.Reply("数据已存在，无需重复导入")
		return
	}

	bytes, err := downloadUserData(bcjhId)
	if err != nil {
		logger.Errorw("下载白菜菊花个人数据失败，"+err.Error(), "bcjhId", bcjhId)
		_, _ = c.Reply("导入失败！下载数据失败，请过一会重试")
		return
	}

	var r userdata.Response
	if err = json.Unmarshal(bytes, &r); err != nil {
		logger.Errorw("解析白菜菊花个人数据失败，"+err.Error(), "bcjhId", bcjhId)
		_, _ = c.Reply("导入失败！解析个人数据失败")
		return
	}

	if r.Result == false {
		_, _ = c.Reply("导入失败！" + r.Msg)
		return
	}

	err = dao.SetUserData(database.UserData{QQ: qq, User: r.User, BcjhID: r.Id, Data: r.Data, CreateTime: r.CreateTime})
	if err != nil {
		logger.Error(err.Error())
		_, _ = c.Reply("导入失败！")
		return
	}

	_, _ = c.Reply(fmt.Sprintf("导入【%s】的数据成功！", r.User))
}

func downloadUserData(bcjhId int) ([]byte, error) {
	url := fmt.Sprintf("https://bcjh.xyz/api/download_data?id=%d", bcjhId)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
