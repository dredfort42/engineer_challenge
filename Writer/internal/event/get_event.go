package event

import "encoding/json"

// Event is the struct for storing the event data
type Event struct {
	Criticality  int    `json:"criticality"`
	Timestamp    string `json:"timestamp"`
	EventMessage string `json:"eventMessage"`
}

// GetEventFromJSON parses a JSON string and returns an Event struct
func GetEventFromJSON(jsonEvent string) (Event, error) {
	event := Event{}
	err := json.Unmarshal([]byte(jsonEvent), &event)
	if err != nil {
		return Event{}, err
	}

	return event, nil
}
