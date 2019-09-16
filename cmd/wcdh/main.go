package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	metrics "github.com/armon/go-metrics"
	prom_met "github.com/armon/go-metrics/prometheus"
	"github.com/dgraph-io/badger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/gtosh4/WoWCDHelper/internal/app/wcdh"
	"github.com/gtosh4/WoWCDHelper/internal/pkg/ctx"
)

func main() {
	log := logrus.New()

	mustNotErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	port := flag.Int("port", 8113, "port")
	debug := flag.Bool("debug", false, "Debug logging")
	cachePath := flag.String("cachePath", "./cache", "Cache location")
	apiKey := flag.String("api-key", "", "API Key")
	baseURL := flag.String("base-url", "https://www.warcraftlogs.com/v1", "Base url")

	flag.Parse()

	prometheusSink, err := prom_met.NewPrometheusSink()
	mustNotErr(err)

	prometheusRegistry := prometheus.NewRegistry()
	err = prometheusRegistry.Register(prometheusSink)
	mustNotErr(err)

	metCfg := &metrics.Config{
		ServiceName:          "thcda",
		HostName:             mustHostname(),
		EnableServiceLabel:   true,
		EnableHostnameLabel:  true,
		EnableRuntimeMetrics: true,
		TimerGranularity:     time.Millisecond,
		ProfileInterval:      time.Second,
		FilterDefault:        true,
	}
	met, err := metrics.New(metCfg, prometheusSink)
	mustNotErr(err)

	if *debug {
		log.SetLevel(logrus.DebugLevel)
	}

	ctx := ctx.NewContext(log, prometheusRegistry, met)

	db, err := badger.Open(badger.DefaultOptions(*cachePath))
	mustNotErr(err)
	ctx.Tomb.Go(func() error {
		<-ctx.Tomb.Dying()
		db.Close()
		return nil
	})

	srv := wcdh.NewServer(
		ctx,
		&wcdh.Config{
			Port:    *port,
			ApiKey:  *apiKey,
			BaseURL: *baseURL,
		},
		db,
	)

	go func() {
		err := srv.Run()
		if err != nil {
			log.WithError(err).Warn("Server shutdown with error")
		} else {
			log.Info("Server shutdown")
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Infof("Shutdown requested via %s", sig)
	err = srv.Shutdown(10 * time.Second)
	if err != nil {
		log.WithError(err).Warn("Shutdown had errors")
	} else {
		log.Info("Exit")
	}
}

func mustHostname() string {
	h, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return h
}
