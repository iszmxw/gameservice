/**
 @author:way
 @date:2021/11/30
 @note
**/

package mysql

import (
	"redisData/model"
	"redisData/pkg/logger"
	"redisData/pkg/mysql"
)

// CreateOneAssert 创建一条资产信息
func CreateOneAssert(data model.AssetsData)  {
	err := mysql.DB.Model(model.AssetsData{}).Create(&data).Error
	if err != nil{
		logger.Error(err)
		return
	}
}

// CreateOneAssertDetails 创建一条资产详情信息
func CreateOneAssertDetails(data model.AssetsDetails)  {
	err := mysql.DB.Model(model.AssetsDetails{}).Create(&data).Error
	if err != nil{
		logger.Error(err)
		return
	}
}

//GetDataByGid 通过gid查询资产详情
func GetDataByGid(gid string) *model.AssetsDetails {
	data := model.AssetsDetails{}
	err := mysql.DB.Model(model.AssetsDetails{}).Where("gid",gid).Find(&data).Error
	if err != nil{
		logger.Error(err)
		return nil
	}
	return &data


}
