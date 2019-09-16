package ctx

import (
	"context"

	metrics "github.com/armon/go-metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	tomb "gopkg.in/tomb.v2"
)

type Ctx struct {
	context.Context

	Tomb       *tomb.Tomb
	Log        logrus.FieldLogger
	Prometheus *prometheus.Registry
	Metrics    *metrics.Metrics
}

func NewContext(
	Log logrus.FieldLogger,
	Prometheus *prometheus.Registry,
	Metrics *metrics.Metrics) *Ctx {

	t, ctx := tomb.WithContext(context.Background())
	return &Ctx{
		Context: ctx,
		Tomb:    t,

		Log:        Log,
		Prometheus: Prometheus,
		Metrics:    Metrics,
	}
}

func (ctx *Ctx) NewSubContext() (sub *Ctx) {
	sub = &Ctx{}
	*sub = *ctx
	sub.Tomb, sub.Context = tomb.WithContext(ctx.Context)
	return
}

func (ctx *Ctx) NewSubcontextWithFields(f logrus.Fields) *Ctx {
	sub := ctx.NewSubContext()
	sub.Log = ctx.Log.WithFields(f)
	return sub
}

func (ctx *Ctx) NewSubcontextWithField(key string, value interface{}) *Ctx {
	sub := ctx.NewSubContext()
	sub.Log = ctx.Log.WithField(key, value)
	return sub
}
