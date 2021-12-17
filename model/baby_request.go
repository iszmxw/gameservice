/**
 @author:way
 @date:2021/12/16
 @note
**/

package model

//  买卖数据

type ParamsBuyAndSellQuery struct {
	Name    string `json:"name" form:"name"`         // 名称
	TokenId string `json:"token_id" form:"token_id"` // token
	Status  string `json:"status" form:"status"`     // 状态 1.买入
}
