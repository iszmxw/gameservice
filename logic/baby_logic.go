/**
 @author:way
 @date:2021/12/15
 @note 存放baby脚本相关的逻辑
**/

package logic

import (
	"fmt"
	"redisData/dao/redis"
	"strconv"
)

func CountBabyMarPrice(priceList []float64) []float64 {
	m1 := make(map[float64]int)
	var s2 []int
	var max int
	var s3 []float64

	// 统计频率最高的价格
	for _, v := range priceList {
		if m1[v] != 0 {
			m1[v]++
		} else {
			m1[v] = 1
		}
	}
	//遍历m1把里面的float转化成string
	for _, v := range m1 {
		//拼接成数组
		s2 = append(s2, v)
	}
	// 取出来放进数组

	for _, v := range m1 {
		//拼接成数组
		s2 = append(s2, v)
	}
	//算出最大值
	if s2 == nil{
		return nil
	}
	max = s2[0]
	for i := 0; i < len(s2); i++ {
		if max < s2[i] {
			max = s2[i]
		}
	}

	//存在出现同样次数的
	for k, v := range m1 {
		if v == max {
			s3 = append(s3, k)
		}
	}
	//插入一条redis数据，把这次遍历市场价占比计算后返回
		m2 := make(map[string]interface{})
		for i,v := range m1{
			str := strconv.FormatFloat(i, 'E', -1, 64)
			m2[fmt.Sprintf("%s",str)] = v
		}
		redis.CreatHashKey(fmt.Sprintf("baby:Proportion"),m2)
	return s3
}