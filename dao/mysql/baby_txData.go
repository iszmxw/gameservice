/**
 @author:way
 @date:2021/12/16
 @note
**/

package mysql

import (
	"redisData/model"
	"redisData/pkg/logger"
	"redisData/pkg/mysql"
)

//CreateBNTxHashList 拆分链上的TX列表并且存入数据库
func CreateBNTxHashList(data model.BabyTxData)  {
	CreateErr := mysql.DB.Model(model.BabyTxData{}).Create(&data).Error
	if CreateErr != nil{
		logger.Error(CreateErr)
		return
	}
}

//GetPriceList 获取价格的列表
func GetPriceList()  []string{
	var price []string
	err := mysql.DB.Model(model.BabyTxData{}).Pluck("price",&price).Error
	if err != nil{
		logger.Error(err)
		return nil
	}
	return price
}

//GetAllBabyTx 获取全部买入清单
func GetAllBabyTx() (data []model.BabyTxData) {
	err := mysql.DB.Model(model.BabyTxData{}).Order("id desc").Find(&data).Error
	if err != nil{
		logger.Error(err)
		return nil
	}
	return data
}