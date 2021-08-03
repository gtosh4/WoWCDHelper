package app

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gtosh4/WoWCDHelper/pkg/teams"
	"gorm.io/gorm"
)

func registerTeamApi(s *Server) {
	r := s.router.Group("team")

	r.POST("/new", s.handleCreateTeam)
	r.GET("/:id", s.handleGetTeam)
	r.PUT("/:id", s.handleSetTeam)
}

func (s *Server) handleGetTeam(c *gin.Context) {
	id := c.Param("id")
	if id == "" || id == "new" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
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
	id := c.Param("id")
	if id == "" || id == "new" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var roster teams.Roster
	err := c.Bind(&roster)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for i := range roster {
		roster[i].TeamID = id
	}
	err = s.clients.DB.Create(&roster).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, roster)
}

func (s *Server) handleCreateTeam(c *gin.Context) {
	team := teams.Team{}
	err := s.clients.DB.Create(&team).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusCreated, fmt.Sprintf("/team/%s", team.ID))
}
