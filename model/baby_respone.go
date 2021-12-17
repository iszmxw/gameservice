/**
 @author:way
 @date:2021/12/16
 @note
**/

package model

import (
	"gorm.io/gorm"
	"redisData/pkg/helpers"
)

//baby游戏相关---------------------------------------------------

//RespBNTxList 访问币安网拿到交易列表
type RespBNTxList struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		BlockNumber       string `json:"blockNumber"`
		TimeStamp         string `json:"timeStamp"`
		Hash              string `json:"hash"`
		Nonce             string `json:"nonce"`
		BlockHash         string `json:"blockHash"`
		TransactionIndex  string `json:"transactionIndex"`
		From              string `json:"from"`
		To                string `json:"to"`
		Value             string `json:"value"`
		Gas               string `json:"gas"`
		GasPrice          string `json:"gasPrice"`
		IsError           string `json:"isError"`
		TxreceiptStatus   string `json:"txreceipt_status"`
		Input             string `json:"input"`
		ContractAddress   string `json:"contractAddress"`
		CumulativeGasUsed string `json:"cumulativeGasUsed"`
		GasUsed           string `json:"gasUsed"`
		Confirmations     string `json:"confirmations"`
	} `json:"result"`
}

//RespTxHashStatus 响应交易的状态数据
type RespTxHashStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		Status string `json:"status"`
	} `json:"result"`
}

//RespTxDetails 响应交易详情数据
type RespTxDetails struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  struct {
		BlockHash        string `json:"blockHash"`
		BlockNumber      string `json:"blockNumber"`
		From             string `json:"from"`
		Gas              string `json:"gas"`
		GasPrice         string `json:"gasPrice"`
		Hash             string `json:"hash"`
		Input            string `json:"input"`
		Nonce            string `json:"nonce"`
		To               string `json:"to"`
		TransactionIndex string `json:"transactionIndex"`
		Value            string `json:"value"`
		Type             string `json:"type"`
		V                string `json:"v"`
		R                string `json:"r"`
		S                string `json:"s"`
	} `json:"result"`
}

// RespBabyOrder 响应买卖数据
type RespBabyOrder struct {
	Id          int                `json:"id"`
	Name        string             `json:"name"`         //名称
	FixPrice    float64            `json:"fix_price"`    //单价
	SalePrice   float64            `json:"sale_price"`   //出售价格
	Profit      float64            `json:"profit"`       //利润
	Status      int                `json:"status"`       //状态 1.买入
	TokenId     string             `json:"token_id"`     //token
	MarketPrice float64            `json:"market_price"` //买入市场价
	CreatedAt   helpers.TimeNormal `json:"created_at"`
	UpdatedAt   helpers.TimeNormal `json:"updated_at"`
	DeletedAt   gorm.DeletedAt     `json:"deleted_at"`
	Buycount    int                `json:"buyCount"` //购买次数，大于3次，放弃这个单
}
