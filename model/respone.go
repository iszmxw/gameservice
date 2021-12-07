/**
 @author:way
 @date:2021/11/26
 @note
**/

package model

import (
	"redisData/pkg/helpers"
)

//Rdata 返回给前端的资产详情
type Rdata struct {
	Data struct {
		Id             int           `json:"id"`
		Name           string        `json:"name"`
		Description    string        `json:"description"`
		CreatedAt      string        `json:"created_at"`
		CategoryName   string        `json:"category_name"`
		CollectionName string        `json:"collection_name"`
		ImageUrl       string        `json:"image_url"`
		Count          int           `json:"count"`
		FixedPrice     string        `json:"fixed_price"`
		TotalPrice     string        `json:"total_price"`
		HighestPrice   string        `json:"highest_price"`
		SaleAddress    string        `json:"sale_address"`
		SaleType       string        `json:"saleType"`
		IdInContract   string        `json:"id_in_contract"`
		TokenId        int           `json:"token_id"`
		TokenStandard  string        `json:"token_standard"`
		TokenType      int           `json:"token_type"`
		Owner          string        `json:"owner"`
		NftAddress     string        `json:"nft_address"`
		BlockChain     string        `json:"block_chain"`
		ProtocolFee    string        `json:"protocol_fee"`
		NeedCheck      int           `json:"need_check"`
		ShowLeft       int           `json:"show_left"`
		StartTime      int           `json:"start_time"`
		EndTime        int           `json:"end_time"`
		Status         string        `json:"status"`
		Properties     []interface{} `json:"properties"`
	} `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// ResponseDataList 返回资产列表数据
type ResponseDataList struct {
	Total int    `json:"total"`
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	List  []struct {
		Id           int    `json:"id"`
		Name         string `json:"name"`
		ImageUrl     string `json:"image_url"`
		Count        int    `json:"count"`
		FixedPrice   string `json:"fixed_price"`
		HighestPrice string `json:"highest_price"`
		Status       string `json:"status"`
		SaleType     string `json:"sale_type"`
		TokenId      string `json:"token_id"`
		SaleAddress  string `json:"sale_address"`
	} `json:"list"`
}

// RespBuy buy表对应的结构体
type RespBuy struct {
	Gid         string             `json:"gid"`
	Name        string             `json:"name"`
	Count       int                `json:"count"`
	FixedPrice  float64            `json:"fixed_price"`
	TotalPrice  float64            `json:"total_price"`
	Type        int                `json:"type"`
	SaleAddress string             `json:"sale_address"`
	TokenId     string             `json:"token_id"`
	MarketPrice float64            `json:"market_price"`
	Profit      float64            `json:"profit"`
	CreateTime  helpers.TimeNormal `json:"create_time"`
}

// ResponseAssertsDetails 根据ID返回的资产详情信息
type ResponseAssertsDetails struct {
	Data struct {
		Id            int           `json:"id"`
		Name          string        `json:"name"`
		Description   string        `json:"description"`
		CreatedAt     string        `json:"created_at"`
		ImageUrl      string        `json:"image_url"`
		Count         int           `json:"count"`
		FixedPrice    string        `json:"fixed_price"`
		TotalPrice    string        `json:"total_price"`
		SaleAddress   string        `json:"sale_address"`
		IdInContract  string        `json:"id_in_contract"`
		TokenId       int           `json:"token_id"`
		TokenStandard string        `json:"token_standard"`
		Owner         string        `json:"owner"`
		NftAddress    string        `json:"nft_address"`
		BlockChain    string        `json:"block_chain"`
		StartTime     int           `json:"start_time"`
		Status        string        `json:"status"`
		Properties    []interface{} `json:"properties"`
	} `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// RespChainData 关联链上数据
type RespChainData struct {
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

//RespAssetsDetailList 返回资产列表详情
type RespAssetsDetailList struct {
	Gid           string `json:"Gid"`
	Name          string `json:"Name"`
	Description   string `json:"Description"`
	ImageUrl      string `json:"ImageUrl"`
	Count         int    `json:"Count"`
	FixedPrice    string `json:"FixedPrice"`
	TotalPrice    string `json:"TotalPrice"`
	SaleAddress   string `json:"SaleAddress"`
	IdInContract  string `json:"IdInContract"`
	TokenId       string `json:"TokenId"`
	TokenStandard string `json:"TokenStandard"`
	Owner         string `json:"Owner"`
	NftAddress    string `json:"NftAddress"`
	BlockChain    string `json:"BlockChain"`
	StartTime     string `json:"StartTime"`
	Status        string `json:"Status"`
	Properties    string `json:"Properties"`
	CreatedAt     string `json:"CreatedAt"`
}

//RespBuyAndSale 买卖控制参数
type RespBuyAndSale struct {
	FallPercentage string `json:"fall_percentage"`
	ProductName    string `json:"product_name"`
	RisePercentage string `json:"rise_percentage"`
	RiseStatus     string `json:"rise_status"`
	FallStatus     string `json:"fall_status"`
}

// RespAllOnOff 高级开关
type RespAllOnOff struct {
	CrlName string `json:"crl_name"`
	Super   string `json:"super"`
}

type RespAllSwitch struct {
	BuyAndSale []RespBuyAndSaleSet `json:"buy_and_sale"`
	AllOnOff   []RespAllOnOff      `json:"all_on_off"`
}

type RespRiskMonitor struct {
	Situation     string `json:"situation"`
	TimeLevel     string `json:"time_level"`
	Percentage    string `json:"percentage"`
	OperationType string `json:"operation_type"`
	Status        string `json:"status"`
}

type RespAssetType struct {
	TypeName string `json:"type_name"`
	TypeID   int    `json:"type_id"`
}

//RespBuyAndSaleSet 买出设置
type RespBuyAndSaleSet struct {
	Percent   string `json:"percent"`
	ProductId string `json:"product_id"`
	Status    string `json:"status"`
	Types     string `json:"types"`
}
