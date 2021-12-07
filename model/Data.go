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
	Id int
	MarketName string
	MarketData float64
	CreatedAt helpers.TimeNormal
	UpdatedAt helpers.TimeNormal
	DeletedAt gorm.DeletedAt
}

//Buy 买入卖出结构体
type Buy struct {
	Id int
	Gid string
	Name string
	Count int
	FixedPrice float64
	Type int
	SaleAddress string
	TokenId string
	MarketPrice float64
	Profit float64
	SalePrice float64
	CreatedAt helpers.TimeNormal
	UpdatedAt helpers.TimeNormal
	DeletedAt gorm.DeletedAt
}

// AssetsData 游戏资产的数据egg和potion
type AssetsData struct {
	GId string `gorm:"column:gid"`
	Name string
	ImageUrl string
	Count int
	FixedPrice string
	HighestPrice string
	Status string
	SaleType string
	TokenId string
	SaleAddress string
	CreatedAt helpers.TimeNormal
	UpdatedAt helpers.TimeNormal
	DeletedAt gorm.DeletedAt
}

//AssetsDetails 游戏资产详情数据
type AssetsDetails struct {
	Gid string `gorm:"column:gid"`
	Name string
	Description string
	ImageUrl string
	Count int
	FixedPrice string
	TotalPrice string
	SaleAddress string
	IdInContract string
	TokenId string
	TokenStandard string
	Owner string
	NftAddress string
	BlockChain string
	StartTime string
	Status string
	Properties string
	CreatedAt string
}

//ChainData 资产链上数据
type ChainData struct {
	Gid string
	PriceHax string
	Id uint
	Blocknumber string
	Timestamp string
	Hash string
	Nonce string
	Blockhash string
	Transactionindex string
	From string
	To string
	Value string
	Gas string
	Gasprice string
	Iserror string
	TxreceiptStatus string
	Input string
	Contractaddress string
	Cumulativegasused string
	Gasused string
	Confirmations string
}

//AssetsType 资产类型
type AssetsType struct {
	Id int
	TypeName string
	TypeId int
}



