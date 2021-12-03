/**
 @author:way
 @date:2021/12/1
 @note
**/

package main

import (
	"go.uber.org/zap"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logic"
	"redisData/pkg/logger"
	"redisData/setting"
)

func main() {
	//初始化viper
	if err := setting.Init(""); err != nil {
		zap.L().Error("viper init fail", zap.Error(err))
		return
	}

	//初始化MySQL
	mysql.InitMysql()

	//初始化redis
	if err := redis.InitClient(); err != nil {
		zap.L().Error("init redis fail err", zap.Error(err))
		return
	}
	defer redis.Close()

	for  {
		str := redis.RmListHead("assertList")
		if len(str) == 0{
			logger.Info("kongshuju")
			continue
		}
		go logic.StoreListToMysql(str)
	}
}