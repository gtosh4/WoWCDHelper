package teams

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/gtosh4/WoWCDHelper/internal/pkg/node"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type (
	Member struct {
		ID      uint         `json:"id" gorm:"primaryKey;autoIncrement:true"`
		TeamID  string       `json:"team"`
		Team    Team         `json:"-" gorm:"foreignKey:TeamID"`
		Name    string       `json:"name"`
		ClassID int          `json:"classId"`
		Config  MemberConfig `json:"config"`
	}

	MemberConfig struct {
		MemberID    uint  `json:"-"`
		Specs       Specs `json:"specs"`
		PrimarySpec int   `json:"primarySpec"`
	}

	Team struct {
		ID         string    `json:"id" gorm:"primaryKey"`
		Name       string    `json:"name"`
		LastViewed time.Time `json:"-" gorm:"index"`
	}

	Specs []int
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

func (s *Specs) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.Errorf("failed to unmarshal spec list value: %v", value)
	}

	err := json.Unmarshal(bytes, s)
	if err != nil {
		return errors.Wrap(err, "could not parse spec list")
	}

	return nil
}

func (s Specs) Value() (driver.Value, error) {
	return json.Marshal(s)
}
