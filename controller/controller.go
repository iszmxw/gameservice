/**
 @author:way
 @date:2021/12/3
 @note
**/

package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logic"
	"redisData/model"
	"redisData/pkg/logger"
	"redisData/utils"
	"strconv"
)

// GetDataHandle 返回id详情列表前端用
func GetDataHandle(c *gin.Context) {
	//获取参数
	var p model.ParamTypeId
	BErr := c.Bind(&p)
	if BErr != nil {
		logger.Error(BErr)
		return
	}
	//如果typeId =0 返回全部  如果不为零就返回根据id查询数据库资产名称后查询返回
	if p.TypeId == 0 {
		logger.Info(p.TypeId)
		logger.Info(0000000000)
		data := mysql.GetAssetDetail100()
		//返回数据
		c.JSON(200, gin.H{
			"msg":  "ok",
			"code": 200,
			"data": data,
		})
		return
	}
	if p.TypeId == 17 || p.TypeId == 15 {
		logger.Info(p.TypeId)
		//逻辑处理 1.根据ID找到前缀
		asset := mysql.GetAssetName(p.TypeId)
		data, err := logic.GetKeysByPfx(asset.TypeName)
		if err != nil {
			logger.Error(err)
			return
		}
		//返回数据
		c.JSON(200, gin.H{
			"msg":  "ok",
			"code": 200,
			"data": data,
		})
	}
}

//GetMarketPriceHandle 获取市场价格
func GetMarketPriceHandle(c *gin.Context) {
	//获取参数
	//var p model.ParamTypeId
	//Berr := c.Bind(&p)
	//if Berr != nil {
	//	logger.Info(Berr)
	//	return
	//}
	//逻辑处理

	sliceInt := []int{17, 15}
	marketPriceMap := make(map[string]interface{})
	for _, v := range sliceInt {
		d := mysql.GetAssetName(v)
		marketKey := fmt.Sprintf("%s.MarketPrice", d.TypeName)
		data, err := redis.GetData(marketKey)
		if err != nil {
			logger.Info(data)
			return
		}
		//返回数据
		marketPriceMap[strconv.Itoa(d.TypeId)] = data
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": marketPriceMap,
	})

}

//SetStartParamHandler 设置启动参数
func SetStartParamHandler(c *gin.Context) {
	//获取参数
	var p model.ParamStart
	err := c.Bind(&p)
	if err != nil {
		logger.Info(err)
		return
	}
	if p.Buy == 0 || p.Sale == 0 || p.Safe == 0 {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "缺少相关参数",
			"data": "",
		})
		return
	}
	//逻辑处理
	err1 := redis.CreateDurableKey("buy", p.Buy)
	if err != nil {
		logger.Info(err1)
		return
	}
	err2 := redis.CreateDurableKey("sale", p.Sale)
	if err != nil {
		logger.Info(err2)
		return
	}
	err3 := redis.CreateDurableKey("safe", p.Safe)
	if err != nil {
		logger.Info(err3)
		return
	}
	//返回参数
	c.JSON(200, gin.H{
		"msg":  "ok",
		"code": 200,
		"data": "",
	})
}

// GetBuyDataHandle 返回买入卖出的数据
func GetBuyDataHandle(c *gin.Context) {

	//通过查询最新10条买入数据
	data1 := mysql.GetBuyData(1)
	result1 := make([]model.RespBuy, len(data1))
	for i, v := range data1 {
		result1[i].Gid = v.Gid
		result1[i].Name = v.Name
		result1[i].Count = v.Count
		result1[i].TokenId = v.TokenId
		result1[i].MarketPrice = v.MarketPrice
		result1[i].SaleAddress = v.SaleAddress
		result1[i].FixedPrice = v.FixedPrice
		result1[i].TotalPrice = v.FixedPrice * float64(v.Count)
		result1[i].CreateTime = v.CreatedAt
		result1[i].Type = v.Type
		result1 = append(result1, result1[i])
	}
	//通过查询卖出的最新10条数据
	data2 := mysql.GetBuyData(2)
	result2 := make([]model.RespBuy, len(data2))
	for i, v := range data2 {
		result2[i].Gid = v.Gid
		result2[i].Name = v.Name
		result2[i].Count = v.Count
		result2[i].MarketPrice = v.MarketPrice
		result2[i].SaleAddress = v.SaleAddress
		result2[i].Profit = v.Profit
		result2[i].TokenId = v.TokenId
		result2[i].FixedPrice = v.FixedPrice
		result2[i].TotalPrice = v.FixedPrice * float64(v.Count)
		result2[i].CreateTime = v.CreatedAt
		result2[i].Type = v.Type
		result2 = append(result2, result2[i])
	}

	c.JSON(200, gin.H{
		"buy_data":  result1,
		"sale_data": result2,
		"msg":       "ok",
		"code":      200,
	})
}

