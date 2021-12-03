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

func main() {
	//初始化viper
	if err := setting.Init(""); err != nil {
		logger.Error(err)
		return
	}

	//初始化redis
	if err := redis.InitClient(); err != nil {
		logger.Error(err)
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
