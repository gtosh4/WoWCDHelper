package teams

import (
	"time"

	"github.com/gtosh4/WoWCDHelper/internal/pkg/node"
	"gorm.io/gorm"
)

type (
	Member struct {
		ID        uint   `json:"id" gorm:"primaryKey;autoIncrement:true"`
		TeamID    string `json:"team"`
		Team      Team   `json:"-" gorm:"foreignKey:TeamID"`
		Name      string `json:"name"`
		ClassName string `json:"className"`
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
	if t.LastViewed.IsZero() {
		t.LastViewed = time.Now()
	}
	return nil
}
