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

//CreateOneOrder  插入买入清单
func CreateOneOrder(data model.BabyOrder)  {
	err := mysql.DB.Model(model.BabyOrder{}).Create(&data).Error
	if err != nil{
		logger.Error(err)
		return
	}
}
//UpdateOneOrder 修改买入清单
func UpdateOneOrder(data model.BabyOrder)  {
	err := mysql.DB.Debug().Model(model.BabyOrder{}).Where("id",data.Id).Updates(&data).Error
	if err != nil{
		logger.Error(err)
		return
	}
}
//GetAllOrder 查询买入清单
func GetAllOrder()  []model.BabyOrder  {
	var data []model.BabyOrder
	err := mysql.DB.Model(model.BabyOrder{}).Find(&data).Error
	if err != nil{
		logger.Error(err)
		return nil
	}
	return data
}

//DeleteOneOrder 删除买入清单
func DeleteOneOrder(data []model.BabyOrder) bool{
	err := mysql.DB.Model(model.BabyOrder{}).Delete(&data).Error
	if err != nil{
		logger.Error(err)
		return false
	}
	return true
}