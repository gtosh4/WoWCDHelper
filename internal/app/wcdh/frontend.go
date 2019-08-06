package wcdh

import (
	"net/http"
)

func (s *Server) registerFrontend() {
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/dist")))
}