//前端系统监控使用

//SetMngRiskHandle 设置风控
func SetMngRiskHandle(c *gin.Context) {
	//获取参数
	var p model.ParamRiskMng
	err := c.Bind(&p)
	if err != nil {
		logger.Info(err)
		return
	}
	if p.Situation == "" || p.TimeLevel == 0 || p.Percentage == 0 || p.OperationType == 0 || p.Status == 0 {
		c.JSON(500, gin.H{
			"msg":  "缺少相关参数",
			"code": 500,
			"data": "",
		})
		return
	}
	//把数据存进redis 中的哈希表
	m := make(map[string]interface{})
	m["Situation"] = p.Situation
	m["TimeLevel"] = p.TimeLevel
	m["Percentage"] = p.Percentage
	m["OperationType"] = p.OperationType
	m["Status"] = p.Status
	logger.Info(m)
	redis.CreatHashKey(fmt.Sprintf("risk:%s", p.Situation), m)
	//返回参数
	c.JSON(200, gin.H{
		"msg":  "ok",
		"code": 200,
		"data": "",
	})

}

//SetBuyAndSaleHandle 设置买入卖出百分比
func SetBuyAndSaleHandle(c *gin.Context) {
	//获取参数
	var p model.ParamBuyAndSale
	err := c.Bind(&p)
	if err != nil {
		logger.Info(err)
		return
	}
	if p.ProductID == 0 {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "缺少相关参数",
			"data": "",
		})
		return
	}
	//逻辑处理 1.根据ID找到前缀
	asset := mysql.GetAssetName(p.ProductID)
	productName := asset.TypeName
	//把数据存进redis 中的哈希表
	m := make(map[string]interface{})
	m["ProductName"] = productName
	m["RisePercentage"] = p.RisePercentage
	m["FallPercentage"] = p.FallPercentage
	m["RiseStatus"] = p.RiseStatus
	m["FallStatus"] = p.FallStatus
	logger.Info(m)
	redis.CreatHashKey(fmt.Sprintf("buyAndSale:%s", productName), m)
	//返回参数
	c.JSON(200, gin.H{
		"msg":  "ok",
		"code": 200,
		"data": "",
	})
}

//SetParamOnOffHandle 设置买入卖出总开关
func SetParamOnOffHandle(c *gin.Context) {
	//获取参数
	var p model.ParamOnOff
	err := c.Bind(&p)
	if err != nil {
		logger.Info(err)
		return
	}
	//if p.CrlName==""||p.Super==0{
	//	c.JSON(500,gin.H{
	//		"code":500,
	//		"msg":"缺少相关参数",
	//		"data":"",
	//	})
	//	return
	//}
	//把数据存进redis 中的哈希表
	m := make(map[string]interface{})
	m["CrlName"] = p.CrlName
	m["Super"] = p.Super
	logger.Info(m)
	redis.CreatHashKey(fmt.Sprintf("buyAndSale:%s", p.CrlName), m)
	//返回参数
	c.JSON(200, gin.H{
		"msg":  "ok",
		"code": 200,
		"data": "",
	})
}

