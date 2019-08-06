package wcdh

import (
	"net/http/pprof"
)

func (s *Server) registerPprof() {
	pprofRouter := s.router.PathPrefix("/debug/pprof").Subrouter()
	pprofRouter.PathPrefix("/cmdline").HandlerFunc(pprof.Cmdline)
	pprofRouter.PathPrefix("/profile").HandlerFunc(pprof.Profile)
	pprofRouter.PathPrefix("/symbol").HandlerFunc(pprof.Symbol)
	pprofRouter.PathPrefix("/trace").HandlerFunc(pprof.Trace)
	pprofRouter.PathPrefix("/").HandlerFunc(pprof.Index)
	pprofRouter.PathPrefix("").HandlerFunc(pprof.Index)
}
