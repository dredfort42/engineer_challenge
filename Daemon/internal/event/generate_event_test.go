package event

import (
	"strconv"
	"testing"
	"time"
)

func TestGetRandomEvent(t *testing.T) {
	event := GetRandomEvent()

	if event.Criticality == 0 {
		t.Errorf("Expected non-zero criticality, got %d", event.Criticality)
	}

	_, err := time.Parse(time.RFC3339, event.Timestamp)
	if err != nil {
		t.Errorf("Expected valid timestamp, got %s", event.Timestamp)
	}

	expectedMessage := "Random event with criticality " + strconv.Itoa(event.Criticality)
	if event.EventMessage != expectedMessage {
		t.Errorf("Expected event message %s, got %s", expectedMessage, event.EventMessage)
	}
}
