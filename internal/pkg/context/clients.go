package context

import (
	"github.com/FuzzyStatic/blizzard/v2"
	"github.com/allegro/bigcache"
	"go.uber.org/zap"
)

type Clients struct {
	Log   *zap.Logger
	Blizz *blizzard.Client
	Cache *bigcache.BigCache
}
