package redis_store

type Limiter interface {
	Take() (bool, error)
}
