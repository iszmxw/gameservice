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
	percentage := -1.0
	//获取order表里面的全部数据
	orderList := mysql.GetAllOrder()
	//对比是否达到盈利点
	for _,v := range orderList{
		//状态为1是买入状态
		if v.Status== 1{
			//查重避免二次出售
			if redis.ExistEle("baby:saleOrder",v.TokenId){
				logger.Info("该订单已经在销售状态")
				continue
			}
			//查询市场价
			marketPrice, GetDataErr := redis.GetData("baby:marketPrice")
			if GetDataErr != nil {
				logger.Error(GetDataErr)
				continue
			}
			marketPriceFloat,err := strconv.ParseFloat(marketPrice,64)
			if err != nil{
				logger.Info(err)
			}
			//达到盈利点卖出
			if marketPriceFloat * 0.99 >= v.FixPrice * ((100+percentage)*0.01){
				//执行更新操作，把status改成2
				var data model.BabyOrder
				data.Id= v.Id
				data.Status = 2
				data.Profit = marketPriceFloat * 0.99 - v.FixPrice
				data.SalePrice = marketPriceFloat * 0.99
				data.TokenId = v.TokenId
				data.Name = v.Name
				mysql.UpdateOneOrder(data)
				//同时把这个token存进redis避免重复购买
				redis.CreateSetData("baby:saleOrder",v.TokenId)
			}
		}
	}

}