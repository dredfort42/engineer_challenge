package event

import (
	"math/rand"
	"strconv"
	"time"
)

// Event is the struct for storing the event data
type Event struct {
	Criticality  int    `json:"criticality"`
	Timestamp    string `json:"timestamp"`
	EventMessage string `json:"eventMessage"`
}

// GetRandomEvent generates a random event
func GetRandomEvent() Event {
	criticality := rand.Int() - rand.Int() - rand.Intn(2)

	return Event{
		Criticality:  criticality,
		Timestamp:    time.Now().Format(time.RFC3339),
		EventMessage: "Random event with criticality " + strconv.Itoa(criticality),
	}
}
