/**
 @author:way
 @date:2021/12/3
 @note 存放http请求
**/

package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/model"
	"redisData/pkg/logger"
	"redisData/utils"
	"strconv"
	"strings"
	"time"
)


//元兽脚本相关请求

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

// RequestChainData 访问链上的数据
func RequestChainData(gid string) {
	//通过查询assets_details表的owner地址查询所在的交易
	data := mysql.GetDataByGid(gid)
	//通过接口访问owner交易的数据列表
	url := fmt.Sprintf("https://api.bscscan.com/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=1&offset=1000&sort=desc&apikey=%s", data.Owner, viper.GetString("apikey"))
	logger.Info(url)
	response, Gerr := http.Get(url) //这步访问可能会慢
	if Gerr != nil {
		logger.Error(Gerr)
	}
	body, _ := ioutil.ReadAll(response.Body)
	if len(string(body)) < 150 {
		logger.Info("休息30s")
		time.Sleep(30 * time.Second)
		return
	}
	//fmt.Println(body)
	//反序列化成结构体
	var d model.RespChainData
	UErr := json.Unmarshal(body, &d)
	logger.Info(d.Message)
	if UErr != nil {
		logger.Error(UErr)
	}
	//通过时间戳和时间筛选数据并且存进数据库,筛选数据
	for _, v := range d.Result {
		if v.TimeStamp != data.StartTime {
			continue
		}
		if strings.ToLower(v.To) != strings.ToLower(data.SaleAddress) {
			continue
		}
		logger.Info("符合条件")
		//符合条件就存进数据库
		chainData := model.ChainData{
			Gid:               gid,
			Blocknumber:       v.BlockNumber,
			Timestamp:         v.TimeStamp,
			Hash:              v.Hash,
			Nonce:             v.Nonce,
			Blockhash:         v.BlockHash,
			Transactionindex:  v.TransactionIndex,
			From:              v.From,
			To:                v.To,
			Value:             v.Value,
			Gas:               v.Gas,
			Gasprice:          v.GasPrice,
			Iserror:           v.IsError,
			TxreceiptStatus:   v.TxreceiptStatus,
			Input:             v.Input,
			Contractaddress:   v.ContractAddress,
			Cumulativegasused: v.CumulativeGasUsed,
			Gasused:           v.GasUsed,
			Confirmations:     v.Confirmations,
			PriceHax:          string([]byte(v.Input)[266:330]),
		}
		mysql.CreateChainData(chainData)
	}
}

//baby脚本相关请求----------------------------------------------------------------------------

