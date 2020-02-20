package cache

import "errors"

var (
	Nil = errors.New("cache: key is missing")
)

type Cache interface {
	Get(key string, obj interface{}) error
	MGet(keys ...string) ([]interface{}, error)
	Set(key string, obj interface{}) error
	MSet(obj ...interface{}) error
	HSet(key string, field string, obj interface{}) error
	HGet(key string, field string) (string, error)
	HGetAll(key string) (map[string]string, error)
	Delete(key string) error
}
