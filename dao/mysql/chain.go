/**
 @author:way
 @date:2021/12/2
 @note
**/

package mysql

import (
	"redisData/model"
	"redisData/pkg/logger"
	"redisData/pkg/mysql"
)

//CreateChainData 创建一条和链相关的数据
func CreateChainData(data model.ChainData)  {
	err := mysql.DB.Model(model.ChainData{}).Create(&data).Error
	if err != nil{
		logger.Error(err)
		return
	}
	return
}

