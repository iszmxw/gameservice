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
	//获取当前市场价格
	//根据输入参数拿对应时间段的市场价
	//对比风控值
}