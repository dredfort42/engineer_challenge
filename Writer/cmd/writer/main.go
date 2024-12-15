package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	cfg "writer/internal/config"

	"github.com/nats-io/nats.go"
)

func main() {
	log.Println("[STATE] starting")

	// Read the configuration from the file
	if err := cfg.ReadConfig("config.json"); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	if cfg.WriterConfig.Debug {
		log.Printf("Debug is enabled\nNATS URL: %s\nSubject: %s\n",
			cfg.WriterConfig.NATSURL, cfg.WriterConfig.Subject)
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

	// Subscribe to the subject
	_, err = natsConnection.Subscribe(cfg.WriterConfig.Subject, func(msg *nats.Msg) {
		if cfg.WriterConfig.Debug {
			log.Printf("Received %s: %s\n", msg.Subject, string(msg.Data))
		}
		// Process the received message

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
