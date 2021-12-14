/**
 @author:way
 @date:2021/12/14
 @note 轮询 更新买入卖出状态
**/

package main

import (
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logic"
	"redisData/model"
	"redisData/pkg/logger"
	"redisData/setting"
	"redisData/utils"
)

func init() {
	// 定义日志目录
	logger.Init("client")
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
	//查询全部符合买卖的数据,
	data := mysql.GetBuyDataHashNotNull()
	for _,v := range *data{
		//通过接口查询交易状态
		data1 :=  logic.ReqGetTxStatus(v.TxHash)
		//查询交易详情
		data2 :=  logic.ReqTxDetailByHash(v.TxHash)
		buyAddress := ""
		//如果交易状态为1,商品已经卖出，没有抢到，或者已经抢到   加入买入清单得知道买入的是否是自己的地址值
		if data1.Result.Status == "1" && data2.Result.From == buyAddress{
			//买入成功
			//更新数据库里面的交易状态
			var d model.Buy
			d.Gid = v.Gid
			d.Type = 3
			mysql.UpdateBuy(d)
		}
		now := utils.GetNowTimeS()
		endUpdate := utils.DatetimeToTimestamp(v.UpdatedAt.String())
		if data1.Status=="" && now - endUpdate > 5 * 3600{
			//买入成功
			//更新数据库里面的交易状态
			var d model.Buy
			d.Gid = v.Gid
			d.Type = 5 // 5是废弃状态,买失败和卖失败
			mysql.UpdateBuy(d)
		}
	}

	//查询全部符合买卖的数据
	SaleData := mysql.GetSaleDataHashNotNull()
	for _,v := range *SaleData{
		//通过接口查询交易状态
		data1 :=  logic.ReqGetTxStatus(v.TxHash)
		//查询交易详情
		data2 :=  logic.ReqTxDetailByHash(v.TxHash)
		buyAddress := ""
		if data1.Result.Status == "1" && data2.Result.From == buyAddress{
			//买入成功
			//更新数据库里面的交易状态
			var d model.Buy
			d.Gid = v.Gid
			d.Type = 4
			mysql.UpdateBuy(d)
		}
		now := utils.GetNowTimeS()
		endUpdate := utils.DatetimeToTimestamp(v.UpdatedAt.String())
		if data1.Status=="" && now - endUpdate > 5 * 3600{
			//买入成功
			//更新数据库里面的交易状态
			var d model.Buy
			d.Gid = v.Gid
			d.Type = 5 // 5是废弃状态
			mysql.UpdateBuy(d)
		}
	}
}

