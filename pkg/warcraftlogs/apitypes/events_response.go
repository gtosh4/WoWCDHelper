package apitypes

import (
	"bytes"
	"encoding/json"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gtosh4/WoWCDHelper/pkg/warcraftlogs/events"
)

type EventsResponse struct {
	Events            []events.Event `json:"events"`
	NextPageTimestamp int64          `json:"nextPageTimestamp"`
}

type basicEvent struct {
	Type string `json:"type"`
	Time int    `json:"timestamp"`
}

type responseRawEvents struct {
	Events            []json.RawMessage `json:"events"`
	NextPageTimestamp int64             `json:"nextPageTimestamp"`
}

func (e *EventsResponse) UnmarshalJSON(msg []byte) error {
	var r responseRawEvents
	err := json.Unmarshal(msg, &r)
	if err != nil {
		return err
	}
	e.NextPageTimestamp = r.NextPageTimestamp

	var base basicEvent
	e.Events = make([]events.Event, len(r.Events))
	for i, b := range r.Events {
		if err := json.Unmarshal(b, &base); err != nil {
			return err
		}
		evt := &e.Events[i]
		evt.Type = string(base.Type)
		evt.Time = int64(base.Time)

		dec := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		buf := bytes.NewBuffer(b)
		switch base.Type {
		case Types.Absorbed:
			var data events.Absorbed
			if err := dec.Unmarshal(buf, &data); err != nil {
				return err
			}
			evt.Data = &events.Event_Absorbed{Absorbed: &data}

		case Types.Combatantinfo:
			var data events.CombatantInfo
			if err := dec.Unmarshal(buf, &data); err != nil {
				return err
			}
			evt.Data = &events.Event_CombatantInfo{CombatantInfo: &data}

		case Types.Damage:
			var data events.Damage
			if err := dec.Unmarshal(buf, &data); err != nil {
				return err
			}
			evt.Data = &events.Event_Damage{Damage: &data}

		case Types.Heal:
			var data events.Heal
			if err := dec.Unmarshal(buf, &data); err != nil {
				return err
			}
			evt.Data = &events.Event_Heal{Heal: &data}

		default:
			var data events.Unknown
			if err := dec.Unmarshal(buf, &data); err != nil {
				return err
			}
			evt.Data = &events.Event_Unknown{Unknown: &data}
		}
	}
	return nil
}
