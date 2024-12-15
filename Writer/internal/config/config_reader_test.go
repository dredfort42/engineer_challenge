package config

import (
	"os"
	"testing"
)

// Mock configuration file data for testing
const validConfig = `{
	"debug": true,
	"nats_url": "nats://validConfig:4222",
	"subject": "validConfig.subject",
	"influxdb_url": "http://validConfig:8086",
	"influxdb_org": "validConfig",
	"influxdb_bucket": "validConfig",
	"influxdb_measurement": "validConfig",
	"path_to_influxdb_token": "./ValidConfig"
}`

const validConfigNoDebug = `{
	"nats_url": "nats://validConfigNoDebug:4222",
	"subject": "validConfigNoDebug.subject",
	"influxdb_url": "http://validConfigNoDebug:8086",
	"influxdb_org": "validConfigNoDebug",
	"influxdb_bucket": "validConfigNoDebug",
	"influxdb_measurement": "validConfigNoDebug",
	"path_to_influxdb_token": "./ValidConfigNoDebug"
}`

const invalidConfigNoNATSURL = `{
	"debug": true,
	"subject": "invalidConfigNoNATSURL.subject",
	"influxdb_url": "http://invalidConfigNoNATSURL:8086",
	"influxdb_org": "invalidConfigNoNATSURL",
	"influxdb_bucket": "invalidConfigNoNATSURL",
	"influxdb_measurement": "invalidConfigNoNATSURL",
	"path_to_influxdb_token": "./MissingNATSURL"
}`

const invalidConfigNoSubject = `{
	"debug": true,
	"nats_url": "nats://invalidConfigNoSubject:4222",
	"influxdb_url": "http://invalidConfigNoSubject:8086",
	"influxdb_org": "invalidConfigNoSubject",
	"influxdb_bucket": "invalidConfigNoSubject",
	"influxdb_measurement": "invalidConfigNoSubject",
	"path_to_influxdb_tokee": "./MissingSubject"
}`

const invalidConfigNoInfluxURL = `{
	"debug": true,
	"nats_url": "nats://invalidConfigNoInfluxURL:4222",
	"subject": "invalidConfigNoInfluxURL.subject",
	"influxdb_org": "invalidConfigNoInfluxURL",
	"influxdb_bucket": "invalidConfigNoInfluxURL",
	"influxdb_measurement": "invalidConfigNoInfluxURL",
	"path_to_influxdb_token": "./MissingInfluxURL"
}`

const invalidConfigNoInfluxOrg = `{
	"debug": true,
	"nats_url": "nats://invalidConfigNoInfluxOrg:4222",
	"subject": "invalidConfigNoInfluxOrg.subject",
	"influxdb_url": "http://invalidConfigNoInfluxOrg:8086",
	"influxdb_bucket": "invalidConfigNoInfluxOrg",
	"influxdb_measurement": "invalidConfigNoInfluxOrg",
	"path_to_influxdb_token": "./MissingInfluxOrg"
}`

const invalidConfigNoInfluxBucket = `{
	"debug": true,
	"nats_url": "nats://invalidConfigNoInfluxBucket:4222",
	"subject": "invalidConfigNoInfluxBucket.subject",
	"influxdb_url": "http://invalidConfigNoInfluxBucket:8086",
	"influxdb_org": "invalidConfigNoInfluxBucket",
	"influxdb_measurement": "invalidConfigNoInfluxBucket",
	"path_to_influxdb_token": "./MissingInfluxBucket"
}`

const invalidConfigNoInfluxMeasurement = `{	
	"debug": true,
	"nats_url": "nats://invalidConfigNoInfluxMeasurement:4222",
	"subject": "invalidConfigNoInfluxMeasurement.subject",
	"influxdb_url": "http://invalidConfigNoInfluxMeasurement:8086",
	"influxdb_org": "invalidConfigNoInfluxMeasurement",
	"influxdb_bucket": "invalidConfigNoInfluxMeasurement",
	"path_to_influxdb_token": "./MissingInfluxMeasurement"
}`

