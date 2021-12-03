/**
 @author:way
 @date:2021/11/28
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
	"strconv"
	"time"
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
	defer redis.Close()
	//定义全局代表收益
	var account float64
	for {
		//调用卖出逻辑
		var sale, _ = redis.GetData("sale")
		f, _ := strconv.ParseFloat(sale, 64)
		fmt.Printf("买出设置百分比为%f\n", f)
		account = logic.SetSaleALG(account, f)
		err := redis.CreateDurableKey("income", account)
		fmt.Println(account)
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(10 * time.Second)
	}
}
