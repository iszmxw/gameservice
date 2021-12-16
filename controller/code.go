/**
 @author:way
 @date:2021/12/16
 @note
**/

package controller

type Recode int64

const (
	CodeSuccess Recode = 200
	CodeInvalidParam Recode = 500
	CodeServerBusy Recode = 501
)

var codeMsgMap = map[Recode]string{
	CodeSuccess:          "ok",
	CodeInvalidParam:     "请求参数错误",
	CodeServerBusy:       "服务器繁忙",

}

func (c Recode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}

