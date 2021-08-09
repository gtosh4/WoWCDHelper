package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gtosh4/WoWCDHelper/pkg/encounters"
	"github.com/gtosh4/WoWCDHelper/pkg/teams"
	"gorm.io/gorm"
)

func (s *Server) handleNewEvent(c *gin.Context) {
	teamId := c.Param("team")
	encIdStr := c.Param("encounter")
	if teamId == "" || encIdStr == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	encId, err := strconv.ParseUint(encIdStr, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var e encounters.Event
	err = c.Bind(&e)
	if err != nil {
		return
	}
	e.ID = 0 // make sure we generate a new id
	e.EncounterID = uint(encId)

	err = s.db(c).
		Create(&e).
		Error
	if err != nil {
		s.errAbort(c, err)
		return
	}

	// c.Redirect(
	// 	http.StatusCreated,
	// 	fmt.Sprintf("/team/%s/encounter/%d/event/%d", teamId, e.EncounterID, e.ID),
	// )
	c.AbortWithStatus(http.StatusCreated)
}

func (s *Server) handleGetEvents(c *gin.Context) {
	teamId := c.Param("team")
	encIdStr := c.Param("encounter")
	if teamId == "" || encIdStr == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	encId, err := strconv.ParseUint(encIdStr, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	es := []encounters.Event{}
	err = s.db(c).
		Where(&teams.Team{ID: teamId}).
		Where(&encounters.Encounter{ID: uint(encId)}).
		Preload("Instances").
		Find(&es).
		Error

	if err != nil {
		s.errAbort(c, err)
		return
	}
	c.JSON(http.StatusOK, es)
}

func (s *Server) handleSetEvents(c *gin.Context) {
	teamId := c.Param("team")
	encIdStr := c.Param("encounter")
	if teamId == "" || encIdStr == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	encId, err := strconv.ParseUint(encIdStr, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	events := []encounters.Event{}
	err = c.Bind(&events)
	if err != nil {
		return
	}
	for i := range events {
		events[i].EncounterID = uint(encId)
	}
	err = s.db(c).Transaction(func(tx *gorm.DB) error {
		err := tx.
			Where(&encounters.Encounter{ID: uint(encId)}).
			Delete(&encounters.Event{}).
			Error
		if err != nil {
			return err
		}

		err = tx.
			Create(&events).
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
	c.JSON(http.StatusOK, events)
}