const invalidConfigNoPathToInfluxToken = `{	
	"debug": true,
	"nats_url": "nats://invalidConfigNoPathToInfluxToken:4222",
	"subject": "invalidConfigNoPathToInfluxToken.subject",
	"influxdb_url": "http://invalidConfigNoPathToInfluxToken:8086",
	"influxdb_org": "invalidConfigNoPathToInfluxToken",
	"influxdb_bucket": "invalidConfigNoPathToInfluxToken",
	"influxdb_measurement": "./MissingPathToInfluxToken"
}`

const invalidInfluxTokenFile = `{
	"debug": true,
	"nats_url": "nats://invalidInfluxTokenFile:4222",
	"subject": "invalidInfluxTokenFile.subject",
	"influxdb_url": "http://invalidInfluxTokenFile:8086",
	"influxdb_org": "invalidInfluxTokenFile",
	"influxdb_bucket": "invalidInfluxTokenFile",
	"influxdb_measurement": "invalidInfluxTokenFile",
	"path_to_influxdb_token": "./InvalidInfluxTokenFile"
}`

const validConfigWithExtraData = `{
	"debug": false,
	"nats_url": "nats://validConfigWithExtraData:4222",
	"subject": "validConfigWithExtraData.subject",
	"influxdb_url": "http://validConfigWithExtraData:8086",
	"influxdb_org": "validConfigWithExtraData",
	"influxdb_bucket": "validConfigWithExtraData",
	"influxdb_measurement": "validConfigWithExtraData",
	"path_to_influxdb_token": "./ExtraData",
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
				Debug:             true,
				NATSURL:           "nats://validConfig:4222",
				Subject:           "validConfig.subject",
				InfluxURL:         "http://validConfig:8086",
				InfluxOrg:         "validConfig",
				InfluxBucket:      "validConfig",
				InfluxMeasurement: "validConfig",
				PathToInfluxToken: "./ValidConfig",
			},
			errMsg: "",
		},
		{
			name:     "ValidConfigNoDebug",
			fileData: validConfigNoDebug,
			expected: Config{
				Debug:             false,
				NATSURL:           "nats://validConfigNoDebug:4222",
				Subject:           "validConfigNoDebug.subject",
				InfluxURL:         "http://validConfigNoDebug:8086",
				InfluxOrg:         "validConfigNoDebug",
				InfluxBucket:      "validConfigNoDebug",
				InfluxMeasurement: "validConfigNoDebug",
				PathToInfluxToken: "./ValidConfigNoDebug",
			},
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
			name:     "MissingInfluxURL",
			fileData: invalidConfigNoInfluxURL,
			expected: Config{},
			errMsg:   "invalid configuration",
		},
		{
			name:     "MissingInfluxOrg",
			fileData: invalidConfigNoInfluxOrg,
			expected: Config{},
			errMsg:   "invalid configuration",
		},
		{
			name:     "MissingInfluxBucket",
			fileData: invalidConfigNoInfluxBucket,
			expected: Config{},
			errMsg:   "invalid configuration",
		},
		{
			name:     "MissingInfluxMeasurement",
			fileData: invalidConfigNoInfluxMeasurement,
			expected: Config{},
			errMsg:   "invalid configuration",
		},
		{
			name:     "MissingPathToInfluxToken",
			fileData: invalidConfigNoPathToInfluxToken,
			expected: Config{},
			errMsg:   "invalid configuration",
		},
		{
			name:     "InvalidInfluxTokenFile",
			fileData: invalidInfluxTokenFile,
			expected: Config{},
			errMsg:   "influxdb token file does not exist",
		},
		{
			name:     "ExtraData",
			fileData: validConfigWithExtraData,
			expected: Config{
				Debug:             false,
				NATSURL:           "nats://validConfigWithExtraData:4222",
				Subject:           "validConfigWithExtraData.subject",
				InfluxURL:         "http://validConfigWithExtraData:8086",
				InfluxOrg:         "validConfigWithExtraData",
				InfluxBucket:      "validConfigWithExtraData",
				InfluxMeasurement: "validConfigWithExtraData",
				PathToInfluxToken: "./ExtraData",
			},
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
				if tt.name != "InvalidInfluxTokenFile" {
					tokenFile, err := os.Create(tt.name)
					if err != nil {
						t.Fatalf("Failed to create temp file: %v", err)
					}
					defer os.Remove(tokenFile.Name())
				}

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
