package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	cfg "reader/internal/config"
	e "reader/internal/event"
	idb "reader/internal/influxdb"

	"github.com/nats-io/nats.go"
)

const MAX_EVENTS_QUEUE = 1000

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

	// Read the configuration from the file
	if err := cfg.ReadConfig("config.json"); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	if cfg.ReaderConfig.Debug {
		log.Printf("Debug is enabled\nNATS URL: %s\nSubject: %s\nInfluxDB URL: %s\nInfluxDB Org: %s\nInfluxDB Bucket: %s\nInfluxDB Measurement: %s\nPath to InfluxDB Token: %s\n",
			cfg.ReaderConfig.NATSURL, cfg.ReaderConfig.Subject, cfg.ReaderConfig.InfluxURL, cfg.ReaderConfig.InfluxOrg, cfg.ReaderConfig.InfluxBucket, cfg.ReaderConfig.InfluxMeasurement, cfg.ReaderConfig.PathToInfluxToken)
	}

	log.Println("[STATE] initializing")

	// Connect to NATS server using the URL from the config
	natsConnection, err := nats.Connect(cfg.ReaderConfig.NATSURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer natsConnection.Close()

	// Create a stateTicker to log the state of the daemon every minute
	stateTicker := time.NewTicker(1 * time.Minute)
	defer stateTicker.Stop()

	log.Println("[STATE] running")

	// Create a context to manage the lifecycle of the event loop and listen for termination signals
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Subscribe to the subject
	_, err = natsConnection.Subscribe(cfg.ReaderConfig.Subject, func(msg *nats.Msg) {
		// Parse incoming request
		var request NATSRequest
		if err := json.Unmarshal(msg.Data, &request); err != nil {
			log.Printf("Error parsing request JSON: %v", err)
		} else {
			if cfg.ReaderConfig.Debug {
				log.Printf("Received request: fetching last %d events with criticality > %d", request.LastEventsCount, request.MinCriticality)
			}

			// Fetch events from InfluxDB
			events, err := idb.RequestEventsFromInfluxDB(request.LastEventsCount, request.MinCriticality, ctx, cancel)
			if err != nil {
				log.Printf("Error fetching events from InfluxDB: %v", err)
			} else {
				// Respond with the events
				response := NATSResponse{Events: events}

				responseJSON, err := json.Marshal(response)
				if err != nil {
					log.Printf("Error marshalling response JSON: %v", err)
				} else {
					// Respond to the request
					if err := msg.Respond(responseJSON); err != nil {
						log.Printf("Error responding to request: %v", err)
					} else {
						if cfg.ReaderConfig.Debug {
							log.Printf("Responded with %d events", len(events))
						}
					}
				}
			}
		}
	})
	if err != nil {
		log.Fatalf("Error subscribing to subject: %v", err)
	}

	// Generate a random event loop
	for {
		select {
		case <-stateTicker.C:
			log.Println("[STATE] running")

		case <-ctx.Done():
			log.Println("[STATE] stopped")
			cancel()
			return
		}
	}
}
