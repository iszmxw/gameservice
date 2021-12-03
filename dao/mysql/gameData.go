/**
 @author:way
 @date:2021/11/26
 @note
**/

package mysql

import (
	"redisData/model"
	"redisData/pkg/logger"
	"redisData/pkg/mysql"
)

func InsertMarketPrice(data model.MarketData)  {
	err := mysql.DB.Model(model.MarketData{}).Create(&data).Error
	if err != nil{
		logger.Error(err)
		return
	}
}

func InsertBuyRecord(data model.Buy)  {
	err := mysql.DB.Model(model.Buy{}).Create(&data).Error
	if err != nil{
		logger.Error(err)
		return
	}
}

func GetBuyData(t int) (data []model.RespBuy)  {
	if err := mysql.DB.Model(model.Buy{}).Where("type",t).Order("desc").Find(&data).Error;err!=nil{
		logger.Error(err)
		return nil
	}
	return data
}