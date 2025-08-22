package util

import (
	"os"
)

// GetUsageTelemetryPublisherServiceAddr returns the address of the usage-telemetry-publisher service
func GetUsageTelemetryPublisherServiceAddr() string {
	value := os.Getenv("USAGE_TELEMETRY_PUBLISHER_ADDRESS")
	if value != "" {
		return value
	}
	return "http://localhost:8080"
}

// GetenvDefault returns the value of the environment variable env if it exists, otherwise it returns the default value
func GetenvDefault(env, d string) string {
	v := os.Getenv(env)
	if v == "" {
		return d
	}

	return v
}
