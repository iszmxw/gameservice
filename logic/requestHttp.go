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

//RequestAssertsData 请求获取数据
func RequestAssertsData(pageSize int,category int) *model.ResponseDataList {
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


//RequestAssertsDetails 访问获取一条详情数据，并且存入数据库
func RequestAssertsDetails(gid string)  {
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
	//存进redis
	key := fmt.Sprintf("%s:%s",assertDetails.Name,assertDetails.Gid)
	//序列化后存入
	marshal, Merr := json.Marshal(&assertDetails)
	if Merr != nil {
		logger.Error(Merr)
		return
	}
	Cerr := redis.CreateKey(key,string(marshal))
	if Cerr != nil {
		logger.Error(Cerr)
		return
	}
	mysql.CreateOneAssertDetails(assertDetails)
	if len(assertDetails.Status) >0  {
		//关联链上数据
		go RequestChainData(assertDetails.Gid)
	}else {
		logger.Info(assertDetails)
	}

}

