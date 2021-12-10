/**
 @author:way
 @date:2021/12/9
 @note
**/

package main

import (
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logic"
	"redisData/pkg/email"
	"redisData/pkg/logger"
	"redisData/setting"
	"redisData/utils"
	"strconv"
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

	//获取当前时间戳
	time := utils.GetNowTimeS()
	//获取redis中的配置参数
	sellingRateSetting := redis.GetHashDataAll("SellingRate")
	status := sellingRateSetting["status"]
	for {
		if status == "2" {
			logger.Info("脚本没有打开")
			continue
		}
		//计算时间
		time2 := sellingRateSetting["time_level"]
		time3, err := strconv.ParseInt(time2, 10, 64)
		time = time - time3*1800
		//转化成字符串
		timestr := utils.TimestampToDatetime(time)
		if err != nil {
			logger.Error(err)
			return
		}
		//传进数据库
		//获取买的的count[]
		buyCount := mysql.GetBuySaleCountList(timestr, 1)
		//获取卖的的count[]
		saleCount := mysql.GetBuySaleCountList(timestr, 2)
		var buy float64
		//从数据库获取买出的全部数量
		for _, v := range buyCount {
			buy = buy + v
		}
		logger.Info(buy)
		var sale float64
		//从数据获取全部买出的数量
		for _, v := range saleCount {
			sale = sale + v
		}
		logger.Info(sale)
		//把剩余的存在redis中的参数读出来

		operation_type := sellingRateSetting["operation_type"]

		percent := sellingRateSetting["percent"]
		//转化成浮点
		f, err := strconv.ParseFloat(percent, 64)
		if err != nil {
			logger.Error(err)
			return
		}
		//开始监控
		//判断status判断脚本是否打开
		logger.Info(sale / buy)
		logger.Info(f * 0.01)
		if (sale / buy) < f*0.01 {
			switch operation_type {
			case "1":
				logger.Info("停止脚本")
				logic.StopScript()
			case "2":
				logger.Info("发送钉钉")
				email.SendDingMsg("SellingRate","卖出率过低")
			case "3":
				logger.Info("停止脚本，且发送钉钉")
				//获取买卖脚本的配置文件，修改其中的状态
				logic.StopScript()
				//发送钉钉
			}
		}
	}

}
