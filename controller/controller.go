/**
 @author:way
 @date:2021/12/3
 @note
**/

package controller

import (
	"github.com/gin-gonic/gin"
	"redisData/dao/redis"
	"redisData/logic"
	"redisData/model"
	"redisData/pkg/logger"
)

func GetDataHandle(c *gin.Context)  {
	//获取参数
	dataType := c.Query("dataType")
	if len(dataType) <= 0 {
		c.JSON(500,gin.H{
			"msg" : "dataType为必填参数",
		})
	}
	//逻辑处理 1.根据ID找到前缀
	data, err := logic.GetKeysByPfx(dataType)
	if err != nil {
		logger.Info(err)
		return 
	}

	//返回数据
	c.JSON(200,gin.H{
		"data" : data,
	})
}

func GetMarketPriceHandle(c *gin.Context)  {
	//获取参数
	//逻辑处理
	data, err := redis.GetData("eggMarket")
	if err != nil {
		logger.Info(data)
		return 
	}
	//返回数据
	c.JSON(200,gin.H{
		"data":data,
	})
}

func SetStartParamHandler(c *gin.Context)  {
	//获取参数
	var p model.ParamStart
	err := c.Bind(&p)
	if err != nil {
		logger.Info(err)
		return 
	}
	if p.Buy==0||p.Sale==0||p.Safe ==0{
		c.JSON(500,gin.H{
			"msg":"缺少相关参数",
		})
	}
	//逻辑处理
	err1 := redis.CreateDurableKey("buy", p.Buy)
	if err != nil {
		logger.Info(err1)
		return 
	}
	err2 := redis.CreateDurableKey("sale", p.Sale)
	if err != nil {
		logger.Info(err2)
		return
	}
	err3 := redis.CreateDurableKey("safe", p.Safe)
	if err != nil {
		logger.Info(err3)
		return 
	}
	//返回参数
	c.JSON(200,gin.H{
		"msg":"success",
	})
}

func GetBuyDataHandle(c *gin.Context)  {
	
}

