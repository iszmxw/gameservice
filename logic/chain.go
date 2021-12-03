/**
 @author:way
 @date:2021/12/2
 @和链相关的结构体
**/

package logic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"redisData/dao/mysql"
	"redisData/model"
	"redisData/pkg/logger"
	"strings"
	"time"
)

func RequestChainData(gid string) {
	//通过查询assets_details表的owner地址查询所在的交易
	data := mysql.GetDataByGid(gid)
	//通过接口访问owner交易的数据列表
	url := fmt.Sprintf("https://api.bscscan.com/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=1&offset=1000&sort=desc&apikey=DUKNN1QZMITSSZC61YINTD1CWQ92FWEKHM",data.Owner)
	logger.Info(url)
	response, Gerr := http.Get(url)  //这步访问可能会慢
	if Gerr != nil {
		logger.Error(Gerr)
	}
	body, _ := ioutil.ReadAll(response.Body)
	if len(string(body)) <150 {
		logger.Info("休息30s")
		time.Sleep(30*time.Second)
		return
	}
	//fmt.Println(body)
	//反序列化成结构体
	var d model.RespChainData
	UErr := json.Unmarshal(body,&d)
	logger.Info(d.Message)
	if UErr != nil {
		logger.Error(UErr)
	}
	//通过时间戳和时间筛选数据并且存进数据库,筛选数据
	for _,v := range d.Result{
		if v.TimeStamp != data.StartTime{
			continue
		}
		if strings.ToLower(v.To) !=strings.ToLower(data.SaleAddress) {
			continue
		}
		logger.Info("符合条件")
		//符合条件就存进数据库
		chainData := model.ChainData{
			Gid: gid,
			Blocknumber: v.BlockNumber,
			Timestamp: v.TimeStamp,
			Hash: v.Hash,
			Nonce: v.Nonce,
			Blockhash: v.BlockHash,
			Transactionindex: v.TransactionIndex,
			From: v.From,
			To: v.To,
			Value: v.Value,
			Gas: v.Gas,
			Gasprice: v.GasPrice,
			Iserror: v.IsError,
			TxreceiptStatus: v.TxreceiptStatus,
			Input: v.Input,
			Contractaddress: v.ContractAddress,
			Cumulativegasused: v.CumulativeGasUsed,
			Gasused: v.GasUsed,
			Confirmations: v.Confirmations,
			PriceHax: string([]byte(v.Input)[266:330]),
		}
		mysql.CreateChainData(chainData)
	}
}
