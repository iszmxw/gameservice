/**
 @author:way
 @date:2021/12/3
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

// GetDataInRedis 从redis中取出数据，运算使用
func GetDataInRedis(key string) (string, error) {
	data, err := redis.GetData(key)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return data, nil
}