/**
 @author:way
 @date:2021/11/27
 @note
**/

package main

import (
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logic"
	"redisData/pkg/email"
	"redisData/pkg/logger"
	"redisData/setting"
	"strconv"
)

func init() {
	// 定义日志目录
	logger.Init("marketMonitor")
	// 初始化 viper 配置
	if err := setting.Init(""); err != nil {
		logger.Error(err)
		return
	}
	mysql.InitMysql()
	// 初始化redis
	if err := redis.InitClient(); err != nil {
		logger.Error(err)
		return
	}
}
func startMonitor(m map[string]interface{}) {
	name := m["name"].(string)
	types := m["type"].(string)
	//判断是元兽蛋还是药水
	switch name {
	case "Metamon Egg":
		EggStart(types)
	case "Potion":
		PotionStart(types)
	}
}

func EggStart(types string) {
	if types == "fall" {
		//获取鸡蛋跌的数据
		m := redis.GetHashDataAll("risk:fall")
		//获取鸡蛋升的数据
		timeLevel, _ := strconv.Atoi(m["TimeLevel"])
		percentage, _ := strconv.ParseFloat(m["Percentage"], 64)
		operationType := m["TimeLevel"]
		status := m["Status"]
		if status != "1" {
			return
		}
		//旧的数据
		data := mysql.GetHistoryMarketData(timeLevel, "Metamon Egg.MarketPrice")
		//新的数据
		newMarketPrice := logic.GetMarketDataByRedis("Metamon Egg.List")
		logger.Info(data.MarketData)
		logger.Info(newMarketPrice)
		logger.Info(percentage)

		tmp := 1 - (data.MarketData / newMarketPrice)
		if tmp > 0 {
			return
		}
		logger.Info(tmp >= (percentage * 0.01))
		if tmp >= (percentage * 0.01) {
			//停止买入脚本，且发邮件通知,使用上一次的market和现在的market对比，上一次的market从redis中读，新的marketPrice重新算
			switch operationType {
			case "1":
				logger.Info("停止脚本")
				logic.StopScript()
			case "2":
				logger.Info("发送钉钉")
				email.SendDingMsg("egg:risk:fall","鸡蛋下跌风险")
			case "3":
				logger.Info("停止脚本且发送钉钉")
				logic.StopScript()
			}
		}
	}

	if types == "rise" {
		m := redis.GetHashDataAll("risk:rise")
		//获取鸡蛋升的数据
		timeLevel, _ := strconv.Atoi(m["TimeLevel"])
		percentage, _ := strconv.ParseFloat(m["Percentage"], 64)
		operationType := m["TimeLevel"]
		status := m["TimeLevel"]
		if status != "1" {
			return
		}

		//旧的数据
		data := mysql.GetHistoryMarketData(timeLevel, "Metamon Egg.MarketPrice")
		//新的数据
		newMarketPrice := logic.GetMarketDataByRedis("Metamon Egg.List")

		if (newMarketPrice/data.MarketData)-1 >= (percentage * 0.01) {
			//停止买入脚本，且发邮件通知,使用上一次的market和现在的market对比，上一次的market从redis中读，新的marketPrice重新算
			switch operationType {
			case "1":
				logger.Info("停止脚本")
				logic.StopScript()
			case "2":
				logger.Info("发送钉钉")
				email.SendDingMsg("egg:risk:rise","鸡蛋上涨风险")
			case "3":
				logger.Info("停止脚本且发送钉钉")
				logic.StopScript()
			}
		}
	}
}

func PotionStart(types string) {
	if types == "fall" {
		//获取potion跌的数据
		m := redis.GetHashDataAll("risk:potion:fall")
		//获取鸡蛋升的数据
		timeLevel, _ := strconv.Atoi(m["TimeLevel"])
		percentage, _ := strconv.ParseFloat(m["Percentage"], 64)
		operationType := m["TimeLevel"]
		status := m["TimeLevel"]
		if status != "1" {
			return
		}
		//旧的数据
		data := mysql.GetHistoryMarketData(timeLevel, "Potion.MarketPrice")
		//新的数据
		newMarketPrice := logic.GetMarketDataByRedis("Potion.List")
		if 1-(data.MarketData/newMarketPrice) >= (percentage * 0.01) {
			//停止买入脚本，且发邮件通知,使用上一次的market和现在的market对比，上一次的market从redis中读，新的marketPrice重新算
			switch operationType {
			case "1":
				logger.Info("停止脚本")
				logic.StopScript()
			case "2":
				logger.Info("发送钉钉")
				email.SendDingMsg("risk:potion:fall","药水下跌风险")
			case "3":
				logger.Info("停止脚本且发送钉钉")
				logic.StopScript()
			}
		}
	}

	if types == "rise" {
		//获取potion升的数据
		m := redis.GetHashDataAll("risk:potion:rise")

		timeLevel, _ := strconv.Atoi(m["TimeLevel"])
		percentage, _ := strconv.ParseFloat(m["Percentage"], 64)
		operationType := m["TimeLevel"]
		status := m["TimeLevel"]
		if status != "1" {
			return
		}
		//旧的数据
		data := mysql.GetHistoryMarketData(timeLevel, "Potion.MarketPrice")
		//新的数据
		newMarketPrice := logic.GetMarketDataByRedis("Potion.List")
		if 1-(data.MarketData/newMarketPrice) >= (percentage * 0.01) {
			//停止买入脚本，且发邮件通知,使用上一次的market和现在的market对比，上一次的market从redis中读，新的marketPrice重新算
			switch operationType {
			case "1":
				logger.Info("停止脚本")
				logic.StopScript()
			case "2":
				logger.Info("发送钉钉")
				email.SendDingMsg("risk:potion:rise","药水上涨风险")
			case "3":
				logger.Info("停止脚本且发送钉钉")
				logic.StopScript()
			}
		}
	}
}

func main() {
	defer redis.Close()
	m1 := map[string]interface{}{
		"name": "Metamon Egg",
		"type": "fall",
	}
	m2 := map[string]interface{}{
		"name": "Metamon Egg",
		"type": "rise",
	}
	m3 := map[string]interface{}{
		"name": "Potion",
		"type": "fall",
	}
	m4 := map[string]interface{}{
		"name": "Potion",
		"type": "rise",
	}

	for {
		startMonitor(m1)
		startMonitor(m2)
		startMonitor(m3)
		startMonitor(m4)
	}
}
