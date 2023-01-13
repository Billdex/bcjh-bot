package main

import (
	"bcjh-bot/model/database"
	"bcjh-bot/util/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetFeedbacks(c *gin.Context) {
	var req struct {
		Status   database.FeedbackStatus `form:"status" json:"status"`
		Page     int                     `form:"page" json:"page"`
		PageSize int                     `form:"page_size" json:"page_size"`
	}
	var resp struct {
		Feedbacks []database.Feedback `json:"feedbacks"`
		Total     int64               `json:"total"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Error("查询参数绑定出错", err)
		c.JSON(http.StatusOK, gin.H{
			"status":  400,
			"message": "参数有误",
			"data":    resp,
		})
		return
	}

	total, err := db.Where("status = ?", req.Status).Count(&database.Feedback{})
	if err != nil {
		// 查询反馈出错时只记录日志不返回数据
		logger.Error("查询反馈总数失败", err)
	}
	feedbacks := make([]database.Feedback, 0)
	page, pageSize := LimitPaginate(req.Page, req.PageSize)
	offset := (page - 1) * pageSize
	err = db.Where("status = ?", req.Status).Limit(pageSize, offset).Find(&feedbacks)
	if err != nil {
		logger.Error("查询反馈数据失败", err)
		c.JSON(http.StatusOK, gin.H{
			"status":  500,
			"message": "数据查询失败",
			"data":    nil,
		})
		return
	}

	resp.Feedbacks = feedbacks
	resp.Total = total
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "ok",
		"data":    resp,
	})
}
