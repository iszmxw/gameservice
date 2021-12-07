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
		logger.Error(err)
		return
	}
	mysql.InitMysql()
	// 初始化redis
	if err := redis.InitClient(); err != nil {
		logger.Error(err)
		return
	}
}
func startMonitor(timeLevel int,assets string) {
	dataOneMin := mysql.GetHistoryMarketData(timeLevel, assets)
	assetsListKey := fmt.Sprintf("%s.List",assets)
	if dataOneMin == nil {
		logger.Info("mysql获取参数错误")
		return
	}
	var safe, _ = redis.GetData("safe")
	f, _ := strconv.ParseFloat(safe, 64)
	//从mysql里面拿
	fmt.Printf("安全监控设置百分比为%f\n", f)
	//获取最新的市场数据,通过redis计算出来的
	newMarketPrice := logic.GetMarketDataByRedis(assetsListKey)
	//读取旧数据，直接从redis中获取
	fmt.Println(dataOneMin.MarketData)
	fmt.Println(newMarketPrice)
	msg := logic.RiskControl(newMarketPrice, dataOneMin.MarketData, f)
	fmt.Println(msg)
	time.Sleep(1 * time.Second)
}

func main() {
	defer redis.Close()
	for {

		startMonitor(60,"Metamon Egg")

	}
}
