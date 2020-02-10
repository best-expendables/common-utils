package cache

import "errors"

var (
	Nil = errors.New("cache: key is missing")
)

type Cache interface {
	Get(key string, obj interface{}) error
	MGet(keys ...string) ([]interface{}, error)
	Set(key string, obj interface{}) error
	MSet(obj []interface{}) error
	Delete(key string) error
}
