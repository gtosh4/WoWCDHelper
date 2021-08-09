package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FuzzyStatic/blizzard/v3"
	"github.com/bwmarrin/snowflake"
	"github.com/dgraph-io/badger/v3"
	"github.com/gregjones/httpcache"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	gormlogger "gorm.io/gorm/logger"
	gormprom "gorm.io/plugin/prometheus"
	"moul.io/zapgorm2"

	"github.com/gtosh4/WoWCDHelper/internal/app"
	"github.com/gtosh4/WoWCDHelper/internal/pkg/clients"
	"github.com/gtosh4/WoWCDHelper/internal/pkg/node"
	ratelimit "github.com/gtosh4/WoWCDHelper/internal/pkg/rate_limit"
	"github.com/gtosh4/WoWCDHelper/pkg/encounters"
	"github.com/gtosh4/WoWCDHelper/pkg/teams"
)

var root = &cobra.Command{
	Use:  "wcdh",
	RunE: serve,
}

var cfg struct {
	Debug             bool
	ShortTime         bool
	Port              int
	CacheDir          string
	BlizzClientID     string
	BlizzClientSecret string
	DBDriver          string
	DSN               string
	NodeID            int
}

func main() {
	flags := root.Flags()

	flags.BoolVar(&cfg.Debug, "debug", false, "enable debug logging")
	flags.BoolVar(&cfg.ShortTime, "short-time", false, "enable short time in logs")
	flags.IntVarP(&cfg.Port, "port", "p", 8080, "port to bind to")
	flags.StringVar(&cfg.CacheDir, "cache", ".cache", "Directory to use for cache")
	flags.StringVar(&cfg.BlizzClientID, "bnetId", "31708c8133144f6fab3b75e2ece62d3d", "Battle.net API client ID")
	flags.StringVar(&cfg.BlizzClientSecret, "bnetSecret", "", "Battle.net API client secret")
	flags.StringVar(&cfg.DBDriver, "db", "sqlite3", "database driver to use")
	flags.StringVar(&cfg.DSN, "dsn", "file:wowcdh.db?cache=shared&mode=rwc", "database/sql DSN string")
	flags.IntVar(&cfg.NodeID, "nodeID", 1, "node id (used for snowflake UUID generation)")

	root.Execute()
}

func serve(cmd *cobra.Command, args []string) (err error) {
	c := &clients.Clients{}

	logLevel := zap.NewAtomicLevel()
	if cfg.Debug {
		logLevel.SetLevel(zap.DebugLevel)
	}
	logCfg := zap.NewDevelopmentEncoderConfig()
	logCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	if cfg.ShortTime {
		logCfg.EncodeTime = clients.ShortTimeEncoder(time.Now())
	}
	c.Log = zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(logCfg),
			zapcore.AddSync(os.Stdout),
			logLevel,
		),
	)
	zap.ReplaceGlobals(c.Log)

	node.Snowflake, err = snowflake.NewNode(int64(cfg.NodeID))
	if err != nil {
		err = errors.Wrap(err, "could not create snowflake node")
		return
	}

	dbLog := zapgorm2.New(c.Log)
	if cfg.Debug {
		dbLog.LogLevel = gormlogger.Info
	}
	dbLog.SetAsDefault()
	dbCfg := &gorm.Config{
		Logger: dbLog,
	}

	var dbMetrics []gormprom.MetricsCollector
	switch cfg.DBDriver {
	case "sqlite3":
		c.DB, err = gorm.Open(sqlite.Open(cfg.DSN), dbCfg)

	case "mysql":
		c.DB, err = gorm.Open(mysql.Open(cfg.DSN), dbCfg)
		dbMetrics = []gormprom.MetricsCollector{
			&gormprom.MySQL{
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
	c.DB.Use(gormprom.New(gormprom.Config{
		DBName:           "db",
		MetricsCollector: dbMetrics,
	}))

	err = initDB(c)
	if err != nil {
		err = errors.Wrap(err, "could not init DB")
		return
	}

	cacheOpts := badger.DefaultOptions(cfg.CacheDir).
		WithLogger(&clients.BadgerZapLogger{Log: c.Log.Sugar().Named("cache")})
	c.Cache, err = badger.Open(cacheOpts)
	if err != nil {
		err = errors.Wrap(err, "could not init cache")
		return
	}

	err = prometheus.Register(clients.NewBadgerMetricCollector(c.Cache))
	if err != nil {
		err = errors.Wrap(err, "could not register cache metrics")
		return
	}

	c.Blizz, err = blizzard.NewClient(
		blizzard.Config{
			ClientID:     cfg.BlizzClientID,
			ClientSecret: cfg.BlizzClientSecret,
			HTTPClient:   blizzHTTPClient(c),
			Region:       blizzard.US,
			Locale:       blizzard.EnUS,
		})
	if err != nil {
		err = errors.Wrap(err, "could not create blizz client")
		return
	}

	c.IconClient = iconHTTPClient(c)

	srv := app.NewServer(c, fmt.Sprintf(":%d", cfg.Port))

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-exitCh
		c.Log.Sugar().Infof("got signal %s shutting down", sig)
		if err := srv.Close(); err != nil {
			c.Log.Sugar().Warnf("error closing http server: %v", err)
		}
	}()

	err = srv.Run()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

func blizzHTTPClient(c *clients.Clients) *http.Client {
	cache := &clients.HTTPCache{
		Log:    c.Log.Sugar().Named("blizz-cache"),
		Cache:  c.Cache,
		Prefix: "blizz",
	}
	prometheus.Register(cache)
	transport := httpcache.NewTransport(cache)
	transport.Transport = &ratelimit.Transport{
		Ratelimiter: rate.NewLimiter(rate.Every(20/time.Second), 20),
	}
	client := transport.Client()
	client.Timeout = 1 * time.Minute
	return client
}

func iconHTTPClient(c *clients.Clients) *http.Client {
	cache := &clients.HTTPCache{
		Log:    c.Log.Sugar().Named("icon-cache"),
		Cache:  c.Cache,
		Prefix: "icon",
	}
	prometheus.Register(cache)
	transport := httpcache.NewTransport(cache)
	transport.Transport = &ratelimit.Transport{
		Ratelimiter: rate.NewLimiter(rate.Every(20/time.Second), 20),
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
		&encounters.EventInstance{},
		&encounters.Assignment{},
		&encounters.Event{},
		&encounters.Encounter{},
		&encounters.Roster{},
	)
	if err != nil {
		return errors.Wrap(err, "error migrating tables")
	}

	testTeam := teams.Team{ID: "test"}

	test := []teams.Member{
		{
			ID:      1,
			Team:    testTeam,
			Name:    "Tosh",
			ClassID: 7, /* Shaman */
			Config: teams.MemberConfig{
				Specs:       []int{264},
				PrimarySpec: 264,
			},
		},
		{
			ID:      2,
			Team:    testTeam,
			Name:    "Jess",
			ClassID: 5, /* Priest */
			Config: teams.MemberConfig{
				Specs:       []int{257},
				PrimarySpec: 257,
			},
		},
	}

	err = db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&test).Error
	if err != nil {
		return errors.Wrap(err, "error creating test roster")
	}

	test = append(test, teams.Member{
		ID:      3,
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
