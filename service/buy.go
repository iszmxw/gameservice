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
	"strings"
	"time"
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
func startBuy(listKey string, marketPrice string) {
	var buy float64
	//添加类型判断
	if strings.Contains(listKey, "Egg") {
		data1 := redis.GetHashDataAll("BuySet:17")
		percent := data1["percent"]
		float, _ := strconv.ParseFloat(percent, 64)
		buy = float
	}
	if strings.Contains(listKey, "Potion") {
		data1 := redis.GetHashDataAll("BuySet:15")
		percent := data1["percent"]
		float, _ := strconv.ParseFloat(percent, 64)
		buy = float
	}
	//f, _ := strconv.ParseFloat(buy, 64)
	//fmt.Printf("买入设置百分比为%f\n", f)
	//市场价直接从redis中取
	float64MarketPrice, _ := strconv.ParseFloat(marketPrice, 64)
	//执行买入脚本
	logic.SetBuyALG(listKey, float64MarketPrice, buy)
	fmt.Println("本轮购买完毕")
	time.Sleep(time.Second * 1)
}
func main() {
	defer redis.Close()
	//开始缓存
	for {
		all := redis.GetHashDataAll("buyAndSale:buy")
		//判断总开关状态
		data1 := redis.GetHashDataAll("BuySet:15")
		logger.Info(data1["status"])
		logger.Info( all["Super"])
		if data1["status"] == "1" && all["Super"] == "2" {
			logger.Info("执行全自动")
			logger.Info("potionList执行买入脚本")
			key := "Potion.List"
			marketPrice, _ := redis.GetData("Potion.MarketPrice")
			logger.Info(marketPrice)
			startBuy(key, marketPrice)
		}

		if data1["status"] == "1" && all["Super"] == "1" {
			logger.Info("执行半自动")
			logger.Info("potionList执行买入脚本")
			key := "Potion.List"
			//从redis里面存的市场价
			marketPrice := data1["market_price"]
			logger.Info(marketPrice)
			startBuy(key, marketPrice)
		}

		data2 := redis.GetHashDataAll("BuySet:17")
		logger.Info(data2["status"])
		logger.Info(data2["Super"])
		if data2["status"] == "1" && all["Super"] == "2" {
			logger.Info("Egg执行买入脚本")
			logger.Info("执行全自动")
			//判断物品开关
			key := "Metamon Egg.List"
			marketPrice, _ := redis.GetData("Metamon Egg.MarketPrice")
			startBuy(key, marketPrice)
		}
		if data2["status"] == "1" && all["Super"] == "1" {
			logger.Info("Egg执行买入脚本")
			logger.Info("执行半自动")
			//判断物品开关
			key := "Metamon Egg.List"
			marketPrice:= data2["market_price"]
			logger.Info(marketPrice)
			startBuy(key, marketPrice)
		}
		logger.Info("没有符合条件")
	}
}
