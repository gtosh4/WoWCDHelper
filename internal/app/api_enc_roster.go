package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gtosh4/WoWCDHelper/pkg/encounters"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Server) handleGetRoster(c *gin.Context) {
	_, encId, err := encounterParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}

	roster := []encounters.Roster{}
	err = s.db(c).
		Where(&encounters.Roster{EncounterID: encId}).
		Find(&roster).
		Error

	if err != nil {
		s.errAbort(c, err)
		return
	}
	c.JSON(http.StatusOK, roster)
}

func (s *Server) handleSetRoster(c *gin.Context) {
	teamId, encId, err := encounterParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}

	roster := []encounters.Roster{}
	err = c.Bind(&roster)
	if err != nil {
		s.log(c).Warnf("bind err: %v", err)
		return
	}
	for i := range roster {
		roster[i].EncounterID = uint(encId)
	}
	err = s.db(c).Transaction(func(tx *gorm.DB) error {
		err := tx.
			Where(&encounters.Encounter{TeamID: teamId, ID: uint(encId)}).
			Delete(&encounters.Roster{}).
			Error
		if err != nil {
			return err
		}

		err = tx.
			Create(&roster).
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
	c.JSON(http.StatusOK, roster)
}

func (s *Server) handleGetRosterMember(c *gin.Context) {
	_, encId, err := encounterParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}
	_, memberId, err := memberParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}

	var rm encounters.Roster

	err = s.db(c).
		Where(&encounters.Roster{EncounterID: encId, MemberID: memberId}).
		First(&rm).
		Error
	if err != nil {
		s.errAbort(c, err)
		return
	}

	c.JSON(http.StatusOK, rm)
}

func (s *Server) handleSetRosterMember(c *gin.Context) {
	_, encId, err := encounterParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}
	memberId, err := strconv.ParseInt(c.Param("member"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var rm encounters.Roster
	err = c.Bind(&rm)
	if err != nil {
		s.log(c).Warnf("bind err: %v", err)
		return
	}
	rm.EncounterID = encId
	rm.MemberID = uint(memberId)

	err = s.db(c).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&rm).
		Error
	if err != nil {
		s.errAbort(c, err)
		return
	}

	c.JSON(http.StatusOK, rm)
}

func (s *Server) handleDeleteRosterMember(c *gin.Context) {
	_, encId, err := encounterParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}

	memberId, err := strconv.ParseUint(c.Param("member"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		s.log(c).Warnf("params err: %v", err)
		return
	}

	err = s.db(c).
		Where(&encounters.Roster{EncounterID: encId, MemberID: uint(memberId)}).
		Delete(&encounters.Roster{}).
		Error
	if err != nil {
		s.errAbort(c, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
