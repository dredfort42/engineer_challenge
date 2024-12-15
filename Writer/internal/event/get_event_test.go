package event

import (
	"testing"
)

func TestGetEventFromJSON(t *testing.T) {
	tests := []struct {
		name      string
		jsonEvent string
		want      Event
		wantErr   bool
	}{
		{
			name:      "Valid JSON",
			jsonEvent: `{"criticality": 5, "timestamp": "2024-12-15T16:02:19Z", "eventMessage": "Test event"}`,
			want:      Event{Criticality: 5, Timestamp: "2024-12-15T16:02:19Z", EventMessage: "Test event"},
			wantErr:   false,
		},
		{
			name:      "Invalid JSON",
			jsonEvent: `{"criticality": "high", "timestamp": "2024-12-15T16:02:19Z", "eventMessage": "Test event"}`,
			want:      Event{},
			wantErr:   true,
		},
		{
			name:      "Empty JSON",
			jsonEvent: `{}`,
			want:      Event{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEventFromJSON(tt.jsonEvent)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventFromJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetEventFromJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