//ReqBNTxList 访问币安获取交易列表数据
//https://api.bscscan.com/api?module=account&action=txlist&address=0xE97Fdca0A3Fc76b3046aE496C1502c9d8dFEf6fc&startblock=0&endblock=99999999&page=1&offset=10&sort=desc&apikey=DUKNN1QZMITSSZC61YINTD1CWQ92FWEKHM
func ReqBNTxList(address string,apikey string,sort string,offset string,page string,contain string)  {
	//校验参数
	if len(address) <= 0 {
		logger.Error(errors.New("输入合约地址为空"))
		return
	}
	if len(apikey) <= 0 {
		apikey = "DUKNN1QZMITSSZC61YINTD1CWQ92FWEKHM"
	}
	if len(sort) <= 0 {
		sort = "desc"
	}
	if len(offset) <= 0 {
		offset = "100"
	}
	if len(page) <= 0 {
		page = "1"
	}
	//逻辑
	url := fmt.Sprintf("https://api.bscscan.com/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=%s&offset=%s&sort=%s&apikey=%s",address,page,offset,sort,apikey)
	logger.Info(url)
	response, GErr := http.Get(url) //这步访问可能会慢
	if GErr != nil {
		logger.Error(GErr)
		return
	}
	//反序列化
	body, ReadAllErr := ioutil.ReadAll(response.Body)
	if ReadAllErr != nil{
		logger.Error(ReadAllErr)
		return
	}
	var resp model.RespBNTxList
	UnmarshalErr := json.Unmarshal(body,&resp)
	if UnmarshalErr != nil {
		logger.Error(UnmarshalErr)
		return
	}

	for _,v := range resp.Result{
		if strings.Contains(v.Input,contain){
			//去重
			if redis.ExistEle("baby:txSet",v.Hash){
				continue
			}
			//把交易hash存进redis Set中 ，去重使用
			redis.CreateSetData("baby:txSet",v.Hash)
			//拆分input数据
			strInput := []byte(v.Input)
			tokenID := utils.StringToBigInt(string(strInput[10:74]))
			price := utils.StringToBigInt(string(strInput[74:138]))
			//存2 份redis
			marshaldata, marshalERR := json.Marshal(v)
			if marshalERR != nil {
				logger.Error(marshalERR)
				return
			}
			CreateKeyExpireErr := redis.CreateKeyExpire(fmt.Sprintf("baby:txHash:%s",v.Hash),string(marshaldata),0)
			if CreateKeyExpireErr != nil {
				logger.Error(CreateKeyExpireErr)
				return
			}
			//存一份mysql
			data := model.BabyTxData{
				BlockNumber: v.BlockNumber,
				TimeStamp: v.TimeStamp,
				Hash: v.Hash,
				Nonce: v.Nonce,
				BlockHash: v.BlockHash,
				TransactionIndex: v.TransactionIndex,
				From: v.From,
				To: v.To,
				Value: v.Value,
				Gas: v.Gas,
				GasPrice: v.GasPrice,
				IsError: v.IsError,
				TxreceiptStatus: v.TxreceiptStatus,
				Input: v.Input,
				ContractAddress: v.ContractAddress,
				Confirmations: v.Confirmations,
				GasUsed: v.GasUsed,
				CumulativeGasUsed: v.CumulativeGasUsed,
				//拆分input数据，
				Token: tokenID.String(),
				Price: price.String(),
			}
			mysql.CreateBNTxHashList(data)
		}
	}
}

func ReqTxDetailByHash(txHash string) *model.RespTxDetails {
	url := fmt.Sprintf("https://api.bscscan.com/api?module=proxy&action=eth_getTransactionByHash&txhash=%s&apikey=DUKNN1QZMITSSZC61YINTD1CWQ92FWEKHM",txHash)
	logger.Info(url)
	response, GErr := http.Get(url) //这步访问可能会慢
	if GErr != nil {
		logger.Error(GErr)
		return nil
	}
	//反序列化
	body, ReadAllErr := ioutil.ReadAll(response.Body)
	if ReadAllErr != nil{
		logger.Error(ReadAllErr)
		return nil
	}
	var data model.RespTxDetails
	UnmarshalErr := json.Unmarshal(body,&data)
	if UnmarshalErr != nil {
		logger.Error(UnmarshalErr)
		return nil
	}
	return &data

}

//ReqGetTxStatus 请求baby相关交易的hash
func ReqGetTxStatus(txHash string) *model.RespTxHashStatus {
	url := fmt.Sprintf("https://api.bscscan.com/api?module=transaction&action=gettxreceiptstatus&txhash=%s&apikey=DUKNN1QZMITSSZC61YINTD1CWQ92FWEKHM",txHash)
	logger.Info(url)
	response, GErr := http.Get(url) //这步访问可能会慢
	if GErr != nil {
		logger.Error(GErr)
		return nil
	}
	//反序列化
	body, ReadAllErr := ioutil.ReadAll(response.Body)
	if ReadAllErr != nil{
		logger.Error(ReadAllErr)
		return nil
	}
	var data model.RespTxHashStatus
	UnmarshalErr := json.Unmarshal(body,&data)
	if UnmarshalErr != nil {
		logger.Error(UnmarshalErr)
		return nil
	}
	return &data
}
//baby脚本相关请求