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
	address := "0x53e562b9b7e5e94b81f10e96ee70ad06df3d2657"
	apikey := "DUKNN1QZMITSSZC61YINTD1CWQ92FWEKHM"
	sort := "desc"
	offset := "10"
	page := "1"
	contain := "0x"
	logic.ReqBNTxList(address,apikey,sort,offset,page,contain)
}


