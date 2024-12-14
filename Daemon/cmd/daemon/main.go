package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	cfg "daemon/internal/config"
	e "daemon/internal/event"

	"github.com/nats-io/nats.go"
)

func main() {
	log.Println("[STATE] starting")

	// Read the configuration from the file
	if err := cfg.ReadConfig("config.json"); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	if cfg.DeamonConfig.Debug {
		log.Printf("Debug is enabled\nNATS URL: %s\nSubject: %s\nEvent Frequency: %d ms\n",
			cfg.DeamonConfig.NATSURL, cfg.DeamonConfig.Subject, cfg.DeamonConfig.EventFrequencyMs)
	}

	log.Println("[STATE] initializing")

	// Connect to NATS server using the URL from the config
	natsConnection, err := nats.Connect(cfg.DeamonConfig.NATSURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer natsConnection.Close()

	// Create a eventTicker to generate events based on the configured frequency
	eventTicker := time.NewTicker(time.Duration(cfg.DeamonConfig.EventFrequencyMs) * time.Millisecond)
	defer eventTicker.Stop()

	// Create a stateTicker to log the state of the daemon every minute
	stateTicker := time.NewTicker(1 * time.Minute)
	defer stateTicker.Stop()

	log.Println("[STATE] running")

	// Create a context to manage the lifecycle of the event loop and listen for termination signals
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Generate a random event loop
	for {
		select {
		case <-eventTicker.C:
			// Marshal the event into JSON
			eventData, err := json.Marshal(e.GetRandomEvent())
			if err != nil {
				log.Printf("Error marshaling event: %v", err)
				continue
			}

			// Publish the event to NATS
			if err := natsConnection.Publish(cfg.DeamonConfig.Subject, eventData); err != nil {
				log.Printf("Error publishing event: %v", err)
			} else if cfg.DeamonConfig.Debug {
				log.Printf("Published event: %s", string(eventData))
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
