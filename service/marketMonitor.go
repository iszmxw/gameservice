/**
 @author:way
 @date:2021/11/27
 @note
**/

package main

import (
	"fmt"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logic"
	"redisData/pkg/logger"
	"redisData/setting"
	"strconv"
	"time"
)

func init() {
	// 定义日志目录
	logger.Init("marketMonitor")
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
	defer redis.Close()
}

func main() {
	for {
		var safe, _ = redis.GetData("safe")
		f, _ := strconv.ParseFloat(safe, 64)
		fmt.Printf("安全监控设置百分比为%f\n", f)
		//获取最新的市场数据,通过redis计算出来的
		newMarketPrice := logic.GetMarketDataByRedis()
		//读取旧数据，直接从redis中获取
		fmt.Println(newMarketPrice)
		msg := logic.RiskControl(newMarketPrice, f)
		fmt.Println(msg)
		//println(msg)
		time.Sleep(1 * time.Second)
	}

}
