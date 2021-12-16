/**
 @author:way
 @date:2021/12/15
 @note 存放baby脚本相关的逻辑
**/

package logic

import (
	"fmt"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/model"
	"redisData/pkg/logger"
	"strconv"
)

// CountBabyMarPrice 传入切片输出市场价
func CountBabyMarPrice(priceList []float64) []float64 {
	m1 := make(map[float64]int)
	var s2 []int
	var max int
	var s3 []float64

	// 统计频率最高的价格
	for _, v := range priceList {
		if m1[v] != 0 {
			m1[v]++
		} else {
			m1[v] = 1
		}
	}
	//遍历m1把里面的float转化成string
	for _, v := range m1 {
		//拼接成数组
		s2 = append(s2, v)
	}
	// 取出来放进数组

	for _, v := range m1 {
		//拼接成数组
		s2 = append(s2, v)
	}
	//算出最大值
	if s2 == nil{
		return nil
	}
	max = s2[0]
	for i := 0; i < len(s2); i++ {
		if max < s2[i] {
			max = s2[i]
		}
	}

	//存在出现同样次数的
	for k, v := range m1 {
		if v == max {
			s3 = append(s3, k)
		}
	}
	//插入一条redis数据，把这次遍历市场价占比计算后返回
		m2 := make(map[string]interface{})
		for i,v := range m1{
			str := strconv.FormatFloat(i, 'E', -1, 64)
			m2[fmt.Sprintf("%s",str)] = v
		}
		redis.CreatHashKey(fmt.Sprintf("baby:Proportion"),m2)
	return s3
}

//InitConfig 初始配置参数
func InitConfig()  {
	//买入参数设置
	babyBuyConfig := make(map[string]interface{}, 4)
	babyBuyConfig["percent"] = "10"
	babyBuyConfig["market_price"] = "10000"
	babyBuyConfig["status"] = "2" //1.打开 2.关闭
	babyBuyConfig["types"] = "1"  //买入固定为1
	redis.CreatHashKey("baby:ConfigBuy", babyBuyConfig)

	//卖出药水参数设置
	babySaleConfig := make(map[string]interface{}, 4)
	babySaleConfig["percent"] = "10"
	babySaleConfig["market_price"] = "10000"
	babySaleConfig["status"] = "2"
	babySaleConfig["types"] = "2"	 //买入固定为2
	redis.CreatHashKey("baby:ConfigSale", babySaleConfig)

	//设置买半自动总开关
	babyStopAutoBuy := make(map[string]interface{}, 2)
	babyStopAutoBuy["CrlName"] = "buy"
	babyStopAutoBuy["Super"] = "2"	//1.为打开 2.为关闭
	redis.CreatHashKey("baby:ConfigStopAutoBuy", babyStopAutoBuy)
	//设置卖半自动总开关
	babyStopAutoSale := make(map[string]interface{}, 2)
	babyStopAutoSale["CrlName"] = "sale"
	babyStopAutoSale["Super"] = "2"	//1.为打开 2.为关闭
	redis.CreatHashKey("baby:ConfigStopAutoSale", babyStopAutoSale)


	//设置药水的风控
	riskBabyFall:= make(map[string]interface{}, 5)
	riskBabyFall["OperationType"] = "1"  //1.为停止脚本 2.发送钉钉 3.停止脚本且发送钉钉
	riskBabyFall["Percentage"] = "10"
	riskBabyFall["Situation"] = "fall"
	riskBabyFall["Status"] = "2"	//1.为打开 2.为关闭
	riskBabyFall["TimeLevel"] = "60"
	redis.CreatHashKey("baby:ConfigRisk:fall", riskBabyFall)
	//设置药水的风控
	riskBabyRise:= make(map[string]interface{}, 5)
	riskBabyRise["OperationType"] = "1"  //1.为停止脚本 2.发送钉钉 3.停止脚本且发送钉钉
	riskBabyRise["Percentage"] = "10"
	riskBabyRise["Situation"] = "rise"
	riskBabyRise["Status"] = "2"	//1.为打开 2.为关闭
	riskBabyRise["TimeLevel"] = "60"
	redis.CreatHashKey("baby:ConfigRisk:rise", riskBabyRise)
	//设置卖出率参数
	babySaleRate:= make(map[string]interface{}, 4)
	babySaleRate["time_level"] = "1"  //1.为停止脚本 2.发送钉钉 3.停止脚本且发送钉钉
	babySaleRate["percent"] = "10"
	babySaleRate["status"] = "rise"
	babySaleRate["operation_type"] = "2"	//1.为打开 2.为关闭
	redis.CreatHashKey("baby:ConfigSaleRate", babySaleRate)
}


func StartBuy(marketPrice float64,percent float64)  {
	////获取市场价
	//marketPriceAuto, GetDataErr := redis.GetData("baby:marketPrice")
	//if GetDataErr != nil {
	//	logger.Info(GetDataErr)
	//	return
	//}
	//marketPriceFloat,err := strconv.ParseFloat(marketPriceAuto,64)
	//if err != nil{
	//	logger.Error(err)
	//}
	//遍历市场清单
	data := mysql.GetAllBabyTx()
	for _,v := range data{
		if redis.ExistEle("baby:order",v.Token){
			logger.Info("该订单已经入库")
			continue
		}
		productPrice,ParseFloatErr :=strconv.ParseFloat(v.Price,64)
		if ParseFloatErr != nil{
			logger.Error(ParseFloatErr)
		}
		//判断是否买入
		if marketPrice * 0.99 > (productPrice/1000000000000000000) * ((100+percent)*0.01){
			//添加买入清单
			var order model.BabyOrder
			order.MarketPrice = marketPrice
			order.Status = 1
			order.Name = "baby"
			order.TokenId= v.Token
			order.FixPrice = productPrice/1000000000000000000
			mysql.CreateOneOrder(order)
			//把买入的token存进redis避免重复购买
			redis.CreateSetData("baby:order",v.Token)
		}
	}
}

func StartSale(marketPrice float64,percent float64)  {
	//获取order表里面的全部数据
	orderList := mysql.GetAllOrder()
	//对比是否达到盈利点
	for _,v := range orderList{
		//状态为1是买入状态
		if v.Status== 1 {
			//查重避免二次出售
			if redis.ExistEle("baby:saleOrder", v.TokenId) {
				logger.Info("该订单已经在销售状态")
				continue
			}
			//查询市场价
			//marketPrice, GetDataErr := redis.GetData("baby:marketPrice")
			//if GetDataErr != nil {
			//	logger.Error(GetDataErr)
			//	continue
			//}
			//marketPriceFloat,err := strconv.ParseFloat(marketPrice,64)
			//if err != nil{
			//	logger.Info(err)
			//}
			//达到盈利点卖出
			if marketPrice*0.99 >= v.FixPrice*((100+percent)*0.01) {
				//执行更新操作，把status改成2
				var data model.BabyOrder
				data.Id = v.Id
				data.Status = 2
				data.Profit = marketPrice*0.99 - v.FixPrice
				data.SalePrice = marketPrice * 0.99
				data.TokenId = v.TokenId
				data.Name = v.Name
				mysql.UpdateOneOrder(data)
				//同时把这个token存进redis避免重复购买
				redis.CreateSetData("baby:saleOrder", v.TokenId)
			}
		}
	}
}