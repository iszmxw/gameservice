package redis

const (
	Prefix = "kline:" // 项目key前缀
	EggPrefix = "egg:"
	PotionPrefix = "potion:"
)

// 给redis key加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}

func getEggData(key string) string {
	return EggPrefix + key
}

func getPotionData(key string) string {
	return PotionPrefix + key
}



