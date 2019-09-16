package wcdh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	badger "github.com/dgraph-io/badger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/gtosh4/WoWCDHelper/internal/pkg/ctx"
	"github.com/gtosh4/WoWCDHelper/pkg/analysis"
	"github.com/gtosh4/WoWCDHelper/pkg/warcraftlogs"
)

type (
	Server struct {
		*ctx.Ctx
		Cfg *Config

		router *mux.Router
		srv    *http.Server

		thcda *analysis.THCDA
	}
	Config struct {
		Port    int    `yaml:"port"`
		ApiKey  string `yaml:"api_key"`
		BaseURL string `yaml:"base_url"`
	}
)

func NewServer(ctx *ctx.Ctx, cfg *Config, db *badger.DB) *Server {
	s := &Server{
		Ctx: ctx,
		Cfg: cfg,

		router: mux.NewRouter(),
	}

	s.router.Use(func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(ctx.Log.WithField("http", "request").WriterLevel(logrus.InfoLevel), h)
	})

	s.thcda = &analysis.THCDA{
		Ctx: s.Ctx.NewSubContext(),
		WCL: warcraftlogs.NewClient(ctx, http.DefaultClient, cfg.BaseURL, cfg.ApiKey, db),
	}
	s.thcda.RegisterHTTPHandlers(s.router)

	s.registerMetricsHandler(s.router)
	s.registerPprof()
	s.registerFrontend() // Do this last since it registers the root path

	s.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		if p, err := route.GetPathRegexp(); err == nil {
			s.Log.Debugf("Registered route: %s", p)
		}
		return nil
	})

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: s.router,

		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	return s
}

func (s *Server) Run() error {
	s.Log.Infof("Server starting on %s", s.srv.Addr)
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.srv.Shutdown(ctx)
}

func errorPayload(err error) interface{} {
	return map[string]string{"error": err.Error(), "details": fmt.Sprintf("%+v", err)}
}

func (s *Server) respondWithError(w http.ResponseWriter, code int, err error) {
	s.respondWithJSON(w, code, errorPayload(err))
}

func (s *Server) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(payload); err != nil {
		s.Log.WithError(err).Warnf("Error writing response")
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Wrapf(err, "Error writing response")
		if err := enc.Encode(errorPayload(err)); err != nil {
			s.Log.WithError(err).Warnf("Failed to write failure response, giving up.")
		}
		return
	}
}
