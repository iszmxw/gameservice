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

//BabyTxData 链上交易列表,数据由api返回
type BabyTxData struct {
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
	Token             string             `json:"token"`
	Price             string             `json:"price"`
	CreatedAt         helpers.TimeNormal `json:"created_at"`
	UpdatedAt         helpers.TimeNormal `json:"updated_at"`
	DeletedAt         gorm.DeletedAt
}
func (BabyTxData) TableName() string {
	return "baby_tx_data"
}

//BabyMarketPrice babyMarket的市场数据
type BabyMarketPrice struct {
	Id         int
	MarketName string	//市场名字
	MarketData float64	//市场数据
	CreatedAt  helpers.TimeNormal
	UpdatedAt  helpers.TimeNormal
	DeletedAt  gorm.DeletedAt
}
func (BabyMarketPrice) TableName() string {
	return "baby_marketprice"
}

//BabyOrder baby买入清单
type BabyOrder struct {
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
}
func (BabyOrder) TableName() string {
	return "baby_order"
}
