/**
 @author:way
 @date:2021/11/26
 @note
**/

package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/model"
	"redisData/pkg/logger"
	"strconv"
	"time"
)

// RequestEggData 获取鸡蛋的数据,返回id切片
func RequestEggData() ([]int, error) {
	url := "https://market-api.radiocaca.com/nft-sales?pageNo=1&pageSize=20&sortBy=created_at&order=desc&name=&saleType&category=17&tokenType"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	body, _ := ioutil.ReadAll(response.Body)
	// 序列化返回的结果
	var data model.Data
	if Uerr := json.Unmarshal(body, &data); Uerr != nil {
		logger.Error(Uerr)
	}
	//取出里面的ID返回到数组里面去
	idList := make([]int, 100)
	for _, v := range data.List {
		idList = append(idList, v.Id)
	}
	return idList, nil
}

// RequestPotionData 获取药水的数据,取出全部id拼接成切片
func RequestPotionData() ([]int, error) {
	url := "https://market-api.radiocaca.com/nft-sales?pageNo=1&pageSize=20&sortBy=created_at&order=desc&name=&saleType&category=15&tokenType"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	body, _ := ioutil.ReadAll(response.Body)
	// 序列化返回的结果
	var data model.Data
	if Uerr := json.Unmarshal(body, &data); Uerr != nil {
		logger.Error(Uerr)
	}
	//取出里面的ID返回到数组里面去
	idList := make([]int, 30)
	for _, v := range data.List {
		idList = append(idList, v.Id)
	}
	return idList, nil
}

//RequestDataDetail 根据id访问详情
func RequestDataDetail(id int) (detailData string) {
	url := fmt.Sprintf("https://market-api.radiocaca.com/nft-sales/%d", id)
	response, Rerr := http.Get(url)
	if Rerr != nil {
		logger.Error(Rerr)
		return
	}
	body, _ := ioutil.ReadAll(response.Body)
	return string(body)
}

//CreatEggData 遍历鸡蛋数据存进redis
func CreatEggData() {
	data, err := RequestEggData()
	if err != nil {
		logger.Error(err)
		return
	}
	for _, v := range data {
		if v == 0 {
			continue
		}
		//判断集合中是否存在,存在就是买入,然后跳过
		fmt.Println(time.Now())
		flag := redis.ExistEle("buySet1", strconv.Itoa(v))
		fmt.Println(time.Now())
		logger.Info(flag)
		logger.Info(strconv.Itoa(v))
		if flag == true {
			continue
		}
		go func() {
			DetailData := RequestDataDetail(v)
			redis.CreateEggData(strconv.Itoa(v), DetailData)
		}()

	}
}

//CreatPotionData 遍历药水数据存进redis
func CreatPotionData() {

	data, err := RequestPotionData()
	if err != nil {
		logger.Error(err)
		return
	}
	for _, v := range data {
		if v == 0 {
			continue
		}
		DetailData := RequestDataDetail(v)
		redis.CreatePotionData(strconv.Itoa(v), DetailData)
	}
}

// GetKeysByPfx 根据前缀遍历key 拼接数据
func GetKeysByPfx(keypfx string) ([]model.Rdata, error) {
	pfx := fmt.Sprintf("%s:",keypfx)
	dataList, err := redis.GetKeysByPfx(pfx) //dataList是一个key集合
	logger.Info(dataList)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	//使用map来存
	var dataDetailMap []model.Rdata
	for _, v := range dataList {
		fmt.Println(v)
		res, RErr := redis.GetDataByKey(v)
		if RErr != nil {
			logger.Error(RErr)
			return nil, RErr
		}
		dataDetailMap = append(dataDetailMap, res)
	}

	//使用切片来存
	//dataDetailList := make([]string,30)
	//for _,v := range dataList{
	//	dataDetailList = append(dataDetailList,v)
	//}
	return dataDetailMap, nil
}

