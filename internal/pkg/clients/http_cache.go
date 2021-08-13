package clients

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type HTTPCache struct {
	Log    *zap.SugaredLogger
	Cache  *badger.DB
	Prefix string
	RNG    *rand.Rand

	met *httpCacheMetrics
}

type httpCacheMetrics struct {
	count   prometheus.Gauge
	size    prometheus.Gauge
	hitRate prometheus.Summary
}

func (c *HTTPCache) key(k string) []byte {
	return []byte(fmt.Sprintf("%s/%s", c.Prefix, k))
}

func (c *HTTPCache) ttl() time.Duration {
	var rng func() float64
	if c.RNG != nil {
		rng = c.RNG.Float64
	} else {
		rng = rand.Float64
	}
	jitter := time.Duration(8 * float64(time.Minute) * rng())
	return time.Hour + jitter - (4 * time.Minute)
}

func (c *HTTPCache) Get(key string) (b []byte, ok bool) {
	defer func() {
		c.ensureMetrics()
		hit := 0.0
		if ok {
			hit = 1.0
		}
		c.met.hitRate.Observe(hit)
	}()

	err := c.Cache.View(func(txn *badger.Txn) error {
		item, err := txn.Get(c.key(key))
		if err != nil {
			return err
		}
		b, err = item.ValueCopy(b)
		return err
	})
	if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
		c.Log.Warnf("got error retrieving %s: %v", key, err)
	}
	return b, err == nil
}

func (c *HTTPCache) Set(key string, responseBytes []byte) {
	err := c.Cache.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(c.key(key), responseBytes).
			WithTTL(c.ttl())
		return txn.SetEntry(entry)
	})
	if err != nil {
		c.Log.Warnf("got error setting %s: %v", key, err)
	}
}

func (c *HTTPCache) Delete(key string) {
	err := c.Cache.Update(func(txn *badger.Txn) error {
		return txn.Delete(c.key(key))
	})
	if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
		c.Log.Warnf("got error deleting %s: %v", key, err)
	}
}

func (c *HTTPCache) ensureMetrics() {
	if c.met != nil {
		return
	}

	c.met = &httpCacheMetrics{
		count: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "cache",
			Subsystem: "itemprefix",
			Name:      c.Prefix + "_count",
			Help:      "the number of items in this prefix",
		}),

		size: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "cache",
			Subsystem: "itemprefix",
			Name:      c.Prefix + "_bytes",
			Help:      "the number of bytes in this prefix",
		}),

		hitRate: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: "cache",
			Subsystem: "itemprefix",
			Name:      c.Prefix + "_hitrate",
		}),
	}
}

func (c *HTTPCache) Describe(ch chan<- *prometheus.Desc) {
	c.ensureMetrics()
	c.met.count.Describe(ch)
	c.met.size.Describe(ch)
	c.met.hitRate.Describe(ch)
}

func (c *HTTPCache) Collect(ch chan<- prometheus.Metric) {
	c.ensureMetrics()
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