//GetScriptStatusHandle 获取脚本运行的状态
func GetScriptStatusHandle(c *gin.Context) {
	//获取买入卖出总开关
	var buyStruct model.RespAllOnOff
	buy := redis.GetHashDataAll("buyAndSale:buy")
	mapstructure.Decode(buy, &buyStruct)

	var saleStruct model.RespAllOnOff
	sale := redis.GetHashDataAll("buyAndSale:sale")
	mapstructure.Decode(sale, &saleStruct)

	//获取买入卖出元兽蛋开关
	//var Egg model.RespBuyAndSale
	//egg := redis.GetHashDataAll("buyAndSale:Metamon Egg")
	//mapstructure.Decode(egg,&Egg)

	//获取买入卖出药水开关
	//var Potion model.RespBuyAndSale
	//potion := redis.GetHashDataAll("buyAndSale:Potion")
	//mapstructure.Decode(potion,&Potion)

	//通过reids获取市场价格
	eggMarket,_ := redis.GetData("Metamon Egg.MarketPrice")
	PotionMarket,_ := redis.GetData("Potion.MarketPrice")

	//获取元兽蛋买入数据
	var eggBuy model.RespBuyAndSaleSet
	eggBuy2 := redis.GetHashDataAll("BuySet:17")
	logger.Info(eggBuy2)
	logger.Info(eggBuy)
	err := mapstructure.Decode(eggBuy2, &eggBuy)
	eggBuy.AotuMarketprice = eggMarket
	if err != nil {
		logger.Error(err)
		return
	}
	eggBuy.ProductId = eggBuy2["product_id"]
	logger.Info(eggBuy)
	//获取元兽蛋卖出入数据
	var eggSale model.RespBuyAndSaleSet
	egg_sale := redis.GetHashDataAll("SaleSet:17")
	mapstructure.Decode(egg_sale, &eggSale)
	eggSale.AotuMarketprice = eggMarket
	eggSale.ProductId = egg_sale["product_id"]
	//获取药水买出数据
	var potionBuy model.RespBuyAndSaleSet
	potion_buy := redis.GetHashDataAll("BuySet:15")
	mapstructure.Decode(potion_buy, &potionBuy)
	potionBuy.AotuMarketprice = PotionMarket
	potionBuy.ProductId = potion_buy["product_id"]
	//获取药水卖出数据
	var potionSale model.RespBuyAndSaleSet
	potion_sale := redis.GetHashDataAll("SaleSet:15")
	mapstructure.Decode(potion_sale, &potionSale)
	potionSale.AotuMarketprice = PotionMarket
	potionSale.ProductId = potion_sale["product_id"]

	var all model.RespAllSwitch
	allOnOffSlice := make([]model.RespAllOnOff, 2)
	buyAndSaleSetSlice := make([]model.RespBuyAndSaleSet, 4)

	buyAndSaleSetSlice[0] = eggBuy
	buyAndSaleSetSlice[1] = eggSale
	buyAndSaleSetSlice[2] = potionBuy
	buyAndSaleSetSlice[3] = potionSale

	allOnOffSlice[0] = buyStruct
	allOnOffSlice[1] = saleStruct

	all.AllOnOff = allOnOffSlice
	all.BuyAndSale = buyAndSaleSetSlice
	c.JSON(200, gin.H{
		"msg":  "ok",
		"code": 200,
		"data": all,
	})

}

//GetRiskMonitorHandle  返回监控信息状态
func GetRiskMonitorHandle(c *gin.Context) {
	var fall model.RespRiskMonitor
	var rise model.RespRiskMonitor
	fallMap := redis.GetHashDataAll("risk:fall")
	mapstructure.Decode(fallMap, &fall)
	riseMap := redis.GetHashDataAll("risk:rise")
	mapstructure.Decode(riseMap, &rise)

	all := make([]model.RespRiskMonitor, 2)
	all[0] = fall
	all[1] = rise

	c.JSON(200, gin.H{
		"msg":  "ok",
		"code": 200,
		"data": all,
	})
}

//GetMarketPriceLineHandle 获取对应的市场数据
func GetMarketPriceLineHandle(c *gin.Context) {
	//参数
	var p model.ParamTypeId
	Berr := c.Bind(&p)
	if Berr != nil {
		logger.Error(Berr)
		return
	}

	//获取通过id获取类型的名称
	d := mysql.GetAssetName(p.TypeId)
	marketPriceKey := fmt.Sprintf("%s.MarketPrice", d.TypeName)
	logger.Info(marketPriceKey)

	//一小时前数据 3600
	time1 := mysql.GetHistoryMarketData(3600, marketPriceKey)
	//两小时前数据 7200
	time2 := mysql.GetHistoryMarketData(7200, marketPriceKey)
	//三小时前数据 10800
	time3 := mysql.GetHistoryMarketData(10800, marketPriceKey)
	//四小时前数据 14400
	time4 := mysql.GetHistoryMarketData(14400, marketPriceKey)
	//五小时前数据 18000
	time5 := mysql.GetHistoryMarketData(18000, marketPriceKey)
	//六小时前数据 21600
	time6 := mysql.GetHistoryMarketData(21600, marketPriceKey)

	//strtime1 := strconv.FormatFloat(time1.MarketData, 'E', -1, 64)
	//strtime2 := strconv.FormatFloat(time2.MarketData, 'E', -1, 64)
	//strtime3 := strconv.FormatFloat(time3.MarketData, 'E', -1, 64)
	//strtime4 := strconv.FormatFloat(time4.MarketData, 'E', -1, 64)
	//strtime5 := strconv.FormatFloat(time5.MarketData, 'E', -1, 64)
	//strtime6 := strconv.FormatFloat(time6.MarketData, 'E', -1, 64)

	var timeSlice []float64
	timeSlice = append(timeSlice, time1.MarketData, time2.MarketData, time3.MarketData, time4.MarketData, time5.MarketData, time6.MarketData)

	//返回数据
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": timeSlice,
	})

}

