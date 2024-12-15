package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	cfg "writer/internal/config"
	e "writer/internal/event"
	idb "writer/internal/influxdb"

	"github.com/nats-io/nats.go"
)

const MAX_EVENTS_QUEUE = 1000

func main() {
	log.Println("[STATE] starting")

	// Read the configuration from the file
	if err := cfg.ReadConfig("config.json"); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	if cfg.WriterConfig.Debug {
		log.Printf("Debug is enabled\nNATS URL: %s\nSubject: %s\nInfluxDB URL: %s\nInfluxDB Org: %s\nInfluxDB Bucket: %s\nInfluxDB Measurement: %s\nPath to InfluxDB Token: %s\n",
			cfg.WriterConfig.NATSURL, cfg.WriterConfig.Subject, cfg.WriterConfig.InfluxURL, cfg.WriterConfig.InfluxOrg, cfg.WriterConfig.InfluxBucket, cfg.WriterConfig.InfluxMeasurement, cfg.WriterConfig.PathToInfluxToken)
	}

	log.Println("[STATE] initializing")

	// Connect to NATS server using the URL from the config
	natsConnection, err := nats.Connect(cfg.WriterConfig.NATSURL)
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

	// Create a channel for the received events
	eventsChannel := make(chan e.Event, MAX_EVENTS_QUEUE)

	// Subscribe to the subject
	_, err = natsConnection.Subscribe(cfg.WriterConfig.Subject, func(msg *nats.Msg) {
		if cfg.WriterConfig.Debug {
			log.Printf("Received %s: %s\n", msg.Subject, string(msg.Data))
		}

		// Store the received event in the list
		if len(eventsChannel) < MAX_EVENTS_QUEUE {
			event, err := e.GetEventFromJSON(string(msg.Data))
			if err != nil {
				log.Printf("Error parsing JSON: %v", err)
			} else {
				eventsChannel <- event
			}
		} else {
			log.Println("Messages queue is full. Dropping the message.")
		}
	})
	if err != nil {
		log.Fatalf("Error subscribing to subject: %v", err)
	}

	go idb.WriteEventsToInfluxDB(eventsChannel, ctx, cancel)

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
