package warcraftlogs

import (
	"fmt"
	"net/url"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/gtosh4/WoWCDHelper/pkg/warcraftlogs/apitypes"
	"github.com/gtosh4/WoWCDHelper/pkg/warcraftlogs/events"
	"github.com/gtosh4/WoWCDHelper/pkg/warcraftlogs/fight"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (c *Client) AllEvents(code string, fight *fight.Fight) (fs []events.Event, err error) {
	return c.Events(code, fight, "")
}

func (c *Client) TypeEvents(code string, fight *fight.Fight, etype string) (fs []events.Event, err error) {
	return c.Events(code, fight, fmt.Sprintf(`type="%s"`, etype))
}

func (c *Client) SourceEvents(code string, fight *fight.Fight, guid int) (fs []events.Event, err error) {
	return c.Events(code, fight, fmt.Sprintf(`source.id="%d"`, guid))
}

func (c *Client) TargetEvents(code string, fight *fight.Fight, guid int) (fs []events.Event, err error) {
	return c.Events(code, fight, fmt.Sprintf(`target.id="%d"`, guid))
}

func (c *Client) Events(code string, fight *fight.Fight, filter string) (fs []events.Event, err error) {
	key := []byte(fmt.Sprintf("events_%s_%d_%s", code, fight.Id, filter))
	log := c.Log.WithFields(logrus.Fields{
		"code":      code,
		"fight":     fight.Id,
		"filter":    filter,
		"cache-key": string(key),
	})

	err = c.DB.Update(func(txn *badger.Txn) error {
		item, rerr := txn.Get(key)
		if rerr == nil {
			var evts events.Events
			rerr = item.Value(func(val []byte) error {
				return (&evts).Unmarshal(val)
			})
			if rerr != nil {
				log.WithError(rerr).Warnf("Could not unmarshal cached value")
			} else {
				fs = evts.Events
				return nil
			}
		}
		fs, rerr = c.getEventsInternal(code, fight, filter)
		if rerr != nil {
			return rerr
		}

		entry := &badger.Entry{
			Key: key,
		}
		entry.Value, rerr = (&events.Events{Events: fs}).Marshal()
		if rerr = txn.SetEntry(entry); rerr != nil {
			log.WithError(rerr).Infof("Cache put failed")
		}
		return nil
	})

	return
}

func (c *Client) getEventsInternal(code string, fight *fight.Fight, filter string) (fs []events.Event, err error) {
	log := c.Log.WithFields(logrus.Fields{"code": code, "fight": fight.Id, "filter": filter})

	fnStart := time.Now()
	pages := 0
	defer func() {
		fnEnd := time.Now()
		c.Metrics.AddSample([]string{"events", "time"}, float32(fnEnd.Sub(fnStart).Seconds()))
		c.Metrics.AddSample([]string{"events", "pages"}, float32(pages))
		c.Metrics.AddSample([]string{"events", "count"}, float32(len(fs)))
	}()

	start, end := fight.StartTime, fight.EndTime
	var evPage apitypes.EventsResponse

	logProgress := func() {
		total := fight.EndTime - fight.StartTime
		current := start - fight.StartTime
		pct := float64(current) / float64(total)
		log.Infof("[page %d] %.2f%% events completed (%d / %d) next: %d", pages, pct*100, current, total, evPage.NextPageTimestamp)
	}

	log.Infof("Fetching events using start %d to end %d", fight.StartTime, fight.EndTime)
	query := make(url.Values)
	query.Set("end", fmt.Sprintf("%d", end))
	if filter != "" {
		query.Set("filter", filter)
	}
	for ; ; pages++ {
		query.Set("start", fmt.Sprintf("%d", start))
		err = c.get("/report/events/"+code, query, &evPage)
		if err != nil {
			break
		}
		if len(evPage.Events) == 0 {
			if pages == 0 {
				err = errors.Errorf("No events found for %s/%d", code, fight.Id)
			}
			break
		}
		normalizeTimestamps(fight.StartTime, evPage.Events)

		fs = append(fs, evPage.Events...)

		logProgress()

		if evPage.NextPageTimestamp == start || evPage.NextPageTimestamp > end || evPage.NextPageTimestamp <= 0 {
			break
		}
		start = evPage.NextPageTimestamp
	}
	return
}

func normalizeTimestamps(startTime int64, evts []events.Event) {
	for i := range evts {
		evts[i].Time -= startTime
	}
}
