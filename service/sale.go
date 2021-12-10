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
		float,err := strconv.ParseFloat(percent,64)
		if err != nil{
			logger.Error(err)
			return
		}
		sale = float
	}
	if strings.Contains(key,"Potion"){
		data1 := redis.GetHashDataAll("SaleSet:15")
		percent := data1["percent"]
		float,err := strconv.ParseFloat(percent,64)
		if err != nil{
			logger.Error(err)
			return
		}
		sale = float
	}
	//调用卖出逻辑
	//marketPriceKey :=  "Metamon Egg.MarketPrice"
	//marketPriceKey := "Potion.MarketPrice"
	marketPriceKey := key
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
func startSale2(key string){
	var sale float64  //设置卖出%百分比
	if strings.Contains(key,"Egg"){
		data1 := redis.GetHashDataAll("SaleSet:17")
		percent := data1["percent"]
		float,err := strconv.ParseFloat(percent,64)
		if err != nil{
			logger.Error(err)
			return
		}
		sale = float
	}
	if strings.Contains(key,"Potion"){
		data1 := redis.GetHashDataAll("SaleSet:15")
		percent := data1["percent"]
		float,err := strconv.ParseFloat(percent,64)
		if err != nil{
			logger.Error(err)
			return
		}
		sale = float
	}
	//调用卖出逻辑
	//marketPriceKey :=  "Metamon Egg.MarketPrice"
	//marketPriceKey := "Potion.MarketPrice"
	marketPriceKey := key
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
	logic.SetSaleALG2(marketPriceKey,account, sale)
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
		//同为1,设定自己买卖的脚本
		logger.Info(all["Super"])
		logger.Info(data1["status"])
		logger.Info(data2["status"])
		if data1["status"] == "1" && all["Super"] == "1"{
			logger.Info("执行半自动脚本")
			marketPriceKey :=  "Metamon Egg.MarketPrice"
			startSale2(marketPriceKey)
		}
		//为1和2设置自动脚本
		if data1["status"] == "1" && all["Super"] == "2"{
			logger.Info("执行全自动自动脚本")
			marketPriceKey :=  "Metamon Egg.MarketPrice"
			startSale(marketPriceKey)
		}

		if data2["status"] == "1" && all["Super"] == "1"{
			logger.Info("执行半自动脚本")
			marketPriceKey := "Potion.MarketPrice"
			startSale2(marketPriceKey)
		}
		//为1和2设置自动脚本
		if data2["status"] == "1" && all["Super"] == "2"{
			logger.Info("执行全自动自动脚本")
			marketPriceKey := "Potion.MarketPrice"
			startSale(marketPriceKey)
		}
		//logger.Info("没有符合条件")
	}
}
