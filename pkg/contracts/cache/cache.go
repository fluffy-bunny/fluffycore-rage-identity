package cache

type (
	IScopedMemoryCache interface {
		Set(key string, value any)
		Get(key string) (any, bool)
		Delete(key string)
		Clear()
	}
)
