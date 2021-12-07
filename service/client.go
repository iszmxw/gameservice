/**
 @author:way
 @date:2021/11/26
 @note
**/

package main

import (
	"encoding/json"
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

}

func StartList(pageSize int,category int)  {
	//缓存鸡蛋市场数据 每次拿100条最新 category 17 为鸡蛋，数据已经序列化成结构体
	data := logic.RequestAssertsData(pageSize, category)
	if data == nil{
		logger.Info("请求游戏市场链接返回信息为空")
		//可能是访问频繁的原因休息1min后继续访问
		time.Sleep(1*time.Minute)
		return
	}
	//把这个数据存进redis,
	key := fmt.Sprintf("%s.List",data.List[0].Name)
	//把内容序列化成字符串
	marshalData, Merr := json.Marshal(data)
	if Merr != nil {
		logger.Info("序列化数据失败")
		logger.Error(Merr)
		return
	}
	//创建对应的key
	Cerr := redis.CreateDurableKey(key,string(marshalData))
	if Cerr != nil {
		logger.Info("创建key失败")
		logger.Error(Cerr)
		return
	}
	time.Sleep(1*time.Second)
}

// 获取列表数据到redis
func main() {
	defer redis.Close()
	pageSize := 100
	//开始缓存
	for {
		//category := 17
		category := 15
		StartList(pageSize,category)
		category = 17
		StartList(pageSize,category)
	}
}
