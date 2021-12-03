/**
 @author:way
 @date:2021/12/1
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
	logger.Init("collectData")
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
		//1.请求数据
		pageSize := 300
		category := 17
		data := logic.GetAssertsData(pageSize, category)
		fmt.Println("拿到数据")
		//2.存redis
		logic.ManageData(data)
		fmt.Println("处理数据完毕")
		//3.持久化到mysql
		//fmt.Println("持久化到mysql完毕")
		//for  {
		//	go logic.StoreListToMysql()
		//}
		time.Sleep(30 * time.Second)
	}
}
