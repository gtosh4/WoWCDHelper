package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerFrontend(s *Server) {
	directory := http.Dir("./frontend/public")
	fileserver := http.FileServer(directory)
	s.router.Use(SPAMiddleware(s, directory, fileserver))
	s.router.NoRoute(SPAFallbackHandler(s, fileserver))
}

func SPAMiddleware(s *Server, directory http.FileSystem, fileserver http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := directory.Open(c.Request.URL.Path); err == nil {
			s.log(c).Infof("serving %s", c.Request.URL)
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func SPAFallbackHandler(s *Server, fileserver http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(old string) {
			c.Request.URL.Path = old
		}(c.Request.URL.Path)

		s.log(c).Info("responding with fallback index.html")

		c.Request.URL.Path = "/"
		fileserver.ServeHTTP(c.Writer, c.Request)
	}
}
