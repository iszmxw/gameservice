/**
 @author:way
 @date:2021/12/1
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
	logger.Init("dataHandle")
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
	defer redis.Close()
	for {
		str := redis.RmListHead("assertList")
		if len(str) == 0 {
			logger.Info("kongshuju")
			continue
		}
		go logic.StoreListToMysql(str)
	}
}
