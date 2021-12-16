/**
 @author:way
 @date:2021/12/16
 @note
**/

package main

import (
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logic"
	"redisData/pkg/logger"
	"redisData/setting"
	"strconv"
)

func init() {
	// 定义日志目录
	logger.Init("buy")
	// 初始化 viper 配置
	if err := setting.Init(""); err != nil {
		logger.Info("viper init fail")
		logger.Error(err)
		return
	}
	// 初始化MySQL
	mysql.InitMysql()
	//初始化redis
	if err := redis.InitClient(); err != nil {
		logger.Info("init redis fail err")
		logger.Error(err)
		return
	}
}



func main() {

	//开始缓存
	for {
		all := redis.GetHashDataAll("baby:ConfigStopAutoBuy")
		//判断总开关状态
		data1 := redis.GetHashDataAll("baby:ConfigBuy")
		percent := data1["percent"]
		percentFloat,ParseFloatErr := strconv.ParseFloat(percent,64)
		if ParseFloatErr != nil{
			logger.Error(ParseFloatErr)
			continue
		}

		//全自动市场价等于运算出来的
		if data1["status"] == "1" && all["Super"] == "2" {
			//获取市场价
			marketPriceAuto, GetDataErr := redis.GetData("baby:marketPrice")
			if GetDataErr != nil {
				logger.Info(GetDataErr)
				return
			}
			marketPriceFloat,err := strconv.ParseFloat(marketPriceAuto,64)
			if err != nil{
				logger.Error(err)
			}
			logic.StartBuy(marketPriceFloat,percentFloat)
		}
		//半自动市场价根据配置问价设定值
		if data1["status"] == "1" && all["Super"] == "1" {
			marketPrice := data1["market_price"]
			market_price_float,FloatErr := strconv.ParseFloat(marketPrice,64)
			if FloatErr != nil{
				logger.Error(FloatErr)
				continue
			}
			logic.StartBuy(market_price_float,percentFloat)
		}
	}
}