package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/FuzzyStatic/blizzard/v3"
	"github.com/allegro/bigcache"
	"github.com/bwmarrin/snowflake"
	"github.com/gregjones/httpcache"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
	"moul.io/zapgorm2"

	"github.com/gtosh4/WoWCDHelper/internal/app"
	"github.com/gtosh4/WoWCDHelper/internal/pkg/clients"
	"github.com/gtosh4/WoWCDHelper/internal/pkg/node"
	ratelimit "github.com/gtosh4/WoWCDHelper/internal/pkg/rate_limit"
	"github.com/gtosh4/WoWCDHelper/pkg/teams"
)

var root = &cobra.Command{
	Use:  "wcdh",
	RunE: serve,
}

var cfg struct {
	Debug             bool
	Port              int
	BlizzClientID     string
	BlizzClientSecret string
	DBDriver          string
	DSN               string
	NodeID            int
}

func main() {
	flags := root.Flags()

	flags.BoolVar(&cfg.Debug, "debug", false, "enable debug logging")
	flags.IntVarP(&cfg.Port, "port", "p", 8080, "port to bind to")
	flags.StringVar(&cfg.BlizzClientID, "bnetId", "31708c8133144f6fab3b75e2ece62d3d", "Battle.net API client ID")
	flags.StringVar(&cfg.BlizzClientSecret, "bnetSecret", "", "Battle.net API client secret")
	flags.StringVar(&cfg.DBDriver, "db", "sqlite3", "database driver to use")
	flags.StringVar(&cfg.DSN, "dsn", "file::memory:?cache=shared", "database/sql DSN string")
	flags.IntVar(&cfg.NodeID, "nodeID", 1, "node id (used for snowflake UUID generation)")

	root.Execute()
}

func serve(cmd *cobra.Command, args []string) (err error) {
	clients := &clients.Clients{}

	zap.L()
	logLevel := zap.NewAtomicLevel()
	if cfg.Debug {
		logLevel.SetLevel(zap.DebugLevel)
	}
	clients.Log = zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			logLevel,
		),
	)

	node.Snowflake, err = snowflake.NewNode(int64(cfg.NodeID))
	if err != nil {
		err = errors.Wrap(err, "could not create snowflake node")
		return
	}

	dbLog := zapgorm2.New(clients.Log)
	if cfg.Debug {
		dbLog.LogLevel = gormlogger.Info
	}
	dbLog.SetAsDefault()
	dbCfg := &gorm.Config{
		Logger: dbLog,
	}

	var dbMetrics []prometheus.MetricsCollector
	switch cfg.DBDriver {
	case "sqlite3":
		clients.DB, err = gorm.Open(sqlite.Open(cfg.DSN), dbCfg)

	case "mysql":
		clients.DB, err = gorm.Open(mysql.Open(cfg.DSN), dbCfg)
		dbMetrics = []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"Threads_running"},
			},
		}

	default:
		err = errors.Errorf("unsupported db driver name: %s", cfg.DBDriver)
	}
	if err != nil {
		err = errors.Wrap(err, "could not open DB")
		return
	}
	clients.DB.Use(prometheus.New(prometheus.Config{
		DBName:           "db",
		MetricsCollector: dbMetrics,
	}))

	err = initDB(clients)
	if err != nil {
		err = errors.Wrap(err, "could not init DB")
		return
	}

	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Hour))
	if err != nil {
		err = errors.Wrap(err, "could not create cache")
		return
	}

	clients.Blizz, err = blizzard.NewClient(
		blizzard.Config{
			ClientID:     cfg.BlizzClientID,
			ClientSecret: cfg.BlizzClientSecret,
			HTTPClient:   blizzHTTPClient(clients.Log, cache),
			Region:       blizzard.US,
			Locale:       blizzard.EnUS,
		})
	if err != nil {
		err = errors.Wrap(err, "could not create blizz client")
		return
	}

	srv := app.NewServer(clients)

	return srv.Run(fmt.Sprintf(":%d", cfg.Port))
}

func blizzHTTPClient(log *zap.Logger, cache *bigcache.BigCache) *http.Client {
	transport := httpcache.NewTransport(&app.HTTPBigCache{
		Log:   log.Sugar().Named("cache"),
		Cache: cache,
	})
	transport.Transport = &ratelimit.Transport{
		Ratelimiter: rate.NewLimiter(rate.Every(90/time.Second), 90),
	}
	client := transport.Client()
	client.Timeout = 1 * time.Minute
	return client
}

func initDB(clients *clients.Clients) error {
	db := clients.DB
	err := db.AutoMigrate(
		&teams.Team{},
		&teams.MemberConfig{},
		&teams.Member{},
	)
	if err != nil {
		return errors.Wrap(err, "error migrating tables")
	}

	testTeam := teams.Team{ID: "test"}

	test := []teams.Member{
		{
			Team:    testTeam,
			Name:    "Tosh",
			ClassID: 7, /* Shaman */
			Config: teams.MemberConfig{
				Specs:       []int{264},
				PrimarySpec: 264,
			},
		},
		{
			Team:    testTeam,
			Name:    "Jess",
			ClassID: 5, /* Priest */
			Config: teams.MemberConfig{
				Specs:       []int{257},
				PrimarySpec: 257,
			},
		},
	}

	err = db.Create(&test).Error
	if err != nil {
		return errors.Wrap(err, "error creating test roster")
	}

	test = append(test, teams.Member{
		Team:    testTeam,
		Name:    "Sci",
		ClassID: 2, /* Paladin */
		Config: teams.MemberConfig{
			Specs:       []int{65, 70},
			PrimarySpec: 65,
		},
	})

	err = db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&test).Error
	if err != nil {
		return errors.Wrap(err, "error appending test roster")
	}
	return nil
}
