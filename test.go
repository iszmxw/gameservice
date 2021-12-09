/**
 @author:way
 @date:2021/12/2
 @note
**/

package main

import (
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/pkg/logger"
	"redisData/setting"
)

func init() {
	// 定义日志目录
	logger.Init("redisData")
	// 初始化 viper 配置
	if err := setting.Init(""); err != nil {
		logger.Info("viper init fail")
		logger.Error(err)
		return
	}
	mysql.InitMysql()
	//初始化redis
	if err := redis.InitClient(); err != nil {
		logger.Info("init redis fail err")
		logger.Error(err)
		return
	}

} //func Hex2Dec(val string) int {
//	n, err := strconv.ParseUint(val, 16, 64)
//	if err != nil {
//		fmt.Println(err)
//	}
//	return int(n)
//}
//
//func main() {
//	hex := "00000000000000000000000000000000000000000000a2bc77ee287ecf500000"
//	dec := Hex2Dec(hex)
//	fmt.Println(dec)
//}

//00000000000000000000000000000000000000000000131454ae75bda7500000

//func main() {
//	str := "0x467f963d000000000000000000000000d40c03b8680d4b6a4d78fc3c6f6a28c854e94a790000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000012bb890508c125661e03b09ec06e404bc928904000000000000000000000000000000000000000000000131454ae75bda750000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
//	str2 := []byte(str)
//	fmt.Println(string(str2[74+64+64+64:74+64+64+64+64]))
//	str3 := "00000000000000000000000000000000000000000000131454ae75bda7500000"
//	if str3 == string(str2[74+64+64+64:74+64+64+64+64]){
//		println("1111")
//	}
//}
//
//func main() {
//	str := "0x467f963d000000000000000000000000d40c03b8680d4b6a4d78fc3c6f6a28c854e94a790000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000012bb890508c125661e03b09ec06e404bc92890400000000000000000000000000000000000000000000013451eb0c55622e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
//	str2 := []byte(str)
//	str3 := "0000000000000000000000000000000000000000000013451eb0c55622e00000"
//	if str3 == string(str2[266:330]){
//				println("1111")
//	}
//
//}

//func main() {
//
//
//		// 倒序：
//		var kArray = []string{"1000", "1001", "1002", "1003", "1004", "1005"}
//		sort.Slice(kArray, func(i, j int) bool {
//			return kArray[i] > kArray[j]
//		})
//		fmt.Println("逆序：", kArray)
//		// 正序：
//		sort.Strings(kArray)
//		fmt.Println("正序：", kArray)
//
//
//}

//func main() {
//	m := redis.GetHashDataAll("buyAndSale:Metamon Egg")
//	fmt.Println(m)
//	data,_ := json.Marshal(m)
//	fmt.Println(string(data))
//
//}

//func main() {
//	data1 := "1.05E+05"
//	v1, err := strconv.ParseFloat(data1, 64)
//	if err != nil {
//		logger.Error(err)
//	}
//	fmt.Println(v1)
//	//fmt.Println(fmt.Sprintf("%T", v1))
//	//fmt.Println(fmt.Sprintf("%T", data))
//}
//func main() {
//	timeStr:=time.Now().Format("2006-01-02 15:04:05")
//	fmt.Println(timeStr)
//}