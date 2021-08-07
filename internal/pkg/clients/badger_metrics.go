package clients

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/ristretto"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type BadgerMetricCollector struct {
	DB *badger.DB

	blockMetrics cacheMetrics
	indexMetrics cacheMetrics
}

type cacheMetrics map[string]prometheus.Gauge

func NewBadgerMetricCollector(db *badger.DB) *BadgerMetricCollector {
	col := &BadgerMetricCollector{
		DB:           db,
		blockMetrics: make(cacheMetrics),
		indexMetrics: make(cacheMetrics),
	}

	cacheMetric := func(name, desc string) {
		col.blockMetrics[name] = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "cache",
			Subsystem: "block",
			Name:      name,
			Help:      desc,
		})
		col.indexMetrics[name] = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "cache",
			Subsystem: "index",
			Name:      name,
			Help:      desc,
		})
	}

	cacheMetric("hits", "the number of Get calls where a value was found for the corresponding key")
	cacheMetric("misses", "the number of Get calls where a value was not found for the corresponding key")
	cacheMetric("keys_added", "the total number of Set calls where a new key-value item was added")
	cacheMetric("keys_updated", "the total number of Set calls where the value was updated")
	cacheMetric("keys_evicted", "the total number of keys evicted")
	cacheMetric("cost_added", "the sum of costs that have been added (successful Set calls)")
	cacheMetric("cost_evicted", "the sum of all costs that have been evicted")
	cacheMetric("sets_dropped", "the number of Set calls that don't make it into internal buffers (due to contention or some other reason).")
	cacheMetric("sets_rejected", "the number of Set calls rejected by the policy (TinyLFU)")
	cacheMetric("gets_dropped", "the number of Get counter increments that are dropped internally.")
	cacheMetric("gets_kept", "the number of Get counter increments that are kept")
	cacheMetric("ratio", "the number of Hits over all accesses (Hits + Misses). This is the percentage of successful Get calls.")

	return col
}

func (col *BadgerMetricCollector) Describe(c chan<- *prometheus.Desc) {
	for _, g := range col.blockMetrics {
		g.Describe(c)
	}
	for _, g := range col.indexMetrics {
		g.Describe(c)
	}
}

func (col *BadgerMetricCollector) collectPrefixed(c chan<- prometheus.Metric, met *ristretto.Metrics, prefix string) {
	zap.S().Infof("collecting %s: %+v", prefix, met)
	if met == nil {
		return
	}

	set := func(n string, v float64) {
		g := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "cache",
			Subsystem: prefix,
			Name:      n,
		})
		g.Set(v)
	}

	setI := func(n string, v uint64) {
		set(n, float64(v))
	}

	setI("hits", met.Hits())
	setI("misses", met.Misses())
	setI("keys_added", met.KeysAdded())
	setI("keys_updated", met.KeysUpdated())
	setI("keys_evicted", met.KeysEvicted())
	setI("cost_added", met.CostAdded())
	setI("cost_evicted", met.CostEvicted())
	setI("sets_dropped", met.SetsDropped())
	setI("sets_rejected", met.SetsRejected())
	setI("gets_dropped", met.GetsDropped())
	setI("gets_kept", met.GetsKept())
	set("ratio", met.Ratio())
}

func (col *BadgerMetricCollector) Collect(c chan<- prometheus.Metric) {
	col.collectPrefixed(c, col.DB.BlockCacheMetrics(), "block")
	col.collectPrefixed(c, col.DB.IndexCacheMetrics(), "index")
}
