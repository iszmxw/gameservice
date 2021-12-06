/**
 @author:way
 @date:2021/11/28
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
	logger.Init("sale")
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
	defer redis.Close()
	//定义全局代表收益
	for {
		//调用卖出逻辑
		marketPriceKey :=  "Metamon Egg.MarketPrice"
		var sale, _ = redis.GetData("sale")
		f, _ := strconv.ParseFloat(sale, 64)
		fmt.Printf("买出设置百分比为%f\n", f)

		//从redis读取
		if !redis.ExistKey("income"){
			err := redis.CreateDurableKey("income",0)
			if err != nil {
				logger.Error(err)
				return 
			}
		}
		data, err := redis.GetData("income")

		if err != nil {
			logger.Error(err)
			return
		}
		account, _ := strconv.ParseFloat(data, 64)
		logger.Info(account)
		logic.SetSaleALG(marketPriceKey,account, f)
		time.Sleep(2 * time.Second)
	}
}
