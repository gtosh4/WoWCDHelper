package events

type FilterFunc func(*Event) bool

func Filter(e []Event, f FilterFunc) (filtered []Event) {
	for i := range e {
		if f(&e[i]) {
			filtered = append(filtered, e[i])
		}
	}
	return filtered
}

func MultiFilter(fs ...FilterFunc) FilterFunc {
	return func(ev *Event) bool {
		for _, f := range fs {
			if !f(ev) {
				return false
			}
		}
		return true
	}
}

func SourceFilter(source int64) FilterFunc {
	return func(ev *Event) (equal bool) {
		h := &EventHandler{
			Absorbed:      func(event *Event, data *Absorbed) { equal = (data.SourceID == source) },
			CombatantInfo: func(event *Event, data *CombatantInfo) { equal = (data.SourceID == source) },
			Damage:        func(event *Event, data *Damage) { equal = (data.SourceID == source) },
			Heal:          func(event *Event, data *Heal) { equal = (data.SourceID == source) },
			Unknown:       func(event *Event, data *Unknown) { equal = false },
		}
		h.Handle(ev)
		return
	}
}

func TargetFilter(target int64) FilterFunc {
	return func(ev *Event) (equal bool) {
		h := &EventHandler{
			Absorbed:      func(event *Event, data *Absorbed) { equal = (data.TargetID == target) },
			CombatantInfo: func(event *Event, data *CombatantInfo) { equal = false },
			Damage:        func(event *Event, data *Damage) { equal = (data.TargetID == target) },
			Heal:          func(event *Event, data *Heal) { equal = (data.TargetID == target) },
			Unknown:       func(event *Event, data *Unknown) { equal = false },
		}
		h.Handle(ev)
		return
	}
}
