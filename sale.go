/**
 @author:way
 @date:2021/11/28
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
	"strconv"
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

	//初始化redis
	if err := redis.InitClient(); err != nil {
		zap.L().Error("init redis fail err", zap.Error(err))
		return
	}
	defer redis.Close()
	//mysql
	//初始化MySQL
	mysql.InitMysql()

	//定义全局代表收益
	var account float64

	for  {
		//调用卖出逻辑
		var sale,_ = redis.GetData("sale")
		f, _ := strconv.ParseFloat(sale, 64)
		fmt.Printf("买出设置百分比为%f\n",f)
		account = logic.SetSaleALG(account,f)
		err := redis.CreateDurableKey("income", account)
		fmt.Println(account)
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(10*time.Second)
	}



}
