package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerFrontend(s *Server) {
	s.router.NoRoute(gin.WrapH(http.FileServer(http.Dir("./frontend/public"))))
}
