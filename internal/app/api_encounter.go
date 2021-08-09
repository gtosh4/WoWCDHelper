package app

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gtosh4/WoWCDHelper/internal/pkg/clients"
	"github.com/gtosh4/WoWCDHelper/pkg/encounters"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func registerEncounterApi(s *Server) {
	teamR := s.router.Group("/team/:team")
	teamR.Use(clients.Ginzap(s.Log, time.RFC3339, true, zap.InfoLevel))
	teamR.GET("/encounters", s.handleGetEncounters)
	teamR.POST("/encounter", s.handleNewEncounter)

	encR := teamR.Group("/encounter/:encounter")
	encR.DELETE("", s.handleDeleteEncounter)
	encR.PUT("", s.handleSetEncounter)
	encR.GET("/events", s.handleGetEvents)
	encR.PUT("/events", s.handleSetEvents)
	encR.POST("/event", s.handleNewEvent)
	encR.GET("/roster", s.handleGetRoster)
	encR.PUT("/roster", s.handleSetRoster)
	encR.PUT("/roster/:member", s.handleSetRosterMember)
	encR.DELETE("/roster/:member", s.handleDeleteRosterMember)

	encR.GET("/assignments", s.handleGetAssignments)
}

func encounterParams(c *gin.Context) (teamId string, encId uint, err error) {
	teamId = c.Param("team")
	if teamId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		err = errors.New("empty team param")
		return
	}
	encId64, err := strconv.ParseUint(c.Param("encounter"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	encId = uint(encId64)
	return
}

func (s *Server) handleGetEncounters(c *gin.Context) {
	teamId := c.Param("team")
	if teamId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	es := []encounters.Encounter{}
	err := s.db(c).
		Where(&encounters.Encounter{TeamID: teamId}).
		Find(&es).
		Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.AbortWithStatus(http.StatusNotFound)

	case err != nil:
		c.AbortWithError(http.StatusInternalServerError, err)

	default:
		c.JSON(http.StatusOK, es)
	}
}

func (s *Server) handleDeleteEncounter(c *gin.Context) {
	teamId, encId, err := encounterParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}

	err = s.db(c).
		Delete(&encounters.Encounter{TeamID: teamId, ID: uint(encId)}).
		Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) handleSetEncounter(c *gin.Context) {
	teamId, encId, err := encounterParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}

	var enc encounters.Encounter
	err = c.Bind(&enc)
	if err != nil {
		s.log(c).Warnf("bind err: %v", err)
		return
	}
	enc.TeamID = teamId
	enc.ID = uint(encId)

	s.log(c).Infof("updating encounter %+v", enc)

	err = s.db(c).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&enc).
		Error

	if err != nil {
		s.errAbort(c, err)
		return
	}
	s.log(c).Infof("updated encounter %+v", enc)
	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) handleNewEncounter(c *gin.Context) {
	teamId := c.Param("team")
	if teamId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var e encounters.Encounter
	err := c.Bind(&e)
	if err != nil {
		s.log(c).Warnf("bind err: %v", err)
		return
	}
	e.TeamID = teamId

	err = s.db(c).
		Create(&e).
		Error
	if err != nil {
		s.errAbort(c, err)
		return
	}

	c.Redirect(http.StatusCreated, fmt.Sprintf("/team/%s/encounter/%d", e.TeamID, e.ID))
}

func (s *Server) handleGetAssignments(c *gin.Context) {
	_, encId, err := encounterParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}

	es := []encounters.Assignment{}
	err = s.db(c).
		Where(&encounters.Encounter{ID: uint(encId)}).
		Find(&es).
		Error

	if err != nil {
		s.errAbort(c, err)
		return
	}

	c.JSON(http.StatusOK, es)
}
