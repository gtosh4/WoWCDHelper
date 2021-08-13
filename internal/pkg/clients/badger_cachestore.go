package clients

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	GinCache struct {
		*Clients

		Prefix string
		TTL    time.Duration

		wrapperPool sync.Pool
		bufPool     sync.Pool

		met ginCacheMetrics
	}

	cacheData struct {
		Status      int
		ContentType string
		Payload     []byte
		Expires     string
	}

	cacheWriter struct {
		gin.ResponseWriter
		cache *GinCache

		payload *bytes.Buffer
	}

	ginCacheMetrics struct {
		count   prometheus.Gauge
		size    prometheus.Gauge
		hitRate prometheus.Summary
	}
)

func init() {
	gob.Register(&cacheData{})
}

func NewGinCache(clients *Clients, prefix string, ttl time.Duration) *GinCache {
	c := &GinCache{
		Clients: clients,
		Prefix:  prefix,
		TTL:     ttl,
		met: ginCacheMetrics{
			count: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: "gincache",
				Subsystem: prefix,
				Name:      "count",
				Help:      "the number of items in this prefix",
			}),

			size: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: "gincache",
				Subsystem: prefix,
				Name:      "bytes",
				Help:      "the number of bytes in this prefix",
			}),

			hitRate: prometheus.NewSummary(prometheus.SummaryOpts{
				Namespace: "gincache",
				Subsystem: prefix,
				Name:      "hitrate",
			}),
		},
	}
	c.wrapperPool.New = func() interface{} {
		return &cacheWriter{
			cache:   c,
			payload: bytes.NewBuffer(nil),
		}
	}
	c.bufPool.New = func() interface{} {
		return bytes.NewBuffer(nil)
	}
	return c
}

func (c *GinCache) key(k string) []byte {
	return []byte(fmt.Sprintf("gin/%s/%s", c.Prefix, k))
}

func (c *GinCache) Handle(ctx *gin.Context) {
	key := ctx.Request.URL.RequestURI()
	var value cacheData

	log := c.Log.Sugar().Named("gincache")

	err := c.Get(key, &value)
	log.Debugf("read cache value %s (err %v)", value, err)

	if err != nil {
		c.met.hitRate.Observe(0)
		cw := c.wrapperPool.Get().(*cacheWriter)
		defer c.wrapperPool.Put(cw)

		value.Expires = time.Now().Add(c.TTL).UTC().Format(http.TimeFormat)
		cw.ResponseWriter = ctx.Writer
		cw.payload.Reset()

		ctx.Writer = cw

		ctx.Header("Cache-Control", "public, immutable")
		ctx.Header("Expires", value.Expires)
		ctx.Next()

		value.ContentType = ctx.Writer.Header().Get("Content-Type")
		value.Payload = cw.payload.Bytes()
		value.Status = ctx.Writer.Status()

		c.Set(key, &value)
		log.Infof("set cache %s: %s", key, value)
		return
	}

	c.met.hitRate.Observe(1)
	ctx.Header("Cache-Control", "public, immutable")
	ctx.Header("Expires", value.Expires)
	ctx.Data(value.Status, value.ContentType, value.Payload)
	ctx.Abort()
}

func (c *GinCache) Get(key string, value *cacheData) error {
	return c.Cache.View(func(txn *badger.Txn) error {
		item, err := txn.Get(c.key(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return gob.NewDecoder(bytes.NewReader(val)).Decode(value)
		})
	})
}

func (c *GinCache) Set(key string, value *cacheData) error {
	b := c.bufPool.Get().(*bytes.Buffer)
	defer c.bufPool.Put(b)
	b.Reset()

	err := gob.NewEncoder(b).Encode(value)
	if err != nil {
		return err
	}
	return c.Cache.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(c.key(key), b.Bytes()).
			WithTTL(c.TTL)
		return txn.SetEntry(entry)
	})
}

func (c *GinCache) Flush() error {
	return c.Cache.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = c.key("")
		opts.PrefetchValues = false
		iter := txn.NewIterator(opts)
		defer iter.Close()

		for iter.Seek(opts.Prefix); iter.Valid(); iter.Next() {
			if err := txn.Delete(iter.Item().Key()); err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *GinCache) Describe(ch chan<- *prometheus.Desc) {
	c.met.count.Describe(ch)
	c.met.size.Describe(ch)
	c.met.hitRate.Describe(ch)
}

func (c *GinCache) Collect(ch chan<- prometheus.Metric) {
	var (
		count = 0
		size  = 0
	)

	c.Cache.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = c.key("")
		opts.PrefetchValues = false
		iter := txn.NewIterator(opts)
		defer iter.Close()

		for iter.Seek(opts.Prefix); iter.Valid(); iter.Next() {
			size += int(iter.Item().ValueSize())
			count++
		}

		return nil
	})

	c.met.count.Set(float64(count))
	c.met.count.Collect(ch)

	c.met.size.Set(float64(size))
	c.met.size.Collect(ch)

	c.met.hitRate.Collect(ch)
}

func (w *cacheWriter) WriteString(s string) (int, error) {
	w.payload.Write([]byte(s))
	return w.ResponseWriter.WriteString(s)
}

func (w *cacheWriter) Write(b []byte) (int, error) {
	w.payload.Write(b)
	return w.ResponseWriter.Write(b)
}

func (v cacheData) String() string {
	return fmt.Sprintf("{%d %s %d bytes}", v.Status, v.ContentType, len(v.Payload))
}
