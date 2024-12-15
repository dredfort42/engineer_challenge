package influxdb

import (
	"context"
	"log"
	"os"
	cfg "reader/internal/config"
	e "reader/internal/event"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func RequestEventsFromInfluxDB(lastEventsCount, minCriticality int, ctx context.Context, cancel context.CancelFunc) ([]e.Event, error) {
	// Read the token from the file provided in the configuration
	token, err := os.ReadFile(cfg.ReaderConfig.PathToInfluxToken)
	if err != nil {
		log.Printf("Error reading influxdb token file: %v", err)
		cancel()
		return nil, err
	}

	// Create a new client using the InfluxDB URL and token
	client := influxdb2.NewClient(cfg.ReaderConfig.InfluxURL, string(token))
	defer client.Close()

	// Create a new query client for the desired organization
	queryAPI := client.QueryAPI(cfg.ReaderConfig.InfluxOrg)

	// Create a Flux query to get the last events with a criticality higher than the provided value
	query := `
		criticality_data = from(bucket: "` + cfg.ReaderConfig.InfluxBucket + `")
			|> range(start: -1h) // Adjust the range as needed
			|> filter(fn: (r) => r._measurement == "event" and r._field == "criticality")

		event_message_data = from(bucket: "` + cfg.ReaderConfig.InfluxBucket + `")
			|> range(start: -1h) // Adjust the range as needed
			|> filter(fn: (r) => r._measurement == "event" and r._field == "event_message")

		joined_data = join(
			tables: {criticality: criticality_data, event_message: event_message_data},
			on: ["_time", "_measurement"],
			method: "inner"
		)
		
		filtered_data = joined_data
			|> filter(fn: (r) => r._value_criticality > ` + strconv.Itoa(minCriticality) + `)
			|> keep(columns: ["_time", "_value_criticality", "_value_event_message"])
			|> sort(columns: ["_time"], desc: true)
			|> limit(n: ` + strconv.Itoa(lastEventsCount) + `)

		|> yield(name: "filtered_data")
	`

	// Execute the query
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var events []e.Event

	// Loop to get the events from the result
	for result.Next() {
		rEvent := result.Record().Values()

		event := e.Event{
			Timestamp:    rEvent["_time"].(time.Time).Format(time.RFC3339),
			Criticality:  int(rEvent["_value_criticality"].(int64)),
			EventMessage: rEvent["_value_event_message"].(string),
		}

		events = append(events, event)
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	return events, nil
}
