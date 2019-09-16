package analysis

import (
	"math"
	"time"

	"golang.org/x/tools/container/intsets"

	"github.com/gtosh4/WoWCDHelper/pkg/warcraftlogs/events"
)

type (
	PerSecondSeries []float64

	window struct {
		width float64
		value float64
	}
)

func EventsToPerSecond(evts []events.Event, f func(*events.Event) float64) PerSecondSeries {
	m := make(map[time.Duration]float64)
	for i := range evts {
		ev := &evts[i]
		t := ev.TimeOffset()
		m[t] += f(ev)
	}
	return ToPerSecond(m)
}

func ToPerSecond(values map[time.Duration]float64) PerSecondSeries {
	stamps := new(intsets.Sparse)
	for s := range values {
		stamps.Insert(int(s))
	}
	length := time.Duration(math.Ceil(time.Duration(stamps.Max()).Seconds())) * time.Second
	ps := make(PerSecondSeries, int(length.Seconds()))
	var evstamp, prevstamp int
	w := &window{width: float64(time.Second)}
	for t := time.Duration(0); t < length; t += time.Second {
		cutoff := t + time.Second
		hasData := false
		for ; stamps.TakeMin(&evstamp) && time.Duration(evstamp) < cutoff; prevstamp = evstamp {
			if evstamp > prevstamp {
				hasData = true
				w.Add(time.Duration(evstamp-prevstamp), values[time.Duration(evstamp)])
			}
		}
		if !hasData {
			w.Add(cutoff-time.Duration(prevstamp), 0)
		}
		if time.Duration(evstamp) >= cutoff {
			// Put it back so it can go into the next bucket
			stamps.Insert(evstamp)
		}
		secs := int(t.Seconds())
		ps[secs] = w.value
	}
	return ps
}

func (e PerSecondSeries) Value(t time.Duration) float64 {
	ms := int(t.Seconds())
	if len(e) < ms {
		return 0
	}
	return e[ms]
}

func (w *window) Add(dt time.Duration, v float64) {
	alpha := 1.0 - math.Exp(-float64(dt)/w.width)
	w.value += alpha * (v - w.value)
}
