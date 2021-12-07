package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"redisData/model"
	"redisData/pkg/logger"
	"time"
)

// 创建redis key
var (
	ErrorRedisDataIsNull = errors.New("name does not exist")
	ErrorGetDataFail     = errors.New("name does not exist")
)

//根据key获取值

// ExistKey 判断key是否存在
func ExistKey(key string) bool {
	result, err := rdb.Exists(key).Result()
	if err != nil {
		return true
	}
	if result == 1 {
		return true
	}
	if result == 0 {
		return false
	}
	return true
}
// 创建egg:id
func CreateEggData(key string, value interface{}) {
	fullKey := getEggData(key)
	err := rdb.Set(fullKey, value, 36000*time.Second).Err()
	//log.Println("redis finish create or change")
	if err != nil {
		log.Println(err)
	}
}
// 删除egg:id
func DeleEggKey(key string)  {
	fullKey := getEggData(key)
	rdb.Del(fullKey)
}
// 创建potion:id
func CreatePotionData(key string, value interface{}) {
	fullKey := getPotionData(key)
	err := rdb.Set(fullKey, value, 36000*time.Second).Err()
	//log.Println("redis finish create or change")
	if err != nil {
		log.Println(err)
	}
}

// GetDataByKey 输入入egg:id或者potion:id 返回对应的结构体
func GetDataByKey(key string) (model.RespAssetsDetailList,error) {
	var response model.RespAssetsDetailList
	res, err := rdb.Get(key).Result()
	fmt.Println(res)
	logger.Info("单个key的数据")
	if err != nil {
		if err == redis.Nil {
			logger.Error(err)
			return model.RespAssetsDetailList{},err
			log.Println("key does not exist")
		}
		logger.Error(err)
		return model.RespAssetsDetailList{},err
		log.Printf("get name failed, err:%v\n", err)
	}
	UnmarshalErr := json.Unmarshal([]byte(res), &response)
	if UnmarshalErr != nil {
		return model.RespAssetsDetailList{},UnmarshalErr
	} else {
		return response,nil
	}
}

func GetKeysByPfx(keypfx string) ([]string,error) {
	vals, err := rdb.Keys(fmt.Sprintf("%s*", keypfx)).Result()
	logger.Info("vals")
	logger.Info(vals)
	if err != nil {
		logger.Error(err)
		return nil,err
	}
	return vals,nil
}

func CreateKey(key string, value interface{}) error {
	err := rdb.Set(key, value, 3600*12*time.Second).Err()
	//log.Println("redis finish create or change")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func CreateDurableKey(key string, value interface{}) error {
	err := rdb.Set(key, value,-1).Err()
	//log.Println("redis finish create or change")
	if err != nil {
		log.Println(err)
		return err
	}
	logger.Info("创建key")
	logger.Info(key)
	return nil

}

func GetData(key string) (data string,err error) {
	logger.Info(key)
	res, err := rdb.Get(key).Result()
	if err != nil{
		logger.Error(err)
		if err == redis.Nil{
			logger.Error(err)
			return "",ErrorRedisDataIsNull
		}
		return "",nil
	}
	return res,nil
}


// set相关操作

// CreateSetData 向一个集合中存值
func CreateSetData(key string,value string)  {
	rdb.SAdd(key,value)
}
//DeleteSetData  从集合中移除值
func DeleteSetData(key string,value string)  {
	rdb.SRem(key,value)
}
// ExistEle 判断集合中是否存在某个值
func ExistEle(key string,value string) bool {
	return rdb.SIsMember(key,value).Val()
}



//hash相关操作

// CreatHashKey 创建hash的key
func CreatHashKey(key string,m map[string]interface{})  {
	rdb.HMSet(key,m)
}
// GetHashDataAll 根据key读hash中的全部数据
func GetHashDataAll(key string) map[string]string {
	result, err := rdb.HGetAll(key).Result()
	if err != nil {
		logger.Error(err)
		return nil
	}
	return result
}

//zset相关操作

// CreateZScoreData 创建一个有序集合
func CreateZScoreData(key string,member string,score float64)  {
	rdb.ZAdd(key, redis.Z{
		Score: score,
		Member: member,
	})
}
//遍历有序集合

//GetScoreByMember 在有序集合中根据member查询对应的score
func GetScoreByMember(key string,member string) interface{}  {
	f:=rdb.ZScore(key,member).Val()
	return f
}

//DeleteRecByMember 根据member删除集合中的某个数据
func DeleteRecByMember(key string,member string)  {
	fmt.Printf("删除menber:%s",member)
	rdb.ZRem(key,member).Val()
}

//GetAllZSet  遍及集合
func GetAllZSet(key string) []string {
	strSlice := rdb.ZRevRange(key,0,-1).Val()
	return strSlice
}


//list相关操作

// SetOneList 添加一个list
func SetOneList(key string,value string)  {
	rdb.LPush(key,value)
}

// GetAllList 按照先进先出的顺序遍历list
func GetAllList(key string) []string {
	strSilce := rdb.LRange(key,0,-1).Val()
	return strSilce
}

// RmListEle 移除list某个值
func RmListEle(key string,value string )  {
	rdb.LRem(key,1,value)
}

// RmListHead 先进先出逻辑
func RmListHead(key string)  string{
	val := rdb.RPop(key).Val()
	return val
}









