package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gtosh4/WoWCDHelper/internal/app/wcdh"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	port := flag.Int("port", 8113, "port")

	flag.Parse()

	svc := wcdh.NewServer(log, *port)
	go func() {
		err := svc.Run()
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
	err := svc.Shutdown(10 * time.Second)
	if err != nil {
		log.WithError(err).Warn("Shutdown had errors")
	} else {
		log.Info("Exit")
	}
}
