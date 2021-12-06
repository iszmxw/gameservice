/**
 @author:way
 @date:2021/11/30
 @note
**/

package logic

import (
	"encoding/json"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/model"
	"redisData/pkg/logger"
	"strconv"
)



//ManageData 数据处理逻辑1.判断数据是否存在集合中 是就跳过，不是存在list里面
func ManageData(data *model.ResponseDataList){
	if data == nil{
		logger.Info("数据为空")
		return
	}

	for _,v := range data.List{

		if redis.ExistEle("assertSet",strconv.Itoa(v.Id)){
			continue
		}
		//如果不是把数据存进redis队列
		marshal, Merr := json.Marshal(v)
		if Merr != nil {
			logger.Error(Merr)
			return 
		}
		//使用list储存
		redis.SetOneList("assertList",string(marshal))
		//使用集合把ID储存起来
		redis.CreateSetData("assertSet",strconv.Itoa(v.Id))

	}
}

//StoreListToMysql 把redis中队列中的数据储存到mysql
func StoreListToMysql(str string)  {
		logger.Info(str)
		l :=  model.List{}
		err := json.Unmarshal([]byte(str), &l)
		if err != nil {
			logger.Error(err)
			return
		}
		d := model.AssetsData{
			GId: strconv.Itoa(l.Id),
			Name: l.Name,
			FixedPrice: l.FixedPrice,
			HighestPrice: l.HighestPrice,
			ImageUrl: l.ImageUrl,
			Count: l.Count,
			SaleType: l.SaleType,
			TokenId: l.TokenId,
			SaleAddress: l.SaleAddress,
			Status: l.Status,
		}
		mysql.CreateOneAssert(d)
		RequestAssertsDetails(strconv.Itoa(l.Id))


}





