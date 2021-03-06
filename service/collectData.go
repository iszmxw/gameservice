/**
 @author:way
 @date:2021/12/1
 @note
**/

package main

import (
	"encoding/json"
	"fmt"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logic"
	"redisData/model"
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

func startCollect(key string)  {
	redisdata, err := redis.GetData(key)
	if err != nil {
		logger.Error(err)
		return
	}
	if len(redisdata) < 150{
		return
	}
	//计算出市场价格存进redis和mysql
	// 序列化返回的结果
	var data model.ResponseDataList
	if Uerr := json.Unmarshal([]byte(redisdata), &data); Uerr != nil {
		logger.Error(Uerr)
	}
	logic.ManageData(&data)
	fmt.Println("处理数据完毕")
	//3.持久化到mysql
	//fmt.Println("持久化到mysql完毕")
	//for  {
	//	go logic.StoreListToMysql()
	//}
	time.Sleep(1 * time.Second)
}

func main() {
	defer redis.Close()
	for {

		key := "Potion.List"
		startCollect(key)
		key = "Metamon Egg.List"
		startCollect(key)
	}
}
