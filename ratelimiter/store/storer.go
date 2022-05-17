package store

import "time"

type LimitStorer interface {
	Inc(key string, d time.Time) error
	Get(key string, previousWindow, currentWindow time.Time) (prevValue int64, currValue int64, err error)
	UpdateToken()
}
