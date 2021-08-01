package app

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gtosh4/WoWCDHelper/pkg/teams"
	"gorm.io/gorm"
)

func registerTeamApi(s *Server) {
	r := s.router.Group("team")

	r.GET("/:id", s.handleGetTeam)
	r.PUT("/:id", s.handleSetTeam)
}

func (s *Server) handleGetTeam(c *gin.Context) {
	id := c.Param("id")
	s.log.Sugar().Debugf("loading team %s", id)

	var roster teams.Roster
	err := s.clients.DB.Where("team_id = ?", id).Find(&roster).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.AbortWithStatus(http.StatusNotFound)

	case err != nil:
		c.AbortWithError(http.StatusInternalServerError, err)

	default:
		c.JSON(http.StatusOK, roster)
	}
}

func (s *Server) handleSetTeam(c *gin.Context) {
	var roster teams.Roster
	err := c.Bind(&roster)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = s.clients.DB.Create(&roster).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}
