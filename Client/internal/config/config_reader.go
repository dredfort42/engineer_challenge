package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

// Config is the configuration struct for storing the configuration of the daemon
type Config struct {
	Debug            bool   `json:"debug"`
	NATSURL          string `json:"nats_url"`
	Subject          string `json:"subject"`
	EventFrequencyMs int    `json:"event_frequency_ms"`
}

// ClientConfig is the configuration for the daemon
var ClientConfig Config

// ReadConfig reads the configuration from the file
func ReadConfig(path string) error {
	// Read the configuration from the file
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Unmarshal the configuration
	err = json.Unmarshal(data, &ClientConfig)
	if err != nil {
		return err
	}

	// Check if the configuration is valid
	if ClientConfig.NATSURL == "" || ClientConfig.Subject == "" {
		return errors.New("invalid configuration")
	}

	// Check if the event frequency is valid
	if ClientConfig.EventFrequencyMs < 1000 {
		ClientConfig.EventFrequencyMs = 1000
		log.Println("Event frequency is invalid, setting to 1000 ms")
	}

	return nil
}
