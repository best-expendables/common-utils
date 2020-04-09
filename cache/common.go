package cache

import (
	"context"
	"errors"
)

var (
	Nil = errors.New("cache: key is missing")
)

type Cache interface {
	Get(ctx context.Context, key string, obj interface{}) error
	MGet(ctx context.Context, keys ...string) ([]interface{}, error)
	Set(ctx context.Context, key string, obj interface{}) error
	MSet(ctx context.Context, obj ...interface{}) error
	HSet(ctx context.Context, key string, field string, obj interface{}) error
	HGet(ctx context.Context, key string, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	Delete(ctx context.Context, key string) error
	ScanD(ctx context.Context, match string) error
	HExpire(ctx context.Context, key string) error
}