//RequestGetEggPrice 通过请求拿到鸡蛋数据
func RequestGetEggPrice() ([]float64, error) {
	url := "https://market-api.radiocaca.com/nft-sales?pageNo=1&pageSize=20&sortBy=created_at&order=desc&name=&saleType&category=17&tokenType"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	body, _ := ioutil.ReadAll(response.Body)
	// 序列化返回的结果
	var data model.Data
	if Uerr := json.Unmarshal(body, &data); Uerr != nil {
		logger.Error(Uerr)
	}
	//做逻辑运算 1.算出均价 2.确定价格 3.返回价格 先不存redis
	//list存下全部均价
	list := make([]float64, 0, len(data.List))
	for _, v := range data.List {
		fixedPrice, FErr := strconv.ParseFloat(v.FixedPrice, 64)
		if v.Count != 1 {
			if FErr != nil {
				logger.Error(err)
				return nil, err
			}
			count := float64(v.Count)
			price := fixedPrice / count
			list = append(list, price)
		}
		if v.Count == 1 {
			list = append(list, fixedPrice)
		}
	}

	return list, nil
}

// SortSlice 输入一个切片然后进行排序，得出比重最多的价格，作为市场价
func SortSlice(priceList []float64) (marketPrice []float64) {
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
	//插入一条redis数据，把这次遍历市场价占比计算后返回


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
	return s3
}

// SetBuyALG 设置买入买出算法,无返回值
func SetBuyALG(marketPrice float64, percentage float64) {
	//1.从redis取出egg列表的数据
	//redisData, err := GetDataInRedis("eggDataList")
	//if err != nil {
	//	logger.Error(err)
	//	return
	//}

	//对redisData 进行反序列化
	//var data model.Data
	//Unmarshalerr := json.Unmarshal([]byte(redisData), &data)
	//if Unmarshalerr != nil {
	//	logger.Error(Unmarshalerr)
	//	return
	//}
	//联网访问买入的数据
	data := GetAssertsData(100, 17)
	if data == nil{
		logger.Info("拿不到数据")
		return
	}
	for _, v := range data.List {
		//查看集合是否已经存在，存在就不加入列表了
		flag := redis.ExistEle("buySet1", strconv.Itoa(v.Id))
		fmt.Println(time.Now())
		logger.Info(flag)
		logger.Info(strconv.Itoa(v.Id))
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
		logger.Info(marketPrice * percentage-avgPrice)
		logger.Info(fmt.Sprintf("%v * %v-%v",marketPrice,(percentage*0.01),avgPrice))
		logger.Info((marketPrice *(percentage*0.01))-avgPrice)
		percentage = 100 - percentage
		if (marketPrice * (percentage*0.01)) > avgPrice  {
			redis.CreateZScoreData("buySet",strconv.Itoa(v.Id),avgPrice)
			redis.CreateSetData("buySet1", strconv.Itoa(v.Id))
			buy := model.Buy{
				Gid: strconv.Itoa(v.Id),
				Name: v.Name,
				Count: v.Count,
				FixedPrice: avgPrice,
				Type: 1,
				SaleAddress: v.SaleAddress,
				TokenId: v.TokenId,
				MarketPrice: marketPrice,

			}
			mysql.InsertBuyRecord(buy)
			//同时删除redis中的key ,下次即使爬到了也自动忽略
			redis.DeleEggKey(strconv.Itoa(v.Id))

		}

	}
	return
}

// SetEggMarketPrice 获取市场数据，缓存到redis,同时存到MySQL
func SetEggMarketPrice() {
	//从redis中获取鸡蛋的数据
	redisdata, err := redis.GetData("eggDataList")
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
	for _, v := range data.List {
		fixedPrice, FErr := strconv.ParseFloat(v.FixedPrice, 64)
		count := float64(v.Count)
		if v.Count != 1 {

			if FErr != nil {
				logger.Error(err)
				return
			}
			price := fixedPrice / count
			list = append(list, price)
		}
		if v.Count == 1 {
			list = append(list, fixedPrice)
		}

	}
	//市场价等于list[0]
	//把市场价存进redis,存进mysql
	//strMarketPrice := strconv.FormatFloat(list[0], 'E', -1, 64)
	//logger.Info(strMarketPrice)
	//排序
	if len(list) <=0 {
		logger.Info("list为空")
		return
	}
	list1 := SortSlice(list)
	err = redis.CreateKey("eggMarket", list1[0])
	logger.Info("添加eggMarket成功")
	if err != nil {
		logger.Error(err)
		return
	}
	//存进mysql
	data1 := model.MarketData{
		MarketName: "egg",
		MarketData: list1[0],
	}
	mysql.InsertMarketPrice(data1)
}

