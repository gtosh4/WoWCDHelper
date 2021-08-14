package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gtosh4/WoWCDHelper/internal/pkg/clients"
	"github.com/gtosh4/WoWCDHelper/pkg/teams"
	"go.uber.org/zap"
)

func registerTeamApi(s *Server) {
	logMW := clients.Ginzap(s.Log, time.RFC3339, true, zap.InfoLevel)
	s.router.POST("/teams", logMW, s.handleCreateTeam)

	teamR := s.router.Group("/team/:team", logMW)
	teamR.GET("", s.handleGetTeam)
	teamR.PUT("", s.handleSetTeam)

	teamR.GET("/members", s.handleGetTeamMembers)
	teamR.PUT("/members", s.handleSetTeamMembers)
	teamR.POST("/members", s.handleNewMember)

	memberR := teamR.Group("/member/:member")
	memberR.GET("", s.handleGetMember)
	memberR.PUT("", s.handleUpdateMember)
	memberR.DELETE("", s.handleDeleteMember)
	memberR.GET("/encounters", s.handleGetMemberEncounters)
	memberR.PUT("/encounters", s.handleSetMemberEncounters)
	memberR.DELETE("/encounters", s.handlRemoveMemberEncounters)

	teamR.GET("/encounters", s.handleGetEncounters)
	teamR.POST("/encounter", s.handleNewEncounter)

	encR := teamR.Group("/encounter/:encounter")
	encR.DELETE("", s.handleDeleteEncounter)
	encR.PUT("", s.handleSetEncounter)

	encR.GET("/events", s.handleGetEvents)
	encR.POST("/events", s.handleNewEvent)
	encR.GET("/event/:event", s.handleGetEvent)
	encR.PUT("/event/:event", s.handleSetEvent)

	rosterR := encR.Group("/roster")
	rosterR.GET("", s.handleGetRoster)
	rosterR.PUT("", s.handleSetRoster)
	rosterR.GET("/:member", s.handleGetRosterMember)
	rosterR.PUT("/:member", s.handleSetRosterMember)
	rosterR.DELETE("/:member", s.handleDeleteRosterMember)

	encR.GET("/assignments/:eventinst", s.handleGetAssignments)
	encR.PUT("/assignments/:eventinst", s.handleSetAssignments)
}

func (s *Server) handleCreateTeam(c *gin.Context) {
	team := teams.Team{}
	err := s.db(c).
		Create(&team).
		Error
	if err != nil {
		s.errAbort(c, err)
		return
	}
	c.Redirect(http.StatusCreated, fmt.Sprintf("/team/%s", team.ID))
}

func (s *Server) handleGetTeam(c *gin.Context) {
	teamId := c.Param("team")
	if teamId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var team teams.Team
	err := s.db(c).
		Where(&teams.Team{ID: teamId}).
		First(&team).
		Error

	if err != nil {
		s.errAbort(c, err)
		return
	}
	c.JSON(http.StatusOK, team)
}

func (s *Server) handleSetTeam(c *gin.Context) {
	teamId := c.Param("team")
	if teamId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var team teams.Team
	err := c.Bind(&team)
	if err != nil {
		return
	}
	team.ID = teamId

	err = s.db(c).
		Save(&team).
		Error
	if err != nil {
		s.errAbort(c, err)
		return
	}
	c.JSON(http.StatusOK, team)
}
