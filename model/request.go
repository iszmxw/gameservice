/**
 @author:way
 @date:2021/11/26
 @note
**/

package model

// ParamRiskMng 两个上涨和下跌的
type ParamRiskMng struct {
	Situation     string  `json:"situation" form:"situation"`           //上涨还是下跌
	TimeLevel     int     `json:"time_level" form:"time_level"`         //控制时间级别
	Percentage    float64 `json:"percentage" form:"percentage"`         //控制成分风险的百分比
	OperationType int     `json:"operation_type" form:"operation_type"` // 1.停止脚本 2.钉钉警告 3.停止脚本通报钉钉
	Status        int     `json:"status" form:"status"`                 // 1是启动 2是停止
}

//ParamOnOff 设置买卖全自动开关
type ParamOnOff struct {
	CrlName string `json:"crl_name" form:"crl_name"` //买入总开关  //卖出总开关
	Super   string `json:"super" form:"super"`       // 1.关 2.开

}

//ParamTypeId 作用根据typeID获取游戏NFT名称
type ParamTypeId struct {
	TypeId int `json:"type_id" form:"type_id"`
}

//ParamBuyAndSaleSet 买出设置
type ParamBuyAndSaleSet struct {
	ProductId   string `json:"product_id" form:"product_id"`
	MarketPrice string `json:"market_price" form:"market_price"`
	Percent     string `json:"percent" form:"percent"`
	Status      string `json:"status" form:"status"`
	Types       string `json:"types" form:"types"`
}

//ParamSellingRate 卖出率参数
type ParamSellingRate struct {
	TimeLevel     int    `json:"time_level" form:"time_level"`
	Percent       string `json:"percent" form:"percent"`
	Status        string `json:"status" form:"status"`
	OperationType string `json:"operation_type" form:"operation_type"`
}

//ParamProportion 获取市场百分比参数
type ParamProportion struct {
	TypeId int `json:"type_id"`
}

type ParamBuy struct {
	GId          string `json:"gid" form:"gid"`
	IdInContract string `json:"id_in_contract" form:"id_in_contract"`
	TxHash       string `json:"tx_hash" form:"tx_hash"`
}
