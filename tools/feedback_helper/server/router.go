package main

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/feedbacks", GetFeedbacks)

	return r
}
