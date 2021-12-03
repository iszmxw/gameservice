/**
 @author:way
 @date:2021/11/30
 @note
**/

package logic

import (
	"encoding/json"
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

//GetAssertsData 请求获取数据
func GetAssertsData(pageSize int,category int) *model.ResponseDataList {
	//category 15是potion 17是egg
	url := fmt.Sprintf("https://market-api.radiocaca.com/nft-sales?pageNo=1&pageSize=%d&sortBy=created_at&order=desc&name=&saleType&category=%d&tokenType",pageSize,category)
	logger.Info(url)
	response, Gerr := http.Get(url)
	if Gerr != nil {
		logger.Error(Gerr)
		return nil
	}
	body, _ := ioutil.ReadAll(response.Body)

	//fmt.Println(body)
	//反序列化成结构体
	if len(body) < 150 {
		logger.Info("访问频繁，休息一分钟")
		time.Sleep(60*time.Second)
	}
	var d model.ResponseDataList
	UErr := json.Unmarshal(body,&d)
	if UErr != nil {
		logger.Error(UErr)
		return nil
	}
	return &d
}

//ManageData 数据处理逻辑1.判断数据是否存在集合中 是就跳过，不是存在list里面
func ManageData(data *model.ResponseDataList){

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

//SetAssertsDetails 访问获取一条详情数据，并且存入数据库
func SetAssertsDetails(gid string)  {
	url := fmt.Sprintf("https://market-api.radiocaca.com/nft-sales/%s", gid)
	response, RErr := http.Get(url)
	if RErr != nil {
		logger.Error(RErr)
		return
	}
	body, _ := ioutil.ReadAll(response.Body)
	//反序列化成结构体
	var d model.ResponseAssertsDetails
	UErr := json.Unmarshal(body,&d)
	if UErr != nil {
		logger.Error(UErr)
		return
	}
	//存进数据库
	assertDetails := model.AssetsDetails{
		Gid: strconv.Itoa(d.Data.Id),
		Name: d.Data.Name,
		Description: d.Data.Description,
		CreatedAt: d.Data.CreatedAt,
		ImageUrl: d.Data.ImageUrl,
		Count: d.Data.Count,
		FixedPrice: d.Data.FixedPrice,
		TotalPrice: d.Data.TotalPrice,
		SaleAddress: d.Data.SaleAddress,
		IdInContract: d.Data.IdInContract,
		TokenId: strconv.Itoa(d.Data.TokenId),
		TokenStandard: d.Data.TokenStandard,
		Owner: d.Data.Owner,
		NftAddress: d.Data.NftAddress,
		BlockChain: d.Data.BlockChain,
		StartTime: strconv.Itoa(d.Data.StartTime),
		Status: d.Data.Status,
		//Properties: d.Data.Properties,

	}
	mysql.CreateOneAssertDetails(assertDetails)
	if len(assertDetails.Status) >0  {
		//关联链上数据
		go RequestChainData(assertDetails.Gid)
	}else {
		logger.Info(assertDetails)
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
			SetAssertsDetails(strconv.Itoa(l.Id))


}

//SetMarketPriceOnline 计算市场数据实时,顺便把市场价格表也做出来，返回map序列化后添加进入redis key为eggProportion:时间戳
func SetMarketPriceOnline(data *model.ResponseDataList)  {
	list := make([]float64, 0, len(data.List))
	for _, v := range data.List {
		fixedPrice, FErr := strconv.ParseFloat(v.FixedPrice, 64)
		if FErr != nil {
			logger.Error(FErr)
			return
		}
		count := float64(v.Count)
		if v.Count != 1 {
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
	//切割数据中的name作为key
	strName := data.List[0].Name
	//categoryName := utils.Split(strName," ")
	err := redis.CreateKey(fmt.Sprintf("%sMarkerPrice",strName), list1[0])
	logger.Info(fmt.Sprintf("%sMarkerPrice",strName))
	if err != nil {
		logger.Error(err)
		return
	}
	//存进mysql
	data1 := model.MarketData{
		MarketName:strName,
		MarketData: list1[0],
	}
	mysql.InsertMarketPrice(data1)
}

