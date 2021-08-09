package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gtosh4/WoWCDHelper/pkg/encounters"
	"github.com/gtosh4/WoWCDHelper/pkg/teams"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func registerTeamApi(s *Server) {
	s.router.POST("/teams", s.handleCreateTeam)

	teamR := s.router.Group("/team/:team")
	teamR.GET("/members", s.handleGetTeam)
	teamR.PUT("/members", s.handleSetTeam)
	teamR.POST("/member", s.handleNewMember)

	memberR := teamR.Group("/member/:member")
	memberR.PUT("", s.handleUpdateMember)
	memberR.DELETE("", s.handleDeleteMember)
	memberR.GET("/encounters", s.handleGetMemberEncounters)
	memberR.POST("/encounters", s.handleSetMemberEncounters)
	memberR.DELETE("/encounters", s.handlRemoveMemberEncounters)
}

func memberParams(c *gin.Context) (teamId string, memberId uint, err error) {
	teamId = c.Param("team")
	if teamId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		err = errors.New("empty team param")
		return
	}
	memberId64, err := strconv.ParseUint(c.Param("member"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	memberId = uint(memberId64)
	return
}

func (s *Server) handleGetTeam(c *gin.Context) {
	teamId := c.Param("team")
	if teamId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	members := []teams.Member{}
	err := s.db(c).
		Where("team_id = ?", teamId).
		Preload("Config").
		Find(&members).
		Error

	if err != nil {
		s.errAbort(c, err)
		return
	}
	c.JSON(http.StatusOK, members)
}

func (s *Server) handleSetTeam(c *gin.Context) {
	teamId := c.Param("team")
	if teamId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var roster []teams.Member
	err := c.Bind(&roster)
	if err != nil {
		return
	}
	for i := range roster {
		roster[i].TeamID = teamId
	}
	err = s.db(c).Transaction(func(tx *gorm.DB) error {
		err := tx.
			Where(&teams.Member{TeamID: teamId}).
			Delete(&teams.Member{}).
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

func (s *Server) handleNewMember(c *gin.Context) {
	teamId := c.Param("team")
	if teamId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var m teams.Member
	err := c.Bind(&m)
	if err != nil {
		return
	}
	m.ID = 0 // make sure we generate a new id
	m.TeamID = teamId

	err = s.db(c).
		Create(&m).
		Error
	if err != nil {
		s.errAbort(c, err)
		return
	}

	c.Redirect(http.StatusCreated, fmt.Sprintf("/team/%s/member/%d", m.TeamID, m.ID))
}

func (s *Server) handleUpdateMember(c *gin.Context) {
	teamId, memberId, err := memberParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}

	var m teams.Member
	err = c.Bind(&m)
	if err != nil {
		return
	}
	if m.ID != uint(memberId) || m.TeamID != teamId {
		c.AbortWithError(
			http.StatusBadRequest,
			errors.Errorf("{id: %d, team: %s} does not match path values {id: %d, team: %s}",
				m.ID, m.TeamID,
				memberId, teamId,
			),
		)
	}

	err = s.db(c).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&m).
		Error
	if err != nil {
		s.errAbort(c, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) handleDeleteMember(c *gin.Context) {
	teamId, memberId, err := memberParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}

	err = s.db(c).
		Delete(&teams.Member{ID: uint(memberId), TeamID: teamId}).
		Error
	if err != nil {
		s.errAbort(c, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) handleGetMemberEncounters(c *gin.Context) {
	teamId, memberId, err := memberParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}

	roster := []encounters.Roster{}
	err = s.db(c).
		Preload("Encounter").
		Where(&encounters.Roster{MemberID: memberId, Encounter: encounters.Encounter{TeamID: teamId}}).
		Find(&roster).
		Error

	if err != nil {
		s.errAbort(c, err)
		return
	}
	c.JSON(http.StatusOK, roster)
}

func (s *Server) handleSetMemberEncounters(c *gin.Context) {
	log := s.log(c)
	teamId, memberId, err := memberParams(c)
	if err != nil {
		log.Warnf("params err: %v", err)
		return
	}

	var rm encounters.Roster
	err = c.Bind(&rm)
	if err != nil {
		s.log(c).Warnf("bind err: %v", err)
		return
	}
	log.Infof("bound %+v", rm)
	rm.MemberID = memberId

	var roster []encounters.Roster
	err = s.db(c).Transaction(func(tx *gorm.DB) error {
		var es []encounters.Encounter
		err := tx.
			Where(&encounters.Encounter{TeamID: teamId}).
			Find(&es).
			Error

		if err != nil {
			return err
		}

		for _, enc := range es {
			roster = append(roster, encounters.Roster{
				EncounterID: enc.ID,
				MemberID:    rm.MemberID,
				SpecID:      rm.SpecID,
			})
		}

		log.Infof("creating %v", roster)
		return tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&roster).Error
	})
	if err != nil {
		s.errAbort(c, err)
		return
	}
	log.Infof("created %v", roster)
	c.JSON(http.StatusOK, roster)
}

func (s *Server) handlRemoveMemberEncounters(c *gin.Context) {
	teamId, memberId, err := memberParams(c)
	if err != nil {
		s.log(c).Warnf("params err: %v", err)
		return
	}

	var roster []encounters.Roster
	err = s.db(c).
		Preload("Encounter").
		Where(&encounters.Roster{MemberID: memberId, Encounter: encounters.Encounter{TeamID: teamId}}).
		Delete(&roster).
		Error

	if err != nil {
		s.errAbort(c, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
