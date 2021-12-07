/**
@author:way
@date:2021/12/1
@note
**/

package main

import (
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logic"
	"redisData/pkg/logger"
	"redisData/setting"
)

func init() {
	// 定义日志目录
	logger.Init("marketPriceOnline")
	// 初始化 viper 配置
	if err := setting.Init(""); err != nil {
		logger.Info("viper init fail")
		logger.Error(err)
		return
	}
	mysql.InitMysql()
	//初始化redis
	if err := redis.InitClient(); err != nil {
		logger.Error(err)
		return
	}

}

func main() {
	defer redis.Close()
	for {
		//1.请求数据 改成redis
		//key := "Metamon Egg.List"
		key := "Potion.List"
		logic.SetMarketPrice(key)
		key = "Metamon Egg.List"
		logic.SetMarketPrice(key)

	}
}
