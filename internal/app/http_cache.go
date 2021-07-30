package app

import (
	"errors"

	"github.com/allegro/bigcache"
	"go.uber.org/zap"
)

type HTTPBigCache struct {
	Log   *zap.SugaredLogger
	Cache *bigcache.BigCache
}

func (c *HTTPBigCache) Get(key string) (b []byte, ok bool) {
	b, err := c.Cache.Get(key)
	if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
		c.Log.Warnf("got error retrieving %s: %v", key, err)
	}
	return b, err == nil
}

func (c *HTTPBigCache) Set(key string, responseBytes []byte) {
	err := c.Cache.Set(key, responseBytes)
	if err != nil {
		c.Log.Warnf("got error setting %s: %v", key, err)
	}
}
func (c *HTTPBigCache) Delete(key string) {
	err := c.Cache.Delete(key)
	if err != nil {
		c.Log.Warnf("got error deleting %s: %v", key, err)
	}
}
