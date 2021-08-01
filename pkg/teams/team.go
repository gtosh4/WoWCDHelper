package teams

import (
	"time"

	"github.com/gtosh4/WoWCDHelper/internal/pkg/node"
	"gorm.io/gorm"
)

type (
	Member struct {
		ID        uint   `json:"id" gorm:"primaryKey"`
		Team      Team   `json:"team" gorm:"foreignKey:ID"`
		Name      string `json:"name"`
		ClassName string `json:"className"`
		ClassIcon string `json:"classIcon" gorm:"-"`
	}

	Roster []Member

	Team struct {
		ID         string    `json:"id" gorm:"primaryKey"`
		LastViewed time.Time `json:"-" gorm:"index"`
	}
)

func (t *Team) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = node.Snowflake.Generate().Base58()
	}
	return nil
}
