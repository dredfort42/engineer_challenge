package config

import (
	"encoding/json"
	"errors"
	"os"
)

// Config is the configuration struct for storing the configuration of the daemon
type Config struct {
	Debug             bool   `json:"debug"`
	NATSURL           string `json:"nats_url"`
	Subject           string `json:"subject"`
	InfluxURL         string `json:"influxdb_url"`
	InfluxOrg         string `json:"influxdb_org"`
	InfluxBucket      string `json:"influxdb_bucket"`
	InfluxMeasurement string `json:"influxdb_measurement"`
	PathToInfluxToken string `json:"path_to_influxdb_token"`
}

// WriterConfig is the configuration for the daemon
var WriterConfig Config

// ReadConfig reads the configuration from the file
func ReadConfig(path string) error {
	// Read the configuration from the file
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Unmarshal the configuration
	err = json.Unmarshal(data, &WriterConfig)
	if err != nil {
		return err
	}

	// Check if the configuration is valid
	if WriterConfig.NATSURL == "" ||
		WriterConfig.Subject == "" ||
		WriterConfig.InfluxURL == "" ||
		WriterConfig.InfluxOrg == "" ||
		WriterConfig.InfluxBucket == "" ||
		WriterConfig.InfluxMeasurement == "" ||
		WriterConfig.PathToInfluxToken == "" {
		return errors.New("invalid configuration")
	}

	// Check PathToInfluxToken
	if _, err := os.Stat(WriterConfig.PathToInfluxToken); os.IsNotExist(err) {
		return errors.New("influxdb token file does not exist")
	}

	return nil
}
