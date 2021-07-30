package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/FuzzyStatic/blizzard/v2"
	"github.com/allegro/bigcache"
	"github.com/gregjones/httpcache"
	"github.com/gtosh4/WoWCDHelper/internal/app"
	"github.com/gtosh4/WoWCDHelper/internal/pkg/context"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var root = &cobra.Command{
	Use:  "wcdh",
	RunE: serve,
}

var cfg struct {
	Debug             bool
	Port              int
	BlizzClientID     string
	BlizzClientSecret string
}

func main() {
	flags := root.Flags()

	flags.BoolVar(&cfg.Debug, "debug", false, "enable debug logging")
	flags.IntVarP(&cfg.Port, "port", "p", 8080, "port to bind to")
	flags.StringVar(&cfg.BlizzClientID, "bnetId", "31708c8133144f6fab3b75e2ece62d3d", "Battle.net API client ID")
	flags.StringVar(&cfg.BlizzClientSecret, "bnetSecret", "", "Battle.net API client secret")

	root.Execute()
}

func serve(cmd *cobra.Command, args []string) (err error) {
	clients := &context.Clients{}

	logLevel := zap.NewAtomicLevel()
	if cfg.Debug {
		logLevel.SetLevel(zap.DebugLevel)
	}
	clients.Log = zap.New(
		zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), zapcore.AddSync(os.Stdout), logLevel),
	)

	clients.Blizz = blizzard.NewClient(
		cfg.BlizzClientID,
		cfg.BlizzClientSecret,
		blizzard.US,
		blizzard.EnUS,
	)

	clients.Cache, err = bigcache.NewBigCache(bigcache.DefaultConfig(time.Hour))
	if err != nil {
		err = errors.Wrap(err, "could not create cache")
		return
	}

	clients.Blizz.SetHTTPClient(wrapBlizzHTTP(clients, clients.Blizz.GetHTTPClient()))

	srv := app.NewServer(clients)

	return srv.Run(fmt.Sprintf(":%d", cfg.Port))
}

func wrapBlizzHTTP(clients *context.Clients, h *http.Client) *http.Client {
	transport := httpcache.NewTransport(&app.HTTPBigCache{
		Log:   clients.Log.Sugar().Named("cache"),
		Cache: clients.Cache,
	})
	transport.Transport = h.Transport
	h.Transport = transport
	return h
}