// RiskControl 风险控制,传入最新的市场价格，和承受波动百分比
func RiskControl(marketPrice float64, currentMarketPricePrice float64,percentage float64) string {
	//当前市场价,从redis中取上一次的
	//var currentMarketPricePrice float64
	//oldMarkerPrice, _ := redis.GetData("eggMarket")
	//currentMarketPricePrice, err := strconv.ParseFloat(oldMarkerPrice, 64)
	//if err != nil {
	//	logger.Error(err)
	//	return ""
	//}
	if (marketPrice/currentMarketPricePrice)-1 >= (percentage * 0.01) {
		//停止买入脚本，且发邮件通知,使用上一次的market和现在的market对比，上一次的market从redis中读，新的marketPrice重新算
		return "目前涨幅超过预期百分比"
	}
	if 1-(currentMarketPricePrice/marketPrice) >= (percentage * 0.01) {
		//下架挂单并且重新上架，发邮件通知
		return "目前跌幅超过预期百分比"
	}
	return "当前数据稳定"
}

// SetDataInRedis 访问网上的数据保存到redis,定时逻辑在main函数上面加
func SetDataInRedis() error {
	url := "https://market-api.radiocaca.com/nft-sales?pageNo=1&pageSize=300&sortBy=created_at&order=desc&name=&saleType&category=17&tokenType"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	body, _ := ioutil.ReadAll(response.Body)
	if len(body) < 150 {
		logger.Info("访问频繁限制")
		return errors.New("访问频繁限制")
	}
	//fmt.Println(body)
	//把得到的数据存进redis key为eggData
	CreateKeyErr := redis.CreateKey(
		"eggDataList",
		string(body),
	)
	if CreateKeyErr != nil {
		logger.Error(CreateKeyErr)
		return CreateKeyErr
	}
	return nil

}

// GetDataInRedis 从redis中取出数据，运算使用
func GetDataInRedis(key string) (string, error) {
	data, err := redis.GetData(key)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return data, nil
}

//GetMarketDataByRedis 根据redis的历史数据，算出历史的数据
func GetMarketDataByRedis() float64 {
	//从redis中获取鸡蛋的数据
	redisdata, err := redis.GetData("eggDataList")
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
	marketDataList := SortSlice(list)
	return marketDataList[0]
}

// SetSaleALG 设置卖出算法
func SetSaleALG(account float64,percentage float64) float64 {
	strData, err := redis.GetData("eggMarket")
	if err != nil {
		fmt.Println(err)
		return 0
	}

	MarketPrice,err := strconv.ParseFloat(strData,64)
	if err != nil{
		fmt.Println(err)
		return 0
	}
	//读入买入数据,对比市场价决定卖出
	strSlice := redis.GetAllZSet("buySet")
	for _,v  := range strSlice{
		score := redis.GetScoreByMember("buySet", v)
		//判断是否卖出
		if score.(float64) *(percentage+100) *0.01 < MarketPrice{
			account = account + score.(float64) *(percentage+100) *0.01 - score.(float64)

			logger.Info("赚了")
			logger.Info(score.(float64) *(percentage+100) *0.01)
			logger.Info(score.(float64))
			logger.Info(score.(float64) *(percentage+100) *0.01 - score.(float64))
		}
		redis.DeleteSetData("buySet1",v)
		redis.DeleteRecByMember("buySet",v)
		//mysql添加买出记录
		buy := model.Buy{
			Gid:v,
			MarketPrice: MarketPrice,
			FixedPrice: score.(float64) *(percentage+100) *0.01,
			Type: 2,
			Profit: score.(float64) *(percentage+100) *0.01 - score.(float64),
		}
		mysql.InsertBuyRecord(buy)
	}
	//移除redis中buySet集合
	return account

}




