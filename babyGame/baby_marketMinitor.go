/**
 @author:way
 @date:2021/12/16
 @note
**/



package main

import (
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/pkg/logger"
	"redisData/setting"
	"redisData/utils"
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
	//获取配置文件上涨下跌的全部参数
	rise := redis.GetHashDataAll("baby:ConfigRisk:rise")
	rise_operationType := rise["OperationType"]
	rise_percentage := rise["Percentage"]
	rise_percentage_float := utils.StringToFloat64(rise_percentage)
	rise_situation := rise["Situation"]
	rise_status := rise["Status"]
	rise_timeLevel := rise["TimeLevel"]
	rise_timeLevel_int := utils.StringToInt(rise_timeLevel)

	fall := redis.GetHashDataAll("baby:ConfigRisk:fall")
	fall_operationType :=  fall["OperationType"]
	fall_percentage := fall["Percentage"]
	rise_percentage_float := utils.StringToFloat64(rise_percentage)
	fall_situation := fall["Situation"]
	fall_status := fall["Status"]
	fall_timeLevel := fall["TimeLevel"]
	fall_timeLevel_int := utils.StringToInt(fall_timeLevel)

	//获取当前市场价格
	marketPrice, GErr := redis.GetData("baby:marketPrice")
	if GErr != nil {
		logger.Error(GErr)
		return 
	}

	//根据输入参数拿对应时间段的市场价
	//对比风控值
}