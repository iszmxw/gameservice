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
	Buy float64 `json:"buy" form:"buy"`
	Sale float64 `json:"sale" form:"sale"`
	Safe float64 `json:"safe" form:"safe"`
}

type ParamGetBuy struct {
	Type int `json:"type"`
}

