package app

import (
	"time"

	"github.com/chenjiandongx/ginprom"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/gtosh4/WoWCDHelper/internal/pkg/clients"
	"go.uber.org/zap"
)

type Server struct {
	router  *gin.Engine
	log     *zap.Logger
	clients *clients.Clients
}

func NewServer(clients *clients.Clients) *Server {
	gin.SetMode(gin.ReleaseMode)

	s := &Server{
		router:  gin.New(),
		log:     clients.Log.Named("server"),
		clients: clients,
	}

	s.router.Use(ginzap.Ginzap(s.log, time.RFC3339, true))
	s.router.Use(ginzap.RecoveryWithZap(s.log, true))
	s.router.Use(ginprom.PromMiddleware(nil))

	registerMetricsHandler(s)
	registerDebug(s)
	registerWoWApi(s)
	registerTeamApi(s)

	registerFrontend(s)

	return s
}

func (s *Server) Run(addr string) error {
	s.log.Sugar().Infof("Listening at %s", addr)
	return s.router.Run(addr)
}
