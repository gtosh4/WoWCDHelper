package clients

import (
	"github.com/FuzzyStatic/blizzard/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Clients struct {
	Log   *zap.Logger
	Blizz *blizzard.Client
	DB    *gorm.DB
}
