package config

import (
	"log"
	"os"
	"testing"
)

// Mock configuration file data for testing
const validConfig = `{
	"debug": true,
	"nats_url": "nats://localhost:4222",
	"subject": "test.subject",
	"event_frequency_ms": 100
}`

const invalidConfigNoNATSURL = `{
	"debug": true,
	"subject": "test.subject",
	"event_frequency_ms": 100
}`

const invalidConfigNoSubject = `{
	"debug": true,
	"nats_url": "nats://localhost:4222",
	"event_frequency_ms": 100
}`

const invalidEventFrequencyConfig = `{
	"debug": true,
	"nats_url": "nats://localhost:4222",
	"subject": "test.subject",
	"event_frequency_ms": 0
}`

func writeTempFile(t *testing.T, content string) string {
	tempFile, err := os.CreateTemp("", "config_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	if _, err := tempFile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tempFile.Close()
	return tempFile.Name()
}

func TestReadConfig(t *testing.T) {
	log.SetOutput(os.Stdout) // For testing log output if needed
	t.Run("ValidConfig", func(t *testing.T) {
		filePath := writeTempFile(t, validConfig)
		defer os.Remove(filePath)

		if err := ReadConfig(filePath); err != nil {
			t.Errorf("Unexpected error reading valid config: %v", err)
		}

		if DeamonConfig.NATSURL != "nats://localhost:4222" || DeamonConfig.Subject != "test.subject" || DeamonConfig.EventFrequencyMs != 100 {
			t.Errorf("Config values not loaded properly: %+v", DeamonConfig)
		}
	})

	t.Run("MissingNATSURL", func(t *testing.T) {
		filePath := writeTempFile(t, invalidConfigNoNATSURL)
		defer os.Remove(filePath)

		err := ReadConfig(filePath)
		if err == nil || err.Error() != "invalid configuration" {
			t.Errorf("Expected error for missing NATSURL, got: %v", err)
		}
	})

	t.Run("MissingSubject", func(t *testing.T) {
		filePath := writeTempFile(t, invalidConfigNoSubject)
		defer os.Remove(filePath)

		err := ReadConfig(filePath)
		if err == nil || err.Error() != "invalid configuration" {
			t.Errorf("Expected error for missing Subject, got: %v", err)
		}
	})

	// t.Run("InvalidEventFrequency", func(t *testing.T) {
	// 	filePath := writeTempFile(t, invalidEventFrequencyConfig)
	// 	defer os.Remove(filePath)

	// 	if err := ReadConfig(filePath); err != nil {
	// 		t.Errorf("Unexpected error for invalid event frequency: %v", err)
	// 	}
	// 	if DeamonConfig.EventFrequencyMs != 1 {
	// 		t.Errorf("Expected event frequency to default to 1, got: %d", DeamonConfig.EventFrequencyMs)
	// 	}
	// })

	// t.Run("FileNotFound", func(t *testing.T) {
	// 	err := ReadConfig("non_existent_file.json")
	// 	if err == nil || !errors.Is(err, os.ErrNotExist) {
	// 		t.Errorf("Expected file not found error, got: %v", err)
	// 	}
	// })
}
