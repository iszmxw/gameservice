/**
 @author:way
 @date:2021/11/26
 @note
**/

package main

import (
	"fmt"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logic"
	"redisData/pkg/logger"
	"redisData/setting"
	"time"
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
	defer redis.Close()
}

func main() {
	//开始缓存
	for {
		//创建egg:id
		//logic.CreatEggData()
		//fmt.Println(111)
		//创建egg:id
		//logic.CreatPotionData()
		//fmt.Println(222)
		//创建 eggDataList
		err := logic.SetDataInRedis()
		if err != nil {
			logger.Error(err)
			return
		}
		logic.SetEggMarketPrice()
		fmt.Println("redis数据已更新")
		time.Sleep(time.Second * 30)

	}

}
