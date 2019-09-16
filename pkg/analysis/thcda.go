package analysis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"

	"github.com/gtosh4/WoWCDHelper/internal/pkg/ctx"
	"github.com/gtosh4/WoWCDHelper/pkg/warcraftlogs"
	"github.com/gtosh4/WoWCDHelper/pkg/warcraftlogs/events"
	"github.com/gtosh4/WoWCDHelper/pkg/warcraftlogs/fight"
)

type (
	THCDA struct {
		*ctx.Ctx

		WCL *warcraftlogs.Client
	}

	FightData struct {
		*ctx.Ctx

		Code string
		Num  int64

		Fights fight.Fights

		Fight      *fight.Fight
		Friendlies []fight.Friendly
		Enemies    []fight.Enemy
		Phases     *fight.Phases
	}
)

func (t *THCDA) RegisterHTTPHandlers(router *mux.Router) {
	// Backend
	apiR := router.PathPrefix("/thcda/v1").Subrouter()
	apiR.Use(
		handlers.CORS(handlers.AllowedOrigins([]string{"*"})),
		handlers.CompressHandler,
	)

	fightR := apiR.PathPrefix("/{code}/{fight:[0-9]+}").Subrouter()

	// Raw-type endpoints
	fightR.Path("/events").Methods("GET").HandlerFunc(t.handleGetEvents)
	fightR.Path("/events/{type}").Methods("GET").HandlerFunc(t.handleGetEvents)

	// Logic endpoints
	fightR.Path("/raid_health").Methods("GET").HandlerFunc(t.handleGetRaidHealth)
	fightR.Path("/damage_taken").Methods("GET").HandlerFunc(t.handleDamageTaken)
}

func (t *THCDA) fightData(reqCtx *ctx.Ctx, vars map[string]string) (fd *FightData, err error) {
	fd = new(FightData)
	fd.Code = vars["code"]
	var n int64
	n, err = strconv.ParseInt(vars["fight"], 10, 64)
	if err != nil {
		return
	}
	fd.Num = n
	fd.Ctx = reqCtx.NewSubcontextWithFields(logrus.Fields{"code": fd.Code, "fight_num": fd.Num})

	fights, err := t.WCL.Fights(fd.Code)
	if err != nil {
		err = errors.Errorf("Failed to load fights from %s", fd.Code)
		return
	}

	fd.Fight, fd.Friendlies, fd.Enemies, fd.Phases = fights.Fight(fd.Num)
	if fd.Fight == &fight.FightNotFound {
		err = errors.Errorf("Could not find fight %d in report %s", fd.Num, fd.Code)
		return
	}
	fd.Log.Infof("Found %d friends, %d enemies and phases: %+v", len(fd.Friendlies), len(fd.Enemies), fd.Phases)

	return
}

func (t *THCDA) handleGetEvents(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	etype := vars["type"]

	fd, err := t.fightData(t.Ctx.NewSubcontextWithField("remote-addr", req.RemoteAddr), vars)
	if err != nil {
		t.Log.WithError(err).Warnf("Could not get fight data for report %s", fd.Code)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	var events []events.Event
	if etype == "" {
		events, err = t.WCL.AllEvents(fd.Code, fd.Fight)
	} else {
		events, err = t.WCL.TypeEvents(fd.Code, fd.Fight, etype)
	}
	if err != nil {
		t.Log.WithError(err).Warnf("Could not get events for report %s", fd.Code)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	evFile, err := os.OpenFile(fmt.Sprintf("examples/events_%s_%d.yaml", fd.Code, fd.Num), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fd.Log.WithError(err).Warn("Could not open file")
	} else {
		defer evFile.Close()
		yaml.NewEncoder(evFile).Encode(events)
	}

	resp.WriteHeader(http.StatusNoContent)
	return
}

func (t *THCDA) handleGetRaidHealth(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	fd, err := t.fightData(t.Ctx, vars)
	if err != nil {
		t.Log.WithError(err).Warnf("Could not get fight data for report %s", fd.Code)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	raidHealth, err := t.RaidHealth(fd)
	if err != nil {
		fd.Log.WithError(err).Warnf("Could not get raid health for report %s", fd.Code)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(resp).Encode(raidHealth)
	resp.Header().Set("Content-Type", "application/json")
}

func (t *THCDA) handleDamageTaken(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	fd, err := t.fightData(t.Ctx, vars)
	if err != nil {
		t.Log.WithError(err).Warnf("Could not get fight data for report %s", fd.Code)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	raidHealth, err := t.RaidHealth(fd)
	if err != nil {
		fd.Log.WithError(err).Warnf("Could not get raid health for report %s", fd.Code)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(resp).Encode(raidHealth)
	resp.Header().Set("Content-Type", "application/json")
}
