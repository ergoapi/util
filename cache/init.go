package cache

func InitGoCache(options ...Option) {
	Instance = NewGoCache(options...)
}
