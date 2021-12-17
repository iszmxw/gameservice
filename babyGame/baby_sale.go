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

	for {
		//判断总开关状态
		all := redis.GetHashDataAll("baby:ConfigStopAutoSale")
		data1 := redis.GetHashDataAll("baby:ConfigSale")
		percent := data1["percent"]
		percentFloat,ParseFloatErr := strconv.ParseFloat(percent,64)
		if ParseFloatErr != nil{
			logger.Error(ParseFloatErr)
			continue
		}
		//判断物品开关
		//同为1,设定自己买卖的脚本
		logger.Info(all["Super"])
		logger.Info(data1["status"])
		if data1["status"] == "1" && all["Super"] == "1"{
			logger.Info("执行半自动脚本")
			marketPrice := data1["market_price"]
			market_price_float,FloatErr := strconv.ParseFloat(marketPrice,64)
			if FloatErr != nil{
				logger.Error(FloatErr)
				continue
			}
			logic.StartSale(market_price_float,percentFloat)
		}
		//为1和2设置自动脚本
		if data1["status"] == "1" && all["Super"] == "2"{
			logger.Info("执行全自动自动脚本")
			//
			marketPriceAuto, GetDataErr := redis.GetData("baby:marketPrice")
			if GetDataErr != nil {
				logger.Info(GetDataErr)
				return
			}
			marketPriceFloat,err := strconv.ParseFloat(marketPriceAuto,64)
			if err != nil{
				logger.Error(err)
			}
			logic.StartSale(marketPriceFloat,percentFloat)
		}
	}
}