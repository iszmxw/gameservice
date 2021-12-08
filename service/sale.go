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
	"strings"
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

func startSale(key string){
	var sale float64  //设置卖出%百分比
	if strings.Contains(key,"Egg"){
		data1 := redis.GetHashDataAll("SaleSet:17")
		percent := data1["percent"]
		float,_ := strconv.ParseFloat(percent,64)
		sale = float
	}

	if strings.Contains(key,"Potion"){
		data1 := redis.GetHashDataAll("SaleSet:15")
		percent := data1["percent"]
		float,_ := strconv.ParseFloat(percent,64)
		sale = float
	}

	//调用卖出逻辑
	//marketPriceKey :=  "Metamon Egg.MarketPrice"
	//marketPriceKey := "Potion.MarketPrice"
	marketPriceKey := key
	//var sale, _ = redis.GetData("sale")
	//f, _ := strconv.ParseFloat(sale, 64)
	fmt.Printf("买出设置百分比为%f\n", sale)
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
	logic.SetSaleALG(marketPriceKey,account, sale)
	time.Sleep(2 * time.Second)
}

func main() {
	defer redis.Close()
	//定义全局代表收益
	for {
		//判断总开关状态
		all := redis.GetHashDataAll("buyAndSale:sale")
		data1 := redis.GetHashDataAll("SaleSet:15")
		data2 := redis.GetHashDataAll("SaleSet:17")
		//判断物品开关
		if data1["status"] == "1" && all["Super"] == "1"{
			marketPriceKey :=  "Metamon Egg.MarketPrice"
			startSale(marketPriceKey)
		}

		if data2["status"] == "1" && all["Super"] == "1"{
			marketPriceKey := "Potion.MarketPrice"
			startSale(marketPriceKey)
		}

	}
}
