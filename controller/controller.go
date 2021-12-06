/**
 @author:way
 @date:2021/12/3
 @note
**/

package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logic"
	"redisData/model"
	"redisData/pkg/logger"
)


// GetDataHandle 返回id详情列表前端用
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

//GetMarketPriceHandle 获取市场价格
func GetMarketPriceHandle(c *gin.Context)  {
	//获取参数
	var p model.ParamTypeId
	Berr := c.Bind(&p)
	if Berr != nil {
		logger.Info(Berr)
		return
	}
	//逻辑处理
	d := mysql.GetAssetName(p.TypeId)
	marketKey := fmt.Sprintf("%s.MarketPrice",d.TypeName)
	data, err := redis.GetData(marketKey)
	if err != nil {
		logger.Info(data)
		return 
	}
	//返回数据
	c.JSON(200,gin.H{
		"data":data,
	})
}

//SetStartParamHandler 设置启动参数
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
		return
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

// GetBuyDataHandle 返回买入卖出的数据
func GetBuyDataHandle(c *gin.Context)  {
	var (
		p model.ParamBuyStatus
		result []model.RespBuy
	)

	Berr := c.Bind(&p)
	if Berr != nil {
		logger.Info(Berr)
		return
	}
	//通过查询最新10条买入数据
		data := mysql.GetBuyData(p.Status)
		result = make([]model.RespBuy,len(data))
		for i,v := range data{
			result[i].Gid = v.Gid
			result[i].Name = v.Name
			result[i].Count = v.Count
			result[i].TokenId = v.TokenId
			result[i].FixedPrice = v.FixedPrice
			result[i].TotalPrice = v.FixedPrice *float64(v.Count)
			result[i].CreateTime = v.CreatedAt
			result[i].Type = v.Type
			result = append(result,result[i])
		}
		c.JSON(200,gin.H{
			"data": result,
		})
}



//前端系统监控使用

//SetMngRiskHandle 设置风控
func SetMngRiskHandle(c *gin.Context)  {
	//获取参数
	var p model.ParamRiskMng
	err := c.Bind(&p)
	if err != nil {
		logger.Info(err)
		return
	}
	if p.Situation==""||p.TimeLevel==0||p.Percentage ==0||p.OperationType==0||p.Status==0{
		c.JSON(500,gin.H{
			"msg":"缺少相关参数",
		})
		return
	}
	//把数据存进redis 中的哈希表
	m := make(map[string]interface{})
	m["Situation"]= p.Situation
	m["TimeLevel"]=p.TimeLevel
	m["Percentage"]=p.Percentage
	m["OperationType"]=p.OperationType
	m["Status"]=p.Status
	logger.Info(m)
	redis.CreatHashKey(fmt.Sprintf("risk:%s",p.Situation),m)
	//返回参数
	c.JSON(200,gin.H{
		"msg":"success",
	})

}

//SetBuyAndSaleHandle 设置买入卖出百分比
func SetBuyAndSaleHandle(c *gin.Context){
	//获取参数
	var p model.ParamBuyAndSale
	err := c.Bind(&p)
	if err != nil {
		logger.Info(err)
		return
	}
	if p.ProductName==""||p.RisePercentage==0||p.FallPercentage ==0||p.Status==0{
		c.JSON(500,gin.H{
			"msg":"缺少相关参数",
		})
		return
	}
	//把数据存进redis 中的哈希表
	m := make(map[string]interface{})
	m["ProductName"]= p.ProductName
	m["RisePercentage"]=p.RisePercentage
	m["FallPercentage"]=p.FallPercentage
	m["Status"]=p.Status
	logger.Info(m)
	redis.CreatHashKey(fmt.Sprintf("buyAndSale:%s",p.ProductName),m)
	//返回参数
	c.JSON(200,gin.H{
		"msg":"success",
	})
}

//SetParamOnOffHandle 设置买入卖出总开关
func SetParamOnOffHandle(c *gin.Context)  {
	//获取参数
	var p model.ParamOnOff
	err := c.Bind(&p)
	if err != nil {
		logger.Info(err)
		return
	}
	if p.CrlName==""||p.Super==0{
		c.JSON(500,gin.H{
			"msg":"缺少相关参数",
		})
		return
	}
	//把数据存进redis 中的哈希表
	m := make(map[string]interface{})
	m["CrlName"]= p.CrlName
	m["Super"]=p.Super
	logger.Info(m)
	redis.CreatHashKey(fmt.Sprintf("buyAndSale:%s", p.CrlName),m)
	//返回参数
	c.JSON(200,gin.H{
		"msg":"success",
	})
}

//GetScriptStatusHandle 获取脚本运行的状态 没写
func GetScriptStatusHandle(c *gin.Context)  {
	//逻辑

}