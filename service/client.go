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

func StartList(pageSize int, category int) {
	//缓存鸡蛋市场数据 每次拿100条最新 category 17 为鸡蛋，数据已经序列化成结构体
	data := logic.RequestAssertsData(pageSize, category)
	if data == nil {
		logger.Info("请求游戏市场链接返回信息为空")
		//可能是访问频繁的原因休息1min后继续访问
		time.Sleep(1 * time.Minute)
		return
	}
	//把这个数据存进redis,
	key := fmt.Sprintf("%s.List", data.List[0].Name)
	//把内容序列化成字符串
	marshalData, Merr := json.Marshal(data)
	if Merr != nil {
		logger.Info("序列化数据失败")
		logger.Error(Merr)
		return
	}
	//创建对应的key
	Cerr := redis.CreateDurableKey(key, string(marshalData))
	if Cerr != nil {
		logger.Info("创建key失败")
		logger.Error(Cerr)
		return
	}
	time.Sleep(1 * time.Second)
}

func StartConfig() {
	//买入药水参数设置
	buySet15 := make(map[string]interface{}, 4)
	buySet15["product_id"] = "15"
	buySet15["percent"] = "10"
	buySet15["status"] = "2" //1.打开 2.关闭
	buySet15["types"] = "1"  //买入固定为1
	redis.CreatHashKey("BuySet:15", buySet15)
	//买入元兽蛋参数设置
	buySet17 := make(map[string]interface{}, 4)
	buySet17["product_id"] = "17"
	buySet17["percent"] = "10"
	buySet17["status"] = "2" //1.打开 2.关闭
	buySet17["types"] = "1"  //买入固定为1
	redis.CreatHashKey("BuySet:17", buySet17)

	//卖出药水参数设置
	SaleSet15 := make(map[string]interface{}, 4)
	SaleSet15["product_id"] = "15"
	SaleSet15["percent"] = "10"
	SaleSet15["status"] = "2"
	SaleSet15["types"] = "2"	 //买入固定为2
	redis.CreatHashKey("SaleSet:15", SaleSet15)
	//买入元兽蛋参数设置
	SaleSet17 := make(map[string]interface{}, 4)
	SaleSet17["product_id"] = "17"
	SaleSet17["percent"] = "10"
	SaleSet17["status"] = "2"   //默认值是关闭
	SaleSet17["types"] = "2"	//买入固定为2
	redis.CreatHashKey("SaleSet:17", SaleSet17)

	//设置买总开关
	buyAndSale_buy := make(map[string]interface{}, 2)
	buyAndSale_buy["CrlName"] = "buy"
	buyAndSale_buy["Super"] = "2"	//1.为打开 2.为关闭
	redis.CreatHashKey("buyAndSale:buy", buyAndSale_buy)
	//设置卖总开关
	buyAndSale_sale := make(map[string]interface{}, 2)
	buyAndSale_sale["CrlName"] = "sale"
	buyAndSale_sale["Super"] = "2"	//1.为打开 2.为关闭
	redis.CreatHashKey("buyAndSale:sale", buyAndSale_sale)

	//设置元兽蛋下跌的风控
	riskFall := make(map[string]interface{}, 5)
	riskFall["OperationType"] = "1"  //1.为停止脚本 2.发送钉钉 3.停止脚本且发送钉钉
	riskFall["Percentage"] = "10"
	riskFall["Situation"] = "fall"
	riskFall["Status"] = "2"	//1.为打开 2.为关闭
	riskFall["TimeLevel"] = "60"
	redis.CreatHashKey("risk:fall", riskFall)
	//设置元兽蛋上涨的风控
	riskRise := make(map[string]interface{}, 5)
	riskRise["OperationType"] = "1"  //1.为停止脚本 2.发送钉钉 3.停止脚本且发送钉钉
	riskRise["Percentage"] = "10"
	riskRise["Situation"] = "rise"
	riskRise["Status"] = "2"	//1.为打开 2.为关闭
	riskRise["TimeLevel"] = "60"
	redis.CreatHashKey("risk:rise", riskRise)
	//设置药水的风控
	riskPotionFall:= make(map[string]interface{}, 5)
	riskPotionFall["OperationType"] = "1"  //1.为停止脚本 2.发送钉钉 3.停止脚本且发送钉钉
	riskPotionFall["Percentage"] = "10"
	riskPotionFall["Situation"] = "fall"
	riskPotionFall["Status"] = "2"	//1.为打开 2.为关闭
	riskPotionFall["TimeLevel"] = "60"
	redis.CreatHashKey("risk:potion:fall", riskPotionFall)
	//设置药水的风控
	riskPotionRise:= make(map[string]interface{}, 5)
	riskPotionRise["OperationType"] = "1"  //1.为停止脚本 2.发送钉钉 3.停止脚本且发送钉钉
	riskPotionRise["Percentage"] = "10"
	riskPotionRise["Situation"] = "rise"
	riskPotionRise["Status"] = "2"	//1.为打开 2.为关闭
	riskPotionRise["TimeLevel"] = "60"
	redis.CreatHashKey("risk:potion:rise", riskPotionRise)
	//设置卖出率参数
	sellingRate:= make(map[string]interface{}, 4)
	riskPotionRise["time_level"] = "1"  //1.为停止脚本 2.发送钉钉 3.停止脚本且发送钉钉
	riskPotionRise["percent"] = "10"
	riskPotionRise["status"] = "rise"
	riskPotionRise["operation_type"] = "2"	//1.为打开 2.为关闭
	redis.CreatHashKey("SellingRate", sellingRate)
}
// 获取列表数据到redis
func main() {
	defer redis.Close()
	//初始化配置参数
	StartConfig()
	logger.Info("初始化启动参数")
	pageSize := 200
	//开始缓存
	for {
		//category := 17
		category := 15
		StartList(pageSize, category)
		category = 17
		StartList(pageSize, category)
	}
}
