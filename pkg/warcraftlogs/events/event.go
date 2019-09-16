package events

import "time"

func (e *Event) TimeOffset() time.Duration {
	return time.Duration(e.Time) * time.Millisecond
}
