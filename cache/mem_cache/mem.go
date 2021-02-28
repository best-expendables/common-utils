package mem_cache

import (
	"reflect"
	"time"

	commonCache "github.com/best-expendables/common-utils/cache"
	"github.com/patrickmn/go-cache"
)

type Mem struct {
	c *cache.Cache
}

func NewMem(ttl time.Duration) *Mem {
	return &Mem{cache.New(ttl, 10*time.Minute)}
}

func (m *Mem) Get(key string, obj interface{}) error {
	value, found := m.c.Get(key)
	if found {
		v := reflect.ValueOf(obj).Elem()
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		v.Set(rv)
		return nil
	}
	return commonCache.Nil
}

func (m *Mem) Set(key string, obj interface{}) error {
	m.c.Set(key, obj, cache.DefaultExpiration)
	return nil
}

func (m *Mem) Delete(key string) error {
	m.c.Delete(key)
	return nil
}
