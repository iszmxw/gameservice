/**
 @author:way
 @date:2021/11/26
 @note
**/

package model

import (
	"gorm.io/gorm"
	"redisData/pkg/helpers"
)

//Data 和 List是次请求中的数据
type Data struct {
	Total int    `json:"total"`
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	List  []List `json:"list"`
}

type List struct {
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
}

//MarketData 市场数据结构体
type MarketData struct {
	Id         int
	MarketName string
	MarketData float64
	CreatedAt  helpers.TimeNormal
	UpdatedAt  helpers.TimeNormal
	DeletedAt  gorm.DeletedAt
}

//Buy 买入卖出结构体
type Buy struct {
	Id           int     `json:"id" `
	Gid          string  `json:"gid" `
	Name         string  `json:"name" `
	Count        int     `json:"count" `
	FixedPrice   float64 `json:"fixed_price" `
	Type         int     `json:"type" `
	SaleAddress  string  `json:"sale_address" `
	TokenId      string  `json:"token_id" `
	MarketPrice  float64 `json:"market_price" `
	Profit       float64 `json:"profit" `
	SalePrice    float64 `json:"sale_price" `
	IdInContract string  `json:"id_in_contract" `
	TxHash       string  `json:"tx_hash" `
	CreatedAt    helpers.TimeNormal
	UpdatedAt    helpers.TimeNormal
	DeletedAt    gorm.DeletedAt
}

// AssetsData 游戏资产的数据egg和potion
type AssetsData struct {
	GId          string `gorm:"column:gid"`
	Name         string
	ImageUrl     string
	Count        int
	FixedPrice   string
	HighestPrice string
	Status       string
	SaleType     string
	TokenId      string
	SaleAddress  string
	CreatedAt    helpers.TimeNormal
	UpdatedAt    helpers.TimeNormal
	DeletedAt    gorm.DeletedAt
}

//AssetsDetails 游戏资产详情数据
type AssetsDetails struct {
	Gid           string `gorm:"column:gid"`
	Name          string
	Description   string
	ImageUrl      string
	Count         int
	FixedPrice    string
	TotalPrice    string
	SaleAddress   string
	IdInContract  string
	TokenId       string
	TokenStandard string
	Owner         string
	NftAddress    string
	BlockChain    string
	StartTime     string
	Status        string
	Properties    string
	CreatedAt     string
}

//ChainData 资产链上数据
type ChainData struct {
	Gid               string
	PriceHax          string
	Id                uint
	Blocknumber       string
	Timestamp         string
	Hash              string
	Nonce             string
	Blockhash         string
	Transactionindex  string
	From              string
	To                string
	Value             string
	Gas               string
	Gasprice          string
	Iserror           string
	TxreceiptStatus   string
	Input             string
	Contractaddress   string
	Cumulativegasused string
	Gasused           string
	Confirmations     string
}

//AssetsType 资产类型
type AssetsType struct {
	Id       int
	TypeName string
	TypeId   int
}

//ChainTxData 链上交易列表
type ChainTxData struct {
	Id                int64              `json:"id"`
	BlockNumber       string             `json:"block_number"`
	TimeStamp         string             `json:"time_stamp"`
	Hash              string             `json:"hash"`
	Nonce             string             `json:"nonce"`
	BlockHash         string             `json:"block_hash"`
	TransactionIndex  string             `json:"transaction_index"`
	From              string             `json:"from"`
	To                string             `json:"to"`
	Value             string             `json:"value"`
	Gas               string             `json:"gas"`
	GasPrice          string             `json:"gas_price"`
	IsError           string             `json:"is_error"`
	TxreceiptStatus   string             `json:"txreceipt_status"`
	Input             string             `json:"input"`
	ContractAddress   string             `json:"contract_address"`
	CumulativeGasUsed string             `json:"cumulative_gas_used"`
	GasUsed           string             `json:"gas_used"`
	Confirmations     string             `json:"confirmations"`
	CreatedAt         helpers.TimeNormal `json:"created_at"`
	UpdatedAt         helpers.TimeNormal
	DeletedAt         gorm.DeletedAt
}
