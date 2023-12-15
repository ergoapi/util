package cache

func InitGoCache(options ...Option) {
	Instance = NewGoCache(options...)
}

func InitGoRedis(options ...Option) {
	Instance = NewGoRedis(options...)
}

func InitGoRedisCluster(options ...Option) {
	Instance = NewGoRedisCluster(options...)
}
