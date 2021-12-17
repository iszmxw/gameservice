package routes

import (
	"github.com/gin-gonic/gin"
	"redisData/controller/baby"
)

// RegisterWebRoutes 注册路由 baby 游戏路由

var babyController = new(baby.Controller)

func RegisterWebRoutes(router *gin.RouterGroup) {
	// 买卖数据筛选
	router.GET("/buy_and_sell", babyController.BuyAndSellHandler)
}
