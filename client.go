/**
 @author:way
 @date:2021/11/26
 @note
**/

package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logger"
	"redisData/logic"
	"redisData/setting"
	"time"
)

func main() {

	//初始化viper
	if err := setting.Init(""); err != nil {
		zap.L().Error("viper init fail", zap.Error(err))
		return
	}
	//初始化日志
	if err := logger.InitLogger(viper.GetString("mode")); err != nil {
		zap.L().Error("init logger fail err", zap.Error(err))
		return
	}
	defer zap.L().Sync() //把缓冲区的日志添加
	zap.L().Debug("init logger success")

	//初始化MySQL
	mysql.InitMysql()

	//初始化redis
	if err := redis.InitClient(); err != nil {
		zap.L().Error("init redis fail err", zap.Error(err))
		return
	}
	defer redis.Close()

	//开始缓存
	for {
		//创建egg:id
		logic.CreatEggData()
		//fmt.Println(111)
		//创建egg:id
		//logic.CreatPotionData()
		//fmt.Println(222)
		//创建 eggDataList
		logic.SetDataInRedis()
		logic.SetEggMarketPrice()
		fmt.Println("redis数据已更新")
		time.Sleep(time.Second*30)

	}

}