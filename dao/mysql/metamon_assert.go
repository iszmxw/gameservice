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
	logger.Info(data)
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

//GetBuyCount 获取买卖清单
func GetBuyCount(startTime string, endTime string) float64 {
	var result []float64
	var sum float64
	mysql.DB.Table("buy").Where("type=1  and created_at between ? and ?",startTime,endTime).Pluck("count",&result )
	for _,v := range result{
		sum += v
	}
	return sum
}

//GetSaleCount 获取买出数量
func GetSaleCount(startTime string, endTime string) float64 {
	var result []float64
	var sum float64
	mysql.DB.Table("buy").Where("type=2  and updated_at between ? and ?",startTime,endTime).Pluck("count",&result )
	for _,v := range result{
		sum += v
	}
	return sum

}

//GetBuySaleCountList 返回买出卖出的数量列表
func GetBuySaleCountList(startTime string,t int ) []float64{
	var count []float64
	mysql.DB.Raw("SELECT count FROM buy WHERE type = ? and created_at > ?",t ,startTime).Scan(&count)
	return count
}

//UpdateBuy 更新buy表,发起交易后返回，更新字段,目前只需要更新hash和状态
func UpdateBuy(data model.Buy){
	err := mysql.DB.Model(model.Buy{}).Where("gid",data.Gid).Updates(&data).Error
	if err != nil{
		logger.Error(err)
		return
	}
}

//GetBuyDataHashNotNull 获取买入数据，条件hash不等于nil
func GetBuyDataHashNotNull() *[]model.Buy {
	var b []model.Buy
	err := mysql.DB.Model(model.Buy{}).Where("type",1).Not("tx_hash","").Find(&b).Error
	if err != nil{
		logger.Error(err)
		return nil
	}
	return &b
}
//GetSaleDataHashNotNull 获取买入数据，条件hash不等于nil
func GetSaleDataHashNotNull()  *[]model.Buy{
	var b []model.Buy
	err := mysql.DB.Model(model.Buy{}).Where("type",2).Not("tx_hash","").Find(&b).Error
	if err != nil{
		logger.Error(err)
		return nil
	}
	return &b
}


