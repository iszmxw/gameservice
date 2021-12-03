/**
 @author:way
 @date:2021/11/27
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
	"strconv"
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

		var safe,_ = redis.GetData("safe")
		f, _ := strconv.ParseFloat(safe, 64)
		fmt.Printf("安全监控设置百分比为%f\n",f)
		//获取最新的市场数据,通过redis计算出来的
		newMarketPrice := logic.GetMarketDataByRedis()
		//读取旧数据，直接从redis中获取
		fmt.Println(newMarketPrice)
		msg := logic.RiskControl(newMarketPrice,f)
		fmt.Println(msg)
		//println(msg)
		time.Sleep(1*time.Second)
	}



}
