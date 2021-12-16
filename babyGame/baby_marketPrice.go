/**
 @author:way
 @date:2021/12/16
 @note
**/

package main

import (
	"errors"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logic"
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
	//从mysql获取市场数据计算
	price := mysql.GetPriceList()
	if price ==nil{
		logger.Error(errors.New("price切片为空"))
	}
	//将string切片转换成float切片
	priceFoalt := make([]float64,0,len(price))
	for _,v := range price{
		float,err := strconv.ParseFloat(v,64)
		if err != nil{
			logger.Error(errors.New("string转换float失败"))
		}
		priceFoalt = append(priceFoalt,float/1000000000000000000)
	}
	//统计切片计算市场价
	marPrice := logic.CountBabyMarPrice(priceFoalt)
	if len(marPrice) > 1{
		logger.Info("同时出现两个市场价")
	}
	//写进mysql和redis
	CreateDurableKeyErr := redis.CreateDurableKey("baby:marketPrice", marPrice[0])
	if CreateDurableKeyErr != nil {
		logger.Error(CreateDurableKeyErr)
	}
	data := model.BabyMarketPrice{}
	data.MarketData = marPrice[0]
	data.MarketName = "baby"
	mysql.InsertBabyMarketPrice(data)
}