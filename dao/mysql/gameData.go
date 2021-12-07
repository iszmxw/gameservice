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

// GetHistoryMarketData 输入对应的秒数返回对应的历史市场数据
func GetHistoryMarketData(second int,types string) (data *model.MarketData) {
	err := mysql.DB.Model(model.MarketData{}).Where("market_name",types).Order("id desc").Offset(second-1).Limit(1).Find(&data).Error
	if err != nil{
		logger.Error(err)
		return nil
	}
	return
}