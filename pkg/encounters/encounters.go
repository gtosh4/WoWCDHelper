package encounters

import (
	"github.com/gtosh4/WoWCDHelper/pkg/teams"
)

type (
	Encounter struct {
		ID     uint       `json:"id" gorm:"primaryKey;autoIncrement:true"`
		TeamID string     `json:"team"`
		Team   teams.Team `json:"-"`
		Name   string     `json:"name"`
		Events []Event    `json:"events"`
	}

	Roster struct {
		EncounterID uint         `json:"encounter_id" gorm:"primaryKey;foreignKey"`
		Encounter   Encounter    `json:"-"`
		MemberID    uint         `json:"member_id" gorm:"primaryKey;foreignKey"`
		Member      teams.Member `json:"-"`
		SpecID      int          `json:"spec"`
	}

	Event struct {
		ID          uint            `json:"id" gorm:"primaryKey;autoIncrement:true"`
		Label       string          `json:"label"`
		EncounterID uint            `json:"-"`
		Color       string          `json:"color"`
		Instances   []EventInstance `json:"instances"`
	}

	EventInstance struct {
		ID        uint `json:"id" gorm:"primaryKey;autoIncrement:true"`
		EventID   uint `json:"event_id"`
		OffsetSec uint `json:"offset_sec"`
	}

	Assignment struct {
		ID              uint          `json:"id" gorm:"primaryKey;autoIncrement:true"`
		EventInstanceID uint          `json:"event_instance_id"`
		EventInstance   EventInstance `json:"-"`
		MemberID        uint          `json:"member"`
		Member          teams.Member  `json:"-"`
		SpellID         uint          `json:"spell_id,omitempty"`
	}
)
