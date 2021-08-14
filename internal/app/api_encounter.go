package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gtosh4/WoWCDHelper/pkg/encounters"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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
	eiID, err := strconv.ParseUint(c.Param("eventinst"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	es := []encounters.Assignment{}
	err = s.db(c).
		Preload("EventInstance").
		Preload("Event").
		Preload("Encounter").
		Where(&encounters.Encounter{ID: uint(encId)}).
		Where(&encounters.EventInstance{ID: uint(eiID)}).
		Find(&es).
		Error

	if err != nil {
		s.errAbort(c, err)
		return
	}

	c.JSON(http.StatusOK, es)
}

func (s *Server) handleSetAssignments(c *gin.Context) {
	_, _, err := encounterParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}
	eiID, err := strconv.ParseUint(c.Param("eventinst"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var as []encounters.Assignment
	err = c.Bind(&as)
	if err != nil {
		return
	}
	for _, a := range as {
		a.EventInstanceID = uint(eiID)
	}

	err = s.db(c).Transaction(func(tx *gorm.DB) error {
		err := tx.
			Preload("EventInstance").
			Where(&encounters.EventInstance{ID: uint(eiID)}).
			Delete(&encounters.Assignment{}).
			Error
		if err != nil {
			return err
		}

		err = tx.
			Create(&as).
			Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		s.errAbort(c, err)
		return
	}
	c.JSON(http.StatusOK, as)
}
