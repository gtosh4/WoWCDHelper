package clients

import (
	"net/http"

	"github.com/FuzzyStatic/blizzard/v3"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Clients struct {
	Log        *zap.Logger
	Blizz      *blizzard.Client
	DB         *gorm.DB
	Cache      *badger.DB
	IconClient *http.Client
}
