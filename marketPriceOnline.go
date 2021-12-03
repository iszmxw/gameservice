/**
 @author:way
 @date:2021/12/1
 @note
**/

package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logger"
	"redisData/logic"
	"redisData/setting"
	"time"
)

func main() {


	//初始化viper
	if err := setting.Init(""); err != nil {
		zap.L().Error("viper init fail", zap.Error(err))
		return
	}
	//初始化日志
	if err := logger.InitLogger(viper.GetString("mode")); err != nil {
		zap.L().Error("init logger fail err", zap.Error(err))
		return
	}
	defer zap.L().Sync() //把缓冲区的日志添加
	zap.L().Debug("init logger success")

	//初始化MySQL
	mysql.InitMysql()

	//初始化redis
	if err := redis.InitClient(); err != nil {
		zap.L().Error("init redis fail err", zap.Error(err))
		return
	}
	defer redis.Close()

	for  {
		//1.请求数据
		pageSize := 300
		category := 17
		data := logic.GetAssertsData(pageSize, category)
		fmt.Println("拿到数据")

		//计算市场价然后返回
		logic.SetMarketPriceOnline(data)
		fmt.Println("算出市场价")
		time.Sleep(2*time.Second)
	}


}