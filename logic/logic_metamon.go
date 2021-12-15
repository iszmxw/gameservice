/**
 @author:way
 @date:2021/11/26
 @note
**/

package logic

import (
	"encoding/json"
	"fmt"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/model"
	"redisData/pkg/logger"
	"redisData/utils"
	"sort"
	"strconv"
	"strings"
	"time"
)


// GetKeysByPfx 根据前缀遍历key 拼接数据
func GetKeysByPfx(keypfx string) ([]model.RespAssetsDetailList, error) {
	pfx := fmt.Sprintf("%s:",keypfx)
	dataList, err := redis.GetKeysByPfx(pfx) //dataList是一个key集合
	//对list进行排序，只取前面50条数据，展示使用
	sort.Slice(dataList, func(i, j int) bool {
		return dataList[i] > dataList[j]
	})
	dataList = dataList[:50]

	logger.Info(dataList)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	//使用map来存
	var dataDetailList []model.RespAssetsDetailList
	for _, v := range dataList {
		fmt.Println(v)
		res, RErr := redis.GetDataByKey(v)
		if RErr != nil {
			logger.Error(RErr)
			return nil, RErr
		}
		dataDetailList = append(dataDetailList, res)
	}

	//使用切片来存
	//dataDetailList := make([]string,30)
	//for _,v := range dataList{
	//	dataDetailList = append(dataDetailList,v)
	//}
	return dataDetailList, nil
}

// SortSlice 输入一个切片然后进行排序，得出比重最多的价格，作为市场价
func SortSlice(priceList []float64,productID int) (marketPrice []float64) {
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
		return
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
	if productID != 0{
		m2 := make(map[string]interface{})
		for i,v := range m1{
			str := strconv.FormatFloat(i, 'E', -1, 64)
			m2[fmt.Sprintf("%s",str)] = v
		}
		redis.CreatHashKey(fmt.Sprintf("Proportion:%d",productID),m2)
	}

	return s3
}

//SetBuyALG 设置买入买出算法,无返回值
func SetBuyALG(key string,marketPrice float64, percentage float64) {
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
	if data.List == nil{
		logger.Info("拿不到数据")
		return
	}

	for _, v := range data.List {
		//查看集合是否已经存在，存在就不加入列表了
		flag := redis.ExistEle("buySet1", strconv.Itoa(v.Id))
		//fmt.Println(time.Now())
		//logger.Info(flag)
		//logger.Info(strconv.Itoa(v.Id))
		if flag == true {
			logger.Info("该资产已经购买")
			continue
		}
		//计算单价
		currentPrice, FErr := strconv.ParseFloat(v.FixedPrice, 64)
		if FErr != nil {
			logger.Error(FErr)
			return
		}
		avgPrice := currentPrice / float64(v.Count)
		//logger.Info(marketPrice * percentage-avgPrice)
		//logger.Info(fmt.Sprintf("%v * %v-%v",marketPrice,(percentage*0.01),avgPrice))
		//logger.Info((marketPrice *(percentage*0.01))-avgPrice)
		percentage = 100 - percentage
		if (marketPrice * (percentage*0.01)) > avgPrice  {
			//买入时加入一个多一个标识
			memberKey := fmt.Sprintf("%s:%s",v.Name,strconv.Itoa(v.Id))


			redis.CreateZScoreData("buySet",memberKey,avgPrice)
			redis.CreateSetData("buySet1", strconv.Itoa(v.Id))
			dataDetails := mysql.GetDataByGid(strconv.Itoa(v.Id))
			buy := model.Buy{
				Gid: strconv.Itoa(v.Id),
				Name: v.Name,
				Count: v.Count,
				FixedPrice: avgPrice,
				Type: 1,
				SaleAddress: v.SaleAddress,
				TokenId: v.TokenId,
				MarketPrice: marketPrice,
				IdInContract: dataDetails.IdInContract,
			}
			mysql.InsertBuyRecord(buy)
			//同时删除redis中的key ,下次即使爬到了也自动忽略
			redis.DeleEggKey(strconv.Itoa(v.Id))
		}

	}
	return
}

// SetMarketPrice 获取市场数据，缓存到redis,同时存到MySQL
func SetMarketPrice(key string) {
	//从redis中获取市场最新的数据
	redisdata, err := redis.GetData(key)
	//logger.Info(redisdata)
	if err != nil {
		logger.Error(err)
		return
	}
	if len(redisdata) < 150{
		return
	}
	//计算出市场价格存进redis和mysql
	// 序列化返回的结果
	var data model.Data
	if Uerr := json.Unmarshal([]byte(redisdata), &data); Uerr != nil {
		logger.Error(Uerr)
	}
	//做逻辑运算 1.算出均价 2.确定价格 3.返回价格 先不存redis
	//list存下全部均价
	list := make([]float64, 0, len(data.List))

	m := make(map[float64]float64)
	for _, v := range data.List {
		fixedPrice, FErr := strconv.ParseFloat(v.FixedPrice, 64)
		if FErr != nil {
			logger.Error(err)
			return
		}
		count := float64(v.Count)
		price := fixedPrice / count
		list = append(list, price)
		//使用一个map统计数量

		if m[price] == 0 {
			m[price] = count
		} else {
			m[price] = m[price] + count
		}
		//map------
	}
	//使用redis存起来
	tid := NameTranType(key)
	m2 := make(map[string]interface{})
	for i,v := range m{
		str := strconv.FormatFloat(i, 'E', -1, 64)
		m2[fmt.Sprintf("%s",str)] = v
	}
	redis.CreatHashKey(fmt.Sprintf("ProportionCount:%d",tid),m2)
	if len(list) <=0 {
		logger.Info("list为空")
		return
	}
	tranType := NameTranType(key)
	if tranType == 0{
		logger.Info("没有对应的产品")
		return
	}
	//logger.Info(tranType)
	list1 := SortSlice(list,tranType)

	//logger.Info(list1[0])
	marketKey := fmt.Sprintf("%s.MarketPrice",data.List[0].Name)
	err = redis.CreateDurableKey(marketKey, list1[0])
	//logger.Info("添加Market成功")
	if err != nil {
		logger.Error(err)
		return
	}
	//存进mysql
	data1 := model.MarketData{
		MarketName: marketKey,
		MarketData: list1[0],
	}
	mysql.InsertMarketPrice(data1)
	time.Sleep(500*time.Millisecond)
}

//GetMarketDataByRedis 根据redis的历史数据，算出历史的数据
func GetMarketDataByRedis(assetsListKey string) float64 {
	//从redis中获取鸡蛋的数据
	redisdata, err := redis.GetData(assetsListKey)
	if err != nil {
		logger.Error(err)
		return 0
	}
	//计算出市场价格存进redis和mysql
	// 序列化返回的结果
	var data model.Data
	if Uerr := json.Unmarshal([]byte(redisdata), &data); Uerr != nil {
		logger.Error(Uerr)
	}
	//做逻辑运算 1.算出均价 2.确定价格 3.返回价格 先不存redis
	//list存下全部均价
	list := make([]float64, 0, len(data.List))
	for _, v := range data.List {
		fixedPrice, FErr := strconv.ParseFloat(v.FixedPrice, 64)
		conut := float64(v.Count)
		if v.Count != 1 {

			if FErr != nil {
				logger.Error(err)
				return 0
			}
			price := fixedPrice / conut
			list = append(list, price)
		}
		if v.Count == 1 {
			list = append(list, fixedPrice)
		}

	}
	if len(list) <=0 {
		logger.Info("list为空")
		return 0
	}
	//算出市场价格
	marketDataList := SortSlice(list,0)
	return marketDataList[0]
}

// SetSaleALG 设置卖出算法
func SetSaleALG(marketPriceKey string,account float64,percentage float64)  {
	strData, err := redis.GetData(marketPriceKey)
	//把产品名称切割出来
	product := utils.Split(marketPriceKey,".")[0]
	logger.Info(product)
	if err != nil {
		fmt.Println(err)
		return
	}
	MarketPrice,err := strconv.ParseFloat(strData,64)
	if err != nil{
		fmt.Println(err)
		return
	}
	//读入买入数据,对比市场价决定卖出
	strSlice := redis.GetAllZSet("buySet")
	for _,v  := range strSlice{
		//判断下种类,不是同一个种类跳过
		productName := utils.Split(v,":")[0]
		if product != productName{
			continue
		}
		score := redis.GetScoreByMember("buySet", v)
		//判断是否卖出
		if score.(float64) *(percentage+100) *0.01 < MarketPrice*0.99  {
			//定义卖出价格
			salePrice := MarketPrice * 0.99
			account = account + salePrice - score.(float64)
			logger.Info("赚了")
			logger.Info(score.(float64) *(percentage+100) *0.01)
			logger.Info(score.(float64))
			logger.Info(salePrice - score.(float64))
		}
		redis.DeleteSetData("buySet1",utils.Split(v,":")[1])
		redis.DeleteRecByMember("buySet",v)
		//mysql根据gid查询一下
		gId := utils.Split(v,":")[1]
		buyData := mysql.GetBuyById(gId)
		//mysql添加买出记录
		buy := model.Buy{
			Gid:utils.Split(v,":")[1],
			Name: utils.Split(v,":")[0],
			MarketPrice: MarketPrice,
			Count: buyData[0].Count,
			FixedPrice: score.(float64),
			Type: 2,
			SaleAddress: buyData[0].SaleAddress,
			IdInContract: buyData[0].IdInContract,
			Profit: MarketPrice * 0.99 - score.(float64),
			SalePrice: MarketPrice * 0.99,
		}
		mysql.InsertBuyRecord(buy)
	}
	//移除redis中buySet集合
	//count 存redis
	Cerr := redis.CreateDurableKey("income", account)
	if Cerr != nil {
		fmt.Println(err)
		return
	}
}

// SetSaleALG2 设置卖出算法，买出价格由配置文件决定
func SetSaleALG2(marketPriceKey string,account float64,percentage float64)  {
	var marketPrice string
	t := NameTranType(marketPriceKey)
	if t == 17{
		data := redis.GetHashDataAll("SaleSet:17")
		marketPrice = data["market_price"]
	}
	if t == 15{
		data := redis.GetHashDataAll("SaleSet:15")
		marketPrice = data["market_price"]
	}

	strData := marketPrice
	//logger.Info(strData)
	//把产品名称切割出来
	product := utils.Split(marketPriceKey,".")[0]
	//logger.Info(product)
	MarketPrice,err := strconv.ParseFloat(strData,64)
	if err != nil{
		logger.Info(err)
		return
	}
	//读入买入数据,对比市场价决定卖出
	strSlice := redis.GetAllZSet("buySet")
	for _,v  := range strSlice{
		//判断下种类,不是同一个种类跳过
		productName := utils.Split(v,":")[0]
		if product != productName{
			continue
		}
		score := redis.GetScoreByMember("buySet", v)
		//判断是否卖出
		if score.(float64) *(percentage+100) *0.01 < MarketPrice*0.99  {
			//定义卖出价格
			salePrice := MarketPrice * 0.99
			account = account + salePrice - score.(float64)
			logger.Info("赚了")
			logger.Info(score.(float64) *(percentage+100) *0.01)
			logger.Info(score.(float64))
			logger.Info(salePrice - score.(float64))
		}
		redis.DeleteSetData("buySet1",utils.Split(v,":")[1])
		redis.DeleteRecByMember("buySet",v)
		//mysql根据gid查询一下
		gId := utils.Split(v,":")[1]
		buyData := mysql.GetBuyById(gId)
		//mysql添加买出记录
		buy := model.Buy{
			Gid:utils.Split(v,":")[1],
			Name: utils.Split(v,":")[0],
			MarketPrice: MarketPrice,
			Count: buyData[0].Count,
			FixedPrice: score.(float64),
			SaleAddress: buyData[0].SaleAddress,
			IdInContract: buyData[0].IdInContract,
			Type: 2,
			Profit: MarketPrice * 0.99 - score.(float64),
			SalePrice: MarketPrice * 0.99,
		}
		mysql.InsertBuyRecord(buy)
	}
	//移除redis中buySet集合
	//count 存redis
	Cerr := redis.CreateDurableKey("income", account)
	if Cerr != nil {
		fmt.Println(err)
		return
	}
}

// NameTranType 名字转回类型
func NameTranType(name string)int {
	switch  {
	case strings.Contains(name,"Metamon Egg") :
		return 17
	case strings.Contains(name,"Potion"):
		return 15
	}
	return 0
}

//StopScript 停止脚本
func StopScript(){
	//获取买卖脚本的配置文件，修改其中的状态
	//buy15 := redis.GetHashDataAll("BuySet:15")
	//buy15["status"] = "2"
	//buy17 := redis.GetHashDataAll("BuySet:17")
	//buy17["status"] = "2"
	//sale15 := redis.GetHashDataAll("SaleSet:15")
	//sale15["status"] = "2"
	//sale17 := redis.GetHashDataAll("SaleSet:17")
	//sale17["status"] = "2"


	//买入药水参数设置
	buySet15 := make(map[string]interface{}, 4)
	buySet15["product_id"] = "15"
	buySet15["percent"] = "10"
	buySet15["market_price"] = "7000"
	buySet15["status"] = "2" //1.打开 2.关闭
	buySet15["types"] = "1"  //买入固定为1
	redis.CreatHashKey("BuySet:15", buySet15)
	//买入元兽蛋参数设置
	buySet17 := make(map[string]interface{}, 4)
	buySet17["product_id"] = "17"
	buySet17["market_price"] = "70000"
	buySet17["percent"] = "10"
	buySet17["status"] = "2" //1.打开 2.关闭
	buySet17["types"] = "1"  //买入固定为1
	redis.CreatHashKey("BuySet:17", buySet17)

	//卖出药水参数设置
	SaleSet15 := make(map[string]interface{}, 4)
	SaleSet15["product_id"] = "15"
	SaleSet15["percent"] = "10"
	SaleSet15["market_price"] = "9000"
	SaleSet15["status"] = "2"
	SaleSet15["types"] = "2"	 //买入固定为2
	redis.CreatHashKey("SaleSet:15", SaleSet15)
	//买入元兽蛋参数设置
	SaleSet17 := make(map[string]interface{}, 4)
	SaleSet17["product_id"] = "17"
	SaleSet17["percent"] = "10"
	SaleSet17["market_price"] = "90000"
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

	logger.Info("停止脚本完成")

}

//ManageData 数据处理逻辑1.判断数据是否存在集合中 是就跳过，不是存在list里面
func ManageData(data *model.ResponseDataList){
	if data == nil{
		logger.Info("数据为空")
		return
	}
	for _,v := range data.List{
		if redis.ExistEle("assertSet",strconv.Itoa(v.Id)){
			continue
		}
		//如果不是把数据存进redis队列
		marshal, Merr := json.Marshal(v)
		if Merr != nil {
			logger.Error(Merr)
			return
		}
		//使用list储存
		redis.SetOneList("assertList",string(marshal))
		//使用集合把ID储存起来
		redis.CreateSetData("assertSet",strconv.Itoa(v.Id))
	}
}

//StoreListToMysql 把redis中队列中的数据储存到mysql
func StoreListToMysql(str string)  {
	logger.Info(str)
	l :=  model.List{}
	err := json.Unmarshal([]byte(str), &l)
	if err != nil {
		logger.Error(err)
		return
	}
	d := model.AssetsData{
		GId: strconv.Itoa(l.Id),
		Name: l.Name,
		FixedPrice: l.FixedPrice,
		HighestPrice: l.HighestPrice,
		ImageUrl: l.ImageUrl,
		Count: l.Count,
		SaleType: l.SaleType,
		TokenId: l.TokenId,
		SaleAddress: l.SaleAddress,
		Status: l.Status,
	}
	mysql.CreateOneAssert(d)
	RequestAssertsDetails(strconv.Itoa(l.Id))


}





