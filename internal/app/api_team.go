package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gtosh4/WoWCDHelper/pkg/teams"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func registerTeamApi(s *Server) {
	r := s.router.Group("team")

	r.POST("/new", s.handleCreateTeam)
	r.GET("/:id", s.handleGetTeam)
	r.PUT("/:id", s.handleSetTeam)

	r.POST("/:id/member", s.handleNewMember)
	r.PUT("/:id/:member", s.handleUpdateMember)
	r.DELETE("/:id/:member", s.handleDeleteMember)
}

func (s *Server) handleGetTeam(c *gin.Context) {
	id := c.Param("id")
	if id == "" || id == "new" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var roster teams.Roster
	err := s.clients.DB.Where("team_id = ?", id).Preload("Config").Find(&roster).Error
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
		return
	}
	for i := range roster {
		roster[i].TeamID = id
	}
	err = s.clients.DB.
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&roster).
		Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, roster)
}

func (s *Server) handleCreateTeam(c *gin.Context) {
	team := teams.Team{}
	var member []teams.Member
	err := c.ShouldBind(&member)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = s.clients.DB.Create(&team).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusCreated, fmt.Sprintf("/team/%s", team.ID))
}

func (s *Server) handleNewMember(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var m teams.Member
	err := c.Bind(&m)
	if err != nil {
		return
	}
	m.TeamID = id

	s.log.Sugar().Infof("Creating %+v", m)

	err = s.clients.DB.
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&m).
		Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	s.log.Sugar().Infof("Created %+v", m)

	c.JSON(http.StatusOK, m)
}

func (s *Server) handleUpdateMember(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	memberId, err := strconv.ParseInt(c.Param("member"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var m teams.Member
	err = c.Bind(&m)
	if err != nil {
		return
	}
	if m.ID != uint(memberId) || m.TeamID != id {
		c.AbortWithError(
			http.StatusBadRequest,
			errors.Errorf("{id: %d, team: %s} does not match path values {id: %d, team: %s}",
				m.ID, m.TeamID,
				memberId, id,
			),
		)
	}
	err = s.clients.DB.
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&m).
		Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) handleDeleteMember(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	memberId, err := strconv.ParseInt(c.Param("member"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	m := teams.Member{ID: uint(memberId), TeamID: id}
	err = s.clients.DB.Delete(&m).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
