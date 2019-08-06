package wcdh

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Port int

	router *mux.Router
	srv    *http.Server
}

func NewServer(log logrus.FieldLogger, port int) *Server {
	s := &Server{
		Port:   port,
		router: mux.NewRouter(),
	}

	s.router.Use(func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(log.WithField("http", "request").WriterLevel(logrus.InfoLevel), h)
	})
	s.registerPprof()
	s.registerFrontend() // Do this last since it registers the root path

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.router,
	}
	return s
}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.srv.Shutdown(ctx)
}
