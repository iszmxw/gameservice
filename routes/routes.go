package routes

import (
	"github.com/gin-gonic/gin"
	"redisData/controller"
	"redisData/middleware"
)

func SetUp() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors()) //跨域
	r.Use(middleware.TraceLogger())
	//r.Use(middleware.TlsHandler())  // 支持wss
	 // 日志上下文进行绑定追踪
	//查询，查询redis上的数据，返回给前端
	//websocket
	v1 := r.Group("/api/game")
	v1.GET("/getAssetsDetailList",controller.GetDataHandle)
	v1.GET("/getMarketPrice",controller.GetMarketPriceHandle)
	v1.GET("/setStartParam",controller.SetStartParamHandler)
	v1.GET("/getBuyData",controller.GetBuyDataHandle)
	v1.GET("/setMngRisk",controller.SetMngRiskHandle)
	//v1.GET("/setBuyAndSale",controller.SetBuyAndSaleHandle)
	v1.GET("/setParamOnOff",controller.SetParamOnOffHandle)
	v1.GET("/getScriptStatus",controller.GetScriptStatusHandle)
	v1.GET("/getRiskMonitor",controller.GetRiskMonitorHandle)
	v1.GET("/getMarketPriceLine",controller.GetMarketPriceLineHandle)
	v1.GET("/getIncome",controller.GetIncomeHandle)
	v1.GET("/getAssetType",controller.GetAssetType)
	v1.GET("/setBuySet",controller.SetBuySetHandle)
	v1.GET("/setSaleSet",controller.SetSaleSetHandle)
	v1.GET("/setRiskPotion",controller.SetRiskPotionHandle)
	v1.GET("/getRiskPotion",controller.GetRiskPotionHandle)
	v1.GET("/setSellingRate",controller.SetSellingRateHandle)
	v1.GET("/getSellingRate",controller.GetSellingRateHandle)
	v1.GET("/getProportion",controller.GetProportionHandle)
	//获取买卖数据,这个接口有问题
	v1.GET("/getBuyAndSaleData",controller.GetBuyAndSaleHandle)
	return r
}
