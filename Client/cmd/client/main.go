package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"
	"time"

	cfg "client/internal/config"
	e "client/internal/event"

	"github.com/nats-io/nats.go"
)

// NATS request structure
type NATSRequest struct {
	LastEventsCount int `json:"last_events_count"`
	MinCriticality  int `json:"min_criticality"`
}

// NATS Response structure
type NATSResponse struct {
	Events []e.Event `json:"events"`
}

func main() {
	log.Println("[STATE] starting")

	// Read the minimum criticality from the environment
	minCriticalityStr := os.Getenv("MIN_CRITICALITY")
	if minCriticalityStr == "" {
		log.Fatalln("MIN_CRITICALITY environment variable is not set")
	}

	// Convert the minimum criticality to an integer
	minCriticality, err := strconv.Atoi(minCriticalityStr)
	if err != nil {
		log.Fatalf("Error converting MIN_CRITICALITY to int: %v", err)
	}

	// Read the configuration from the file
	if err := cfg.ReadConfig("config.json"); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	if cfg.ClientConfig.Debug {
		log.Printf("Debug is enabled\nNATS URL: %s\nSubject: %s\nEvent Frequency: %d ms\n",
			cfg.ClientConfig.NATSURL, cfg.ClientConfig.Subject, cfg.ClientConfig.EventFrequencyMs)
	}

	log.Println("[STATE] initializing")

	// Connect to NATS server using the URL from the config
	natsConnection, err := nats.Connect(cfg.ClientConfig.NATSURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer natsConnection.Close()

	// Create a eventTicker to generate events based on the configured frequency
	eventTicker := time.NewTicker(time.Duration(cfg.ClientConfig.EventFrequencyMs) * time.Millisecond)
	defer eventTicker.Stop()

	// Create a stateTicker to log the state of the daemon every minute
	stateTicker := time.NewTicker(1 * time.Minute)
	defer stateTicker.Stop()

	log.Println("[STATE] running")

	// Create a context to manage the lifecycle of the event loop and listen for termination signals
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Create a NATS request to get the last 10 events with criticality greater than the minimum criticality
	request := NATSRequest{
		LastEventsCount: 10,
		MinCriticality:  minCriticality,
	}

	// Marshal the request into JSON
	requestData, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("Error marshaling request: %v", err)
	}

	printedEvents := []e.Event{}

	// Generate a random event loop
	for {
		select {
		case <-eventTicker.C:
			// Send a request and wait for a response
			response, err := natsConnection.Request(cfg.ClientConfig.Subject, requestData, 30*time.Second)
			if err != nil {
				log.Printf("Failed to send request to subject '%s': %v", cfg.ClientConfig.Subject, err)
				continue
			}

			// Parse the response
			var responseData NATSResponse
			err = json.Unmarshal(response.Data, &responseData)
			if err != nil {
				log.Fatalf("Failed to parse response data: %v", err)
			}

			if len(responseData.Events) > 0 && !reflect.DeepEqual(printedEvents, responseData.Events) {
				e.PrintedEvents(responseData.Events)
				printedEvents = responseData.Events
			}

		case <-stateTicker.C:
			log.Println("[STATE] running")

		case <-ctx.Done():
			log.Println("[STATE] stopped")
			cancel()
			return
		}
	}
}
