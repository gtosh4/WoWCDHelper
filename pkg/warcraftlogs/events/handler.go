package events

import (
	"sync"
)

type EventHandler struct {
	Absorbed      func(event *Event, data *Absorbed)
	CombatantInfo func(event *Event, data *CombatantInfo)
	Damage        func(event *Event, data *Damage)
	Heal          func(event *Event, data *Heal)
	Unknown       func(event *Event, data *Unknown)
}

func (handler *EventHandler) Handle(e *Event) {
	switch ev := e.GetData().(type) {
	case *Event_Absorbed:
		if handler.Absorbed != nil {
			handler.Absorbed(e, ev.Absorbed)
		}

	case *Event_CombatantInfo:
		if handler.CombatantInfo != nil {
			handler.CombatantInfo(e, ev.CombatantInfo)
		}

	case *Event_Damage:
		if handler.Damage != nil {
			handler.Damage(e, ev.Damage)
		}

	case *Event_Heal:
		if handler.Heal != nil {
			handler.Heal(e, ev.Heal)
		}

	case *Event_Unknown:
		if handler.Unknown != nil {
			handler.Unknown(e, ev.Unknown)
		}

	default:
		// Should be unreachable
	}
}

func (handler *EventHandler) HandleMany(evts []Event) {
	eventsPerBatch := 100

	batches := len(evts) / eventsPerBatch
	if batches <= 1 {
		for i := range evts {
			handler.Handle(&evts[i])
		}
		return
	} else if batches > 10000 {
		batches = 10000
	}

	wg := &sync.WaitGroup{}

	queue := make(chan *Event, batches+1)
	for i := 0; i < batches; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for e := range queue {
				handler.Handle(e)
			}
		}()
	}

	for i := range evts {
		queue <- &evts[i]
	}
	close(queue)
	wg.Wait()
}
