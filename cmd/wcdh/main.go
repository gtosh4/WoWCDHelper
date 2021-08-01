package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/FuzzyStatic/blizzard/v2"
	"github.com/allegro/bigcache"
	"github.com/bwmarrin/snowflake"
	"github.com/gregjones/httpcache"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
	"moul.io/zapgorm2"

	"github.com/gtosh4/WoWCDHelper/internal/app"
	"github.com/gtosh4/WoWCDHelper/internal/pkg/clients"
	"github.com/gtosh4/WoWCDHelper/internal/pkg/node"
	"github.com/gtosh4/WoWCDHelper/pkg/teams"
	"github.com/gtosh4/WoWCDHelper/pkg/wow"
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
		zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), zapcore.AddSync(os.Stdout), logLevel),
	)

	node.Snowflake, err = snowflake.NewNode(int64(cfg.NodeID))
	if err != nil {
		err = errors.Wrap(err, "could not create snowflake node")
		return
	}

	dbLog := zapgorm2.New(clients.Log)
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

	clients.Blizz = blizzard.NewClient(
		cfg.BlizzClientID,
		cfg.BlizzClientSecret,
		blizzard.US,
		blizzard.EnUS,
	)

	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Hour))
	if err != nil {
		err = errors.Wrap(err, "could not create cache")
		return
	}

	clients.Blizz.SetHTTPClient(wrapBlizzHTTP(clients.Log, cache, clients.Blizz.GetHTTPClient()))

	// initDB requires Blizz API, so make sure this runs after
	err = initDB(clients)
	if err != nil {
		err = errors.Wrap(err, "could not init DB")
		return
	}

	srv := app.NewServer(clients)

	return srv.Run(fmt.Sprintf(":%d", cfg.Port))
}

func wrapBlizzHTTP(log *zap.Logger, cache *bigcache.BigCache, h *http.Client) *http.Client {
	transport := httpcache.NewTransport(&app.HTTPBigCache{
		Log:   log.Sugar().Named("cache"),
		Cache: cache,
	})
	transport.Transport = h.Transport
	h.Transport = transport
	return h
}

func initDB(clients *clients.Clients) error {
	db := clients.DB
	err := db.AutoMigrate(
		&teams.Member{},
		&teams.Team{},
	)
	if err != nil {
		return errors.Wrap(err, "error migrating tables")
	}

	ctx := context.Background()

	iconForClass := func(class string) string {
		id, err := wow.ClassNameToID(ctx, clients, class)
		if err != nil {
			return ""
		}
		media, _, err := clients.Blizz.WoWPlayableClassMedia(ctx, id)
		if err != nil {
			return ""
		}

		for _, asset := range media.Assets {
			if asset.Key == "icon" {
				return asset.Value
			}
		}
		return ""
	}

	testTeam := teams.Team{ID: "test"}

	test := teams.Roster{
		{Team: testTeam, Name: "Tosh", ClassName: "Shaman", ClassIcon: iconForClass("shaman")},
		{Team: testTeam, Name: "Jess", ClassName: "Priest", ClassIcon: iconForClass("priest")},
	}

	err = db.Create(&test).Error
	if err != nil {
		return errors.Wrap(err, "error creating test roster")
	}

	test = append(test, teams.Member{
		Team: testTeam, Name: "Sci", ClassName: "Paladin", ClassIcon: iconForClass("paladin"),
	})

	err = db.Create(&test).Error
	if err != nil {
		return errors.Wrap(err, "error appending test roster")
	}
	return nil
}
