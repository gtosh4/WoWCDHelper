package analysis

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/tools/container/intsets"

	"github.com/gtosh4/WoWCDHelper/pkg/warcraftlogs/events"
)

type (
	RaidHealth []Health

	HealthTS struct {
		m  map[time.Duration]*Health
		ts *intsets.Sparse
		mu sync.Mutex
	}

	Health struct {
		Timestamp   time.Duration
		Current     int64
		Max         int64
		Healing     int64
		DamageTaken int64
	}

	totalHealth struct {
		last map[int64]*Health
	}
)

func (t *THCDA) RaidHealth(fd *FightData) (RaidHealth, error) {
	ts := NewHealthTS()

	friendlyGuids := make([]string, len(fd.Friendlies))
	players := make(map[int64]*HealthTS)
	for i, f := range fd.Friendlies {
		friendlyGuids[i] = fmt.Sprintf("%d", f.Guid)
		players[f.Id] = NewHealthTS()
	}
	start := time.Now()
	hpEvents, err := t.WCL.Events(fd.Code, fd.Fight, fmt.Sprintf(`type in ("heal", "damage") and target.id in (%s)`, strings.Join(friendlyGuids, ",")))
	if err != nil {
		return nil, errors.Wrap(err, "Error fetching events")
	}
	t.Log.Infof("Loaded events in %s", time.Since(start))
	start = time.Now()

	timestamps := new(intsets.Sparse)
	th := &totalHealth{last: make(map[int64]*Health)}

	type playerHP struct {
		Player int64
		Health Health
	}
	healthDone := make(chan struct{})
	healthCh := make(chan *playerHP, 1000)
	go func() {
		for php := range healthCh {
			players[php.Player].Add(php.Health)
			timestamps.Insert(int(php.Health.Timestamp))

			// Initialize th.last with the earliest known health values
			if lasthp := th.last[php.Player]; lasthp == nil || php.Health.Timestamp < lasthp.Timestamp {
				th.last[php.Player] = &Health{
					Timestamp: php.Health.Timestamp,
					Current:   php.Health.Current,
					Max:       php.Health.Max,
				}
			}
		}
		close(healthDone)
	}()

	blacklist := ignoredAbilitiesByBoss[fd.Fight.Boss]

	handler := &events.EventHandler{}
	handler.Heal = func(ev *events.Event, data *events.Heal) {
		php := &playerHP{
			Player: data.TargetID,
			Health: Health{
				Timestamp: ev.TimeOffset(),
				Current:   data.HitPoints,
				Max:       data.MaxHitPoints,
				Healing:   data.Amount,
			},
		}
		healthCh <- php
	}
	handler.Damage = func(ev *events.Event, data *events.Damage) {
		if _, blacklisted := blacklist[data.Ability.Guid]; blacklisted {
			return
		}
		php := &playerHP{
			Player: data.TargetID,
			Health: Health{
				Timestamp:   ev.TimeOffset(),
				Current:     data.HitPoints,
				Max:         data.MaxHitPoints,
				DamageTaken: data.UnmitigatedAmount,
			},
		}
		healthCh <- php
	}

	handler.HandleMany(hpEvents)
	close(healthCh)
	<-healthDone

	t.Log.Infof("Parsed %d events in %s", len(hpEvents), time.Since(start))
	start = time.Now()

	var (
		timestamp int
	)
	for i := 0; timestamps.TakeMin(&timestamp); i++ {
		for player, php := range players {
			if hp, ok := php.m[time.Duration(timestamp)]; ok {
				th.last[player] = hp
			}
		}

		if len(th.last) == len(players) {
			thp := th.Total(time.Duration(timestamp))
			ts.Add(thp)
		}
	}

	var (
		cur         = &window{width: float64(time.Second)}
		max         = &window{width: float64(time.Second)}
		healing     = &window{width: float64(time.Second)}
		damageTaken = &window{width: float64(time.Second)}
	)
	stamps := new(intsets.Sparse)
	stamps.Copy(ts.ts)
	maxSec := math.Ceil(time.Duration(stamps.Max()).Seconds())
	rh := make(RaidHealth, int(maxSec))

	var evstamp, prevstamp int
	for bucket := time.Duration(0); bucket.Seconds() < maxSec; bucket += time.Second {
		cutoff := bucket + time.Second
		count := 0
		for ; stamps.TakeMin(&evstamp) && time.Duration(evstamp) < cutoff; prevstamp = evstamp {
			count++
			hp := ts.m[time.Duration(evstamp)]
			width := time.Duration(evstamp - prevstamp)
			if hp.Current > 0 {
				cur.Add(width, float64(hp.Current))
				max.Add(width, float64(hp.Max))
			}
			healing.Add(width, float64(hp.Healing))
			damageTaken.Add(width, float64(hp.DamageTaken))
		}
		if count == 0 {
			width := cutoff - time.Duration(prevstamp)
			cur.Add(width, 0)
			max.Add(width, 0)
			healing.Add(width, 0)
			damageTaken.Add(width, 0)
		}
		if time.Duration(evstamp) >= cutoff {
			// Put it back so it can go into the next bucket
			stamps.Insert(evstamp)
		}
		secs := int(bucket.Seconds())
		rh[secs] = Health{
			Timestamp:   bucket,
			Current:     int64(cur.value),
			Max:         int64(max.value),
			Healing:     int64(healing.value),
			DamageTaken: int64(damageTaken.value),
		}
	}
	t.Log.Infof("Summarized events in %s", time.Since(start))
	start = time.Now()

	return rh, nil
}

func (th *totalHealth) Total(timestamp time.Duration) (hp Health) {
	hp.Timestamp = timestamp
	for _, php := range th.last {
		hp.Current += php.Current
		hp.Max += php.Max
		hp.Healing += php.Healing
		hp.DamageTaken += php.DamageTaken
	}
	return
}

func NewHealthTS() *HealthTS {
	return &HealthTS{
		m:  make(map[time.Duration]*Health),
		ts: new(intsets.Sparse),
	}
}

func (hts *HealthTS) Add(hp Health) {
	if current, ok := hts.m[hp.Timestamp]; ok {
		current.Healing += hp.Healing
		current.DamageTaken += hp.DamageTaken
	} else {
		hts.m[hp.Timestamp] = &hp
	}
	hts.ts.Insert(int(hp.Timestamp))
}

func (hp *Health) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"time_sec":     int(hp.Timestamp.Seconds()),
		"current":      hp.Current,
		"max":          hp.Max,
		"healing":      hp.Healing,
		"damage_taken": hp.DamageTaken,
	})
}
