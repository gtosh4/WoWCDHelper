package wcdh

import (
	"net/http"
	"net/http/pprof"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const localhost = "127.0.0.1:"

func localhostOnly(req *http.Request, m *mux.RouteMatch) bool {
	if strings.HasPrefix(req.RemoteAddr, localhost) {
		return true
	}
	return false
}

func (s *Server) registerMetricsHandler(router *mux.Router) {
	handler := promhttp.InstrumentMetricHandler(
		s.Prometheus, promhttp.HandlerFor(s.Prometheus, promhttp.HandlerOpts{
			ErrorLog: s.Log,
			Timeout:  10 * time.Second,
		}),
	)
	router.Path("/metrics").MatcherFunc(localhostOnly).Handler(handler)
}

func (s *Server) registerPprof() {
	pprofRouter := s.router.PathPrefix("/debug/pprof").MatcherFunc(localhostOnly).Subrouter()
	pprofRouter.PathPrefix("/cmdline").HandlerFunc(pprof.Cmdline)
	pprofRouter.PathPrefix("/profile").HandlerFunc(pprof.Profile)
	pprofRouter.PathPrefix("/symbol").HandlerFunc(pprof.Symbol)
	pprofRouter.PathPrefix("/trace").HandlerFunc(pprof.Trace)
	pprofRouter.PathPrefix("/").HandlerFunc(pprof.Index)
	pprofRouter.PathPrefix("").HandlerFunc(pprof.Index)
}
