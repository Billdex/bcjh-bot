package main

import (
	"bcjh-bot/model/database"
	"bcjh-bot/util"
	"bcjh-bot/util/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetFeedbacks(c *gin.Context) {
	var params struct {
		Status database.FeedbackStatus `form:"status"`
		Page   int                     `form:"page"`
	}
	var req struct {
		Feedbacks []database.Feedback `json:"feedbacks"`
		Total     int                 `json:"total"`
	}
	err := c.ShouldBind(&params)
	if err != nil {
		logger.Error("查询参数绑定出错", err)
		c.JSON(http.StatusOK, gin.H{
			"status":  400,
			"message": "参数有误",
			"data":    nil,
		})
		return
	}

	total, err := db.Where("status = ?", params.Status).Count(&database.Feedback{})
	if err != nil {
		// 查询反馈出错时只记录日志不返回数据
		logger.Error("查询反馈总数失败", err)
	}
	feedbacks := make([]database.Feedback, 0)
	err = db.Where("status = ?", params.Status).Limit(20, (util.LimitPage(params.Page, 1, int(total/20)+1)-1)*20).Find(&feedbacks)
	if err != nil {
		logger.Error("查询反馈数据失败", err)
		c.JSON(http.StatusOK, gin.H{
			"status":  500,
			"message": "数据查询失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "ok",
		"data":    req,
	})
}
