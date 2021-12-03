/**
 @author:way
 @date:2021/11/27
 @note
**/

package main

import (
	"fmt"
	"go.uber.org/zap"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logic"
	"redisData/pkg/logger"
	"redisData/setting"
	"strconv"
	"time"
)

func main() {

	//初始化viper
	if err := setting.Init(""); err != nil {
		logger.Error(err)
		return
	}
	defer zap.L().Sync() //把缓冲区的日志添加
	zap.L().Debug("init logger success")

	//初始化redis
	if err := redis.InitClient(); err != nil {
		logger.Error(err)
		return
	}
	defer redis.Close()

	//初始化MySQL
	mysql.InitMysql()
	//开始缓存
	for {
		var buy,_ = redis.GetData("buy")
		f, _ := strconv.ParseFloat(buy, 64)
		fmt.Printf("买入设置百分比为%f\n",f)
		//市场价直接从redis中取
		marketPrice, _ := redis.GetData("eggMarket")
		//数据转换
		float64MarketPrice,_  := strconv.ParseFloat(marketPrice, 64)
		//执行买入脚本
		logic.SetBuyALG(float64MarketPrice,f)
		fmt.Println("本轮购买完毕")
		time.Sleep(time.Second*10)
	}

}
