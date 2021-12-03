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
	defer redis.Close()

}

func main() {
	//获取10s前的市场价格
	//获取30s前的市场价格
	//获取60秒前的市场价格
	//获取300秒前的市场价格 5m
	//获取900秒的市场价格  15m
	//获取1800秒的市场数据 30m
	//获取3600秒市场数据 1h
	//获取14400秒数据 4h
	//获取86400秒数据 1D
	//获取604800秒数据 1W
	//获取2592000秒的数据 1Mon
	for {
		dataOneMin := mysql.GetHistoryMarketData(60, "Metamon Egg")
		if dataOneMin == nil {
			logger.Info("mysql获取参数错误")
			return
		}
		var safe, _ = redis.GetData("safe")
		f, _ := strconv.ParseFloat(safe, 64)
		//从mysql里面拿
		fmt.Printf("安全监控设置百分比为%f\n", f)
		//获取最新的市场数据,通过redis计算出来的
		newMarketPrice := logic.GetMarketDataByRedis()
		//读取旧数据，直接从redis中获取
		fmt.Println(dataOneMin.MarketData)
		fmt.Println(newMarketPrice)
		msg := logic.RiskControl(newMarketPrice, dataOneMin.MarketData, f)
		fmt.Println(msg)
		//println(msg)
		time.Sleep(1 * time.Second)
	}
}
