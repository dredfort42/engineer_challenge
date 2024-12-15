package event

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

// Event is the struct for storing the event data
type Event struct {
	Criticality  int    `json:"criticality"`
	Timestamp    string `json:"timestamp"`
	EventMessage string `json:"eventMessage"`
}

func PrintedEvents(events []Event) {
	if len(events) == 0 {
		log.Println("No events to display")
		return
	}

	// Create a new tabwriter instance
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
	defer w.Flush()

	// Print table headers
	fmt.Fprintf(w, "\n%s\t%s\t%s\t%s\n", "N", "Timestamp", "Criticality", "Event Message")

	// Print table rows
	for i, e := range events {
		fmt.Fprintf(w, "%02d\t%s\t%d\t%s\n", i+1, e.Timestamp, e.Criticality, e.EventMessage)
	}

	fmt.Fprintln(w)

}
