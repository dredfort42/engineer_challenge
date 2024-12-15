package config

import (
	"os"
	"testing"
)

// Mock configuration file data for testing
const validConfig = `{
	"debug": true,
	"nats_url": "nats://validConfig:4222",
	"subject": "validConfig.subject"
}`

const validConfigNoDebug = `{
	"nats_url": "nats://validConfigNoDebug:4222",
	"subject": "validConfigNoDebug.subject"
}`

const invalidConfigNoNATSURL = `{
	"debug": true,
	"subject": "invalidConfigNoNATSURL.subject"
}`

const invalidConfigNoSubject = `{
	"debug": true,
	"nats_url": "nats://invalidConfigNoSubject:4222"
}`

const validConfigWithExtraData = `{
	"debug": false,
	"nats_url": "nats://invalidEventFrequencyConfig:4222",
	"subject": "invalidEventFrequencyConfig.subject",
	"event_frequency_ms": -42
}`

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name     string
		fileData string
		expected Config
		errMsg   string
	}{
		{
			name:     "ValidConfig",
			fileData: validConfig,
			expected: Config{
				Debug:   true,
				NATSURL: "nats://validConfig:4222",
				Subject: "validConfig.subject",
			},
			errMsg: "",
		},
		{
			name:     "ValidConfigNoDebug",
			fileData: validConfigNoDebug,
			expected: Config{
				Debug:   false,
				NATSURL: "nats://validConfigNoDebug:4222",
				Subject: "validConfigNoDebug.subject"},
			errMsg: "",
		},
		{
			name:     "MissingNATSURL",
			fileData: invalidConfigNoNATSURL,
			expected: Config{},
			errMsg:   "invalid configuration",
		},
		{
			name:     "MissingSubject",
			fileData: invalidConfigNoSubject,
			expected: Config{},
			errMsg:   "invalid configuration",
		},
		{
			name:     "ExtraData",
			fileData: validConfigWithExtraData,
			expected: Config{
				Debug:   false,
				NATSURL: "nats://invalidEventFrequencyConfig:4222",
				Subject: "invalidEventFrequencyConfig.subject"},
			errMsg: "",
		},
		{
			name:     "FileNotFound",
			fileData: "",
			expected: Config{},
			errMsg:   "open non_existent_file.json: no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fileData != "" {
				tempFile, err := os.CreateTemp("./", "config_"+tt.name+".json")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				defer os.Remove(tempFile.Name())

				if _, err := tempFile.Write([]byte(tt.fileData)); err != nil {
					t.Fatalf("Failed to write temp file: %v", err)
				}

				tempFile.Close()

				WriterConfig = Config{}

				err = ReadConfig(tempFile.Name())
				if err != nil {
					if err.Error() != tt.errMsg {
						t.Errorf("Expected error %q, got: %v", tt.errMsg, err)
					}
				} else if WriterConfig != tt.expected {
					t.Errorf("Expected config %+v, got: %+v", tt.expected, WriterConfig)
				}
			} else {
				err := ReadConfig("non_existent_file.json")
				if err == nil || err.Error() != tt.errMsg {
					t.Errorf("Expected error %q, got: %v", tt.errMsg, err)
				}
			}
		})
	}
}
