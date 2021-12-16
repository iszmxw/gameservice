/**
 @author:way
 @date:2021/12/16
 @note
**/

package main

import (
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/model"
	"redisData/pkg/logger"
	"redisData/setting"
	"strconv"
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

func main() {
	percentage := 1.0
	//获取市场价
	marketPrice, GetDataErr := redis.GetData("baby:marketPrice")
	if GetDataErr != nil {
		logger.Info(GetDataErr)
		return
	}
	marketPriceFoalt,err := strconv.ParseFloat(marketPrice,64)
	if err != nil{
		logger.Info(err)
	}
	//遍历市场清单
	data := mysql.GetAllBabyTx()
	for _,v := range data{
		if redis.ExistEle("baby:order",v.Token){
			logger.Info("该订单已经入库")
			continue
		}
		productPrice,ParseFloatErr :=strconv.ParseFloat(v.Price,64)
		if ParseFloatErr != nil{
			logger.Error(ParseFloatErr)
		}
		//判断是否买入
		if marketPriceFoalt * 0.99 > (productPrice/1000000000000000000) * ((100+percentage)*0.01){
			//添加买入清单
			var order model.BabyOrder
			order.MarketPrice = marketPriceFoalt
			order.Status = 1
			order.Name = "baby"
			order.TokenId= v.Token
			order.FixPrice = productPrice/1000000000000000000
			mysql.CreateOneOrder(order)
			//把买入的token存进redis避免重复购买
			redis.CreateSetData("baby:order",v.Token)
		}
	}
}