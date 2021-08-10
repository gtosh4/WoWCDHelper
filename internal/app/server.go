package app

import (
	"net/http"

	"github.com/chenjiandongx/ginprom"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/gtosh4/WoWCDHelper/internal/pkg/clients"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	router  *gin.Engine
	Log     *zap.Logger
	clients *clients.Clients
	srv     *http.Server
}

func NewServer(c *clients.Clients, addr string) *Server {
	gin.SetMode(gin.ReleaseMode)

	s := &Server{
		router:  gin.New(),
		Log:     c.Log.Named("server"),
		clients: c,
	}

	s.router.Use(ginzap.RecoveryWithZap(s.Log, true))
	s.router.Use(ginprom.PromMiddleware(nil))

	registerMetricsHandler(s)
	registerDebug(s)
	registerWoWApi(s)
	registerTeamApi(s)
	registerEncounterApi(s)

	registerFrontend(s)

	s.Log.Info("Routes:")
	for _, route := range s.router.Routes() {
		s.Log.With(zap.String("method", route.Method)).Info(route.Path)
	}

	s.srv = &http.Server{Addr: addr, Handler: s.router}

	return s
}

func (s *Server) Run() error {
	s.Log.Sugar().Infof("Listening at %s", s.srv.Addr)
	return s.srv.ListenAndServe()
}

func (s *Server) Close() error {
	return s.srv.Close()
}

func (s *Server) db(c *gin.Context) *gorm.DB {
	return s.clients.DB.WithContext(c)
}

func (s *Server) log(c *gin.Context) *zap.SugaredLogger {
	log := s.Log.Sugar()
	for _, param := range c.Params {
		log = log.With(zap.String(param.Key, param.Value))
	}
	log = log.With("handler", c.HandlerName())
	return log
}

func (s *Server) errAbort(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}