// GetIncomeHandle 查询当前利润
func GetIncomeHandle(c *gin.Context) {
	//获取利润
	data, Gerr := redis.GetData("income")
	if Gerr != nil {
		logger.Error(Gerr)
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": data,
	})
}

//GetAssetType  返回资产类型id列表
func GetAssetType(c *gin.Context) {
	//返回参数
	assetType := mysql.GetAssetType()
	data := make([]int, 0)
	logger.Info(data)
	//for _, v := range assetType {
	//	data = append(data, v.TypeId)
	//}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": assetType,
	})
	return
}

//SetBuySetHandle 设置买入出参数
func SetBuySetHandle(c *gin.Context) {
	var p model.ParamBuyAndSaleSet
	//var reps model.RespBuyAndSaleSet
	Eerr := c.Bind(&p)
	if Eerr != nil {
		logger.Error(Eerr)
		return
	}
	//把数据存进redis 中的哈希表
	m := make(map[string]interface{})
	m["product_id"] = p.ProductId
	m["percent"] = p.Percent
	m["status"] = p.Status
	m["types"] = p.Types
	m["market_price"] = p.MarketPrice
	logger.Info(m)
	redis.CreatHashKey(fmt.Sprintf("BuySet:%s", p.ProductId), m)
	//返回参数
	c.JSON(200, gin.H{
		"msg":  "ok",
		"code": 200,
		"data": "",
	})
}

//SetSaleSetHandle 设置买入出参数
func SetSaleSetHandle(c *gin.Context) {
	var p model.ParamBuyAndSaleSet
	//var reps model.RespBuyAndSaleSet
	Eerr := c.Bind(&p)
	if Eerr != nil {
		logger.Error(Eerr)
		return
	}
	//把数据存进redis 中的哈希表
	m := make(map[string]interface{})
	m["product_id"] = p.ProductId
	m["percent"] = p.Percent
	m["status"] = p.Status
	m["types"] = p.Types
	m["market_price"] = p.MarketPrice
	logger.Info(m)
	redis.CreatHashKey(fmt.Sprintf("SaleSet:%s", p.ProductId), m)
	//返回参数
	c.JSON(200, gin.H{
		"msg":  "ok",
		"code": 200,
		"data": "",
	})
}

//SetRiskPotionHandle 设置药水风控接口
func SetRiskPotionHandle(c *gin.Context) {
	//获取参数
	var p model.ParamRiskMng
	err := c.Bind(&p)
	if err != nil {
		logger.Info(err)
		return
	}
	if p.Situation == "" || p.TimeLevel == 0 || p.Percentage == 0 || p.OperationType == 0 || p.Status == 0 {
		c.JSON(500, gin.H{
			"msg":  "缺少相关参数",
			"code": 500,
			"data": "",
		})
		return
	}
	//把数据存进redis 中的哈希表
	m := make(map[string]interface{})
	m["Situation"] = p.Situation
	m["TimeLevel"] = p.TimeLevel
	m["Percentage"] = p.Percentage
	m["OperationType"] = p.OperationType
	m["Status"] = p.Status
	logger.Info(m)
	redis.CreatHashKey(fmt.Sprintf("risk:potion:%s", p.Situation), m)
	//返回参数
	c.JSON(200, gin.H{
		"msg":  "ok",
		"code": 200,
		"data": "",
	})
}

func GetRiskPotionHandle(c *gin.Context) {
	var fall model.RespRiskMonitor
	var rise model.RespRiskMonitor
	fallMap := redis.GetHashDataAll("risk:potion:fall")
	mapstructure.Decode(fallMap, &fall)
	riseMap := redis.GetHashDataAll("risk:potion:rise")
	mapstructure.Decode(riseMap, &rise)

	all := make([]model.RespRiskMonitor, 2)
	all[0] = fall
	all[1] = rise

	c.JSON(200, gin.H{
		"msg":  "ok",
		"code": 200,
		"data": all,
	})
}

//SetSellingRateHandle 设置卖出率
func SetSellingRateHandle(c *gin.Context) {
	var p model.ParamSellingRate
	BErr := c.Bind(&p)
	if BErr != nil {
		logger.Error(BErr)
		return
	}

	m := make(map[string]interface{})
	m["time_level"] = p.TimeLevel
	m["percent"] = p.Percent
	m["status"] = p.Status
	m["operation_type"] = p.OperationType
	logger.Info(m)
	redis.CreatHashKey("SellingRate", m)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": "",
	})
}

