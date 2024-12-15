package influxdb

import (
	"context"
	"log"
	"os"
	"time"
	cfg "writer/internal/config"
	e "writer/internal/event"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// WriteEventsToInfluxDB writes the event data into InfluxDB
func WriteEventsToInfluxDB(events chan e.Event, ctx context.Context, cancel context.CancelFunc) {
	defer close(events)

	// Read the token from the file provided in the configuration
	token, err := os.ReadFile(cfg.WriterConfig.PathToInfluxToken)
	if err != nil {
		log.Printf("Error reading influxdb token file: %v", err)
		cancel()
		return
	}

	// Create a new client using the InfluxDB URL and token
	client := influxdb2.NewClient(cfg.WriterConfig.InfluxURL, string(token))
	defer client.Close()

	// User blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(cfg.WriterConfig.InfluxOrg, cfg.WriterConfig.InfluxBucket)

	// Loop to write events to InfluxDB
	for {
		select {
		case <-ctx.Done():
			return

		case event := <-events:
			// Parse the timestamp of the event
			eventTime, err := time.Parse(time.RFC3339, event.Timestamp)
			if err != nil {
				log.Printf("Error parsing timestamp: %v", err)
				continue
			}

			// Create a new point with the event data
			point := influxdb2.NewPointWithMeasurement(cfg.WriterConfig.InfluxMeasurement).
				AddField("criticality", event.Criticality).
				AddField("event_message", event.EventMessage).
				SetTime(eventTime)

			// Write the point to InfluxDB
			if err := writeAPI.WritePoint(ctx, point); err != nil {
				log.Printf("Error writing point to InfluxDB: %v", err)
			} else if cfg.WriterConfig.Debug {
				log.Println("Event written to InfluxDB")
			}

		default:
			client.Ready(ctx)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
