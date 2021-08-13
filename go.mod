module github.com/gtosh4/WoWCDHelper

go 1.16

// pending merge of https://github.com/FuzzyStatic/blizzard/pull/41
// replace github.com/FuzzyStatic/blizzard v1.2.3 => github.com/gtosh4/blizzard v1.2.4

require (
	github.com/FuzzyStatic/blizzard v1.3.0
	github.com/bwmarrin/snowflake v0.3.0
	github.com/chenjiandongx/ginprom v0.0.0-20210617023641-6c809602c38a
	github.com/dgraph-io/badger/v3 v3.2103.1
	github.com/dgraph-io/ristretto v0.1.0
	github.com/gin-contrib/zap v0.0.1
	github.com/gin-gonic/gin v1.7.2
	github.com/gogo/protobuf v1.3.2
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79
	github.com/mattn/go-sqlite3 v1.14.8 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.30.0 // indirect
	github.com/prometheus/procfs v0.7.1 // indirect
	github.com/spf13/cobra v1.2.1
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.18.1
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/protobuf v1.27.1 // indirect
	gorm.io/driver/mysql v1.1.1
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.12
	gorm.io/plugin/prometheus v0.0.0-20210614014227-3996fd54c851
	moul.io/zapgorm2 v1.1.0
)
