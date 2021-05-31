package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	defaultConcurrentClients           = "5"
	defaultTimeBetweenReportsInSeconds = "10"
	defaultHost                        = "127.0.0.1"
	defaultPort                        = "4000"
)

type Config struct {
	ConcurrentClients           int
	TimeBetweenReportsInSeconds int
	Host                        string
	Port                        string
}

func New() (Config, error) {
	clients, err := strconv.Atoi(getEnvOrDefault("NUMBERS_CONCURRENT_CLIENTS", defaultConcurrentClients))
	if err != nil {
		return Config{}, err
	}

	if clients < 1 {
		return Config{}, NewInvalidConfigError("concurrent_clients", fmt.Sprint(clients))
	}

	timeBetweenReports, err := strconv.Atoi(getEnvOrDefault("NUMBERS_TIME_BETWEEN_REPORTS", defaultTimeBetweenReportsInSeconds))
	if err != nil {
		return Config{}, err
	}

	if timeBetweenReports < 1 {
		return Config{}, NewInvalidConfigError("time_between_reports", fmt.Sprint(timeBetweenReports))
	}

	return Config{
		ConcurrentClients:           clients,
		TimeBetweenReportsInSeconds: timeBetweenReports,
		Host:                        getEnvOrDefault("NUMBERS_HOST", defaultHost),
		Port:                        getEnvOrDefault("NUMBERS_PORT", defaultPort),
	}, nil
}

func getEnvOrDefault(name, defaultValue string) string {
	value := os.Getenv(name)
	if value != "" {
		return value
	}

	return defaultValue
}
