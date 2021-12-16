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

//InsertBabyMarketPrice 添加一条市场价格
func InsertBabyMarketPrice(data model.BabyMarketPrice)  {
	err := mysql.DB.Model(model.BabyMarketPrice{}).Create(&data).Error
	if err != nil{
		logger.Error(err)
		return
	}
}