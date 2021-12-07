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

// GetAssetName 通过type_id 获取type的名字
func GetAssetName(typeId int) *model.AssetsType {
	data := model.AssetsType{}
	err := mysql.DB.Model(model.AssetsType{}).Where("type_id",typeId).Find(&data).Error
	if err != nil{
		logger.Error(err)
		return nil
	}
	return &data

}

// GetBuyData 根据id判断买入卖出数据
func GetBuyData(t int) (data []model.Buy)  {
	if err := mysql.DB.Model(model.Buy{}).Where("type",t).Order("id desc").Limit(10).Find(&data).Error;err!=nil{
		logger.Error(err)
		return nil
	}
	return data
}

// GetBuyById  根据id判断买入卖出数据
func GetBuyById(Gid string) (data []model.Buy)  {
	if err := mysql.DB.Model(model.Buy{}).Where("gid",Gid).Limit(1).Find(&data).Error;err!=nil{
		logger.Error(err)
		return nil
	}
	return data
}


//GetAssetDetail100 或者最新100条产品详情数据
func GetAssetDetail100() []model.AssetsDetails {
	 data := make([]model.AssetsDetails,50)
	if err := mysql.DB.Model(model.AssetsDetails{}).Order("id desc").Limit(50).Find(&data).Error;err!=nil{
		logger.Error(err)
		return nil
	}
	return data
}


//GetAssetType 返回资产类型清单
func GetAssetType() (data []model.AssetsType ) {
	err := mysql.DB.Model(model.AssetsType{}).Find(&data).Error
	if err != nil{
		logger.Error(err)
		return
	}
	return data
}
