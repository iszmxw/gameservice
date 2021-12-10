package main

import (
	"fmt"
	"github.com/spf13/viper"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/pkg/logger"
	"redisData/routes"
	"redisData/setting"
)

func init() {
	// 定义日志目录
	logger.Init("redisData")
	// 初始化 viper 配置
	if err := setting.Init(""); err != nil {
		logger.Info("viper init fail")
		logger.Error(err)
		return
	}
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
	// 初始化routes
	r := routes.SetUp()
	//err := r.RunTLS(fmt.Sprintf(":%d", viper.GetInt("port")), "./conf/ssl.pem", "./conf/ssl.key")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	_ = r.Run(fmt.Sprintf(":%d", viper.GetInt("port")))
}
