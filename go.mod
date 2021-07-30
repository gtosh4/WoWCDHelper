module github.com/gtosh4/WoWCDHelper

go 1.16

// pending merge of https://github.com/FuzzyStatic/blizzard/pull/41
replace github.com/FuzzyStatic/blizzard v1.2.3 => github.com/gtosh4/blizzard v1.2.4

require (
	github.com/FuzzyStatic/blizzard v1.2.3
	github.com/allegro/bigcache v1.2.1
	github.com/chenjiandongx/ginprom v0.0.0-20210617023641-6c809602c38a
	github.com/gin-contrib/zap v0.0.1
	github.com/gin-gonic/gin v1.7.2
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.1.0
	github.com/spf13/cobra v1.2.1
	go.uber.org/zap v1.18.1
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
)
