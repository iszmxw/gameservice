/**
 @author:way
 @date:2021/12/13
 @note
**/

package main

import (
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logic"
	"redisData/pkg/logger"
	"redisData/setting"
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
	//redis.Close()
	address := "0x1F6bc601fDe821E0079c89529c79e3C616Da7E22" //市场合约地址
	apikey := "DUKNN1QZMITSSZC61YINTD1CWQ92FWEKHM" 	//api 请求token
	sort := "desc" //排序
	offset := "100" //每次请求数量
	page := "1"  //页数
	contain := "0xc37dfc5b"  //0xc37dfc5b 是购买的方法
	logic.ReqBNTxList(address,apikey,sort,offset,page,contain)
}


