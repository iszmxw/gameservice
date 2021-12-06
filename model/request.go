/**
 @author:way
 @date:2021/11/26
 @note
**/

package model

type ParamGetData struct {
	DataType string `json:"data_type" form:"data_type"`
}

type ParamStart struct {
	Buy  float64 `json:"buy" form:"buy"`
	Sale float64 `json:"sale" form:"sale"`
	Safe float64 `json:"safe" form:"safe"`
}

type ParamGetBuy struct {
	Type int `json:"type"`
}

// ParamRiskMng 两个上涨和下跌的
type ParamRiskMng struct {
	Situation     string     `json:"situation" form:"situation"`           //上涨还是下跌
	TimeLevel     int     `json:"time_level" form:"time_level"`         //控制时间级别
	Percentage    float64 `json:"percentage" form:"percentage"`         //控制成分风险的百分比
	OperationType int     `json:"operation_type" form:"operation_type"` // 1.停止脚本 2.钉钉警告 3.停止脚本通报钉钉
	Status        int     `json:"status" form:"status"`                 // 1是启动 2是停止
}

type ParamBuyAndSale struct {
	ProductName    string  `json:"product_name" form:"product_name"`       //产品名称
	RisePercentage float64 `json:"rise_percentage" form:"rise_percentage"` //上涨百分比
	FallPercentage float64 `json:"fall_percentage" form:"fall_percentage"` // 下跌百分比
	Status         int     `json:"status" form:"status"`                   // 1.开  2.关
}

type ParamOnOff struct {
	CrlName string  `json:"crl_name" form:"crl_name"`  //买入总开关  //卖出总开关
	Super int `json:"super" form:"super"` // 1.关 2.开

}

type ParamTypeId struct {
	TypeId int `json:"type_id" form:"type_id"`
}

type ParamBuyStatus struct {
	Status int `json:"status" form:"status"`
}