//GetSellingRateHandle 返回卖出率
func GetSellingRateHandle(c *gin.Context) {
	var resp model.RespSellingRate
	sellRate := redis.GetHashDataAll("SellingRate")
	mapstructure.Decode(sellRate, &resp)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": sellRate,
	})

}

//GetProportionHandle 返回市场占比
func GetProportionHandle(c *gin.Context) {
	//获取参数
	var p model.ParamProportion
	BErr := c.Bind(&p)
	if BErr != nil {
		logger.Error(BErr)
		return
	}
	if p.TypeId == 0 {
		p.TypeId = 17
	}
	//获取redis里面的数据
	data := redis.GetHashDataAll(fmt.Sprintf("Proportion:%d", p.TypeId))
	data1 := redis.GetHashDataAll(fmt.Sprintf("ProportionCount:%d", p.TypeId))
	type M struct {
		Key float64
		Val int
		Count float64
	}
	var m []M
	for i, v := range data {
		//logger.Info(v)
		for i1,v1 := range data1{
			//logger.Info(v1)
			if i1 == i{
				f, _ := strconv.ParseFloat(i, 64)
				t, _ := strconv.Atoi(v)
				count, _ := strconv.ParseFloat(v1, 64)
				m = append(m, M{
					Key: f,
					Val: t,
					Count: count,
				})
			}
		}
	}
	//返回数据
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": m,
	})
}

//GetBuyAndSaleHandle 获取买和出的数据根据时间轴返回
func GetBuyAndSaleHandle(c *gin.Context)  {
	//获取当前时间戳
	now := utils.GetNowTimeS()
	oneHour := now - 3600
	twoHour := now -7200
	threeHour := now - 10800
	fourHour := now - 14400
	fiveHour := now - 18000
	sixHour := now - 21600
	//把时间转成str
	nowStr := utils.TimestampToDatetime(now)
	logger.Info(nowStr)
	oneHourStr:= utils.TimestampToDatetime(oneHour)
	logger.Info(oneHourStr)
	twoHourStr:=utils.TimestampToDatetime(twoHour)
	threeHourStr:=utils.TimestampToDatetime(threeHour)
	fourHourStr:=utils.TimestampToDatetime(fourHour)
	fiveHourStr:=utils.TimestampToDatetime(fiveHour)
	sixHourStr:=utils.TimestampToDatetime(sixHour)


	//返回每小时买入数量
	c1 := mysql.GetBuyCount(oneHourStr,nowStr)
	c2 := mysql.GetBuyCount(twoHourStr,oneHourStr)
	c3 :=mysql.GetBuyCount(threeHourStr,twoHourStr)
	c4 :=mysql.GetBuyCount(fourHourStr,threeHourStr)
	c5 :=mysql.GetBuyCount(fiveHourStr,fourHourStr)
	c6 :=mysql.GetBuyCount(sixHourStr,fiveHourStr)


	//每小时买出出数量
	d1 := mysql.GetSaleCount(oneHourStr,nowStr)
	d2 :=mysql.GetSaleCount(twoHourStr,oneHourStr)
	d3 :=mysql.GetSaleCount(threeHourStr,twoHourStr)
	d4 :=mysql.GetSaleCount(fourHourStr,threeHourStr)
	d5 :=mysql.GetSaleCount(fiveHourStr,fourHourStr)
	d6 :=mysql.GetSaleCount(sixHourStr,fiveHourStr)

	//返回数据
	type TimeToBuy struct {
		Times string `json:"times"`
		BuyCount float64 `json:"buy_count"`
		SaleCount float64 `json:"sale_count"`
	}

	timeToBuy := make([]TimeToBuy,6)
	timeToBuy[0].Times = "one"
	timeToBuy[0].BuyCount = c1
	timeToBuy[0].BuyCount = d1

	timeToBuy[1].Times = "two"
	timeToBuy[1].BuyCount = c2
	timeToBuy[1].BuyCount = d2

	timeToBuy[2].Times = "three"
	timeToBuy[2].BuyCount = c3
	timeToBuy[2].BuyCount = d3

	timeToBuy[3].Times = "four"
	timeToBuy[3].BuyCount = c4
	timeToBuy[3].BuyCount = d4

	timeToBuy[4].Times = "five"
	timeToBuy[4].BuyCount = c5
	timeToBuy[4].BuyCount = d5

	timeToBuy[5].Times = "six"
	timeToBuy[5].BuyCount = c6
	timeToBuy[5].BuyCount = d6
	logger.Info(timeToBuy)
	c.JSON(200,gin.H{
		"code":200,
		"msg":"ok",
		"data":timeToBuy,
	})

}


