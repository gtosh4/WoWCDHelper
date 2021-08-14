package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gtosh4/WoWCDHelper/pkg/encounters"
)

func (s *Server) handleNewEvent(c *gin.Context) {
	teamId, encId, err := encounterParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}

	var e encounters.Event
	err = c.Bind(&e)
	if err != nil {
		return
	}
	e.ID = 0 // make sure we generate a new id
	e.EncounterID = encId

	for i := range e.Instances {
		e.Instances[i].EventID = 0
		e.Instances[i].ID = 0
	}

	err = s.db(c).
		Create(&e).
		Error
	if err != nil {
		s.errAbort(c, err)
		return
	}

	c.Redirect(
		http.StatusCreated,
		fmt.Sprintf("/team/%s/encounter/%d/event/%d", teamId, e.EncounterID, e.ID),
	)
}

func (s *Server) handleGetEvents(c *gin.Context) {
	_, encId, err := encounterParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}

	es := []encounters.Event{}
	err = s.db(c).
		Where(&encounters.Event{EncounterID: encId}).
		Preload("Instances").
		Find(&es).
		Error

	if err != nil {
		s.errAbort(c, err)
		return
	}
	c.JSON(http.StatusOK, es)
}

func (s *Server) handleGetEvent(c *gin.Context) {
	_, encId, err := encounterParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}

	var event encounters.Event

	err = s.db(c).
		Where(&encounters.Event{EncounterID: encId}).
		Preload("Instances").
		Find(&event).
		Error
	if err != nil {
		s.errAbort(c, err)
		return
	}
	c.JSON(http.StatusOK, event)
}

func (s *Server) handleSetEvent(c *gin.Context) {
	_, encId, err := encounterParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}

	var event encounters.Event
	err = c.Bind(&event)
	if err != nil {
		return
	}
	event.EncounterID = uint(encId)

	err = s.db(c).
		Save(&event).
		Error
	if err != nil {
		s.errAbort(c, err)
		return
	}
	c.JSON(http.StatusOK, event)
}
