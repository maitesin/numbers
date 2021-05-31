package config_test

import (
	"os"
	"testing"

	"github.com/maitesin/numbers/config"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	// unset environment variables
	variables := []string{
		"NUMBERS_HOST",
		"NUMBERS_PORT",
		"NUMBERS_CONCURRENT_CLIENTS",
		"NUMBERS_TIME_BETWEEN_REPORTS",
	}
	for _, variable := range variables {
		err := os.Unsetenv(variable)
		require.NoError(t, err)
	}

	cfg, err := config.New()
	require.NoError(t, err)

	require.Equal(t, "127.0.0.1", cfg.Host)
	require.Equal(t, "4000", cfg.Port)
	require.Equal(t, 5, cfg.ConcurrentClients)
	require.Equal(t, 10, cfg.TimeBetweenReportsInSeconds)

	// set concurrent clients to not a number
	err = os.Setenv("NUMBERS_CONCURRENT_CLIENTS", "nine")
	require.NoError(t, err)

	cfg, err = config.New()
	require.NotNil(t, err)

	// set concurrent clients to 0
	err = os.Setenv("NUMBERS_CONCURRENT_CLIENTS", "0")
	require.NoError(t, err)

	cfg, err = config.New()
	require.NotNil(t, err)

	err = os.Unsetenv("NUMBERS_CONCURRENT_CLIENTS")
	require.NoError(t, err)

	// set time between reports to not a number
	err = os.Setenv("NUMBERS_TIME_BETWEEN_REPORTS", "twelve")
	require.NoError(t, err)

	cfg, err = config.New()
	require.NotNil(t, err)

	// set time between reports to 0
	err = os.Setenv("NUMBERS_TIME_BETWEEN_REPORTS", "0")
	require.NoError(t, err)

	cfg, err = config.New()
	require.NotNil(t, err)

	err = os.Unsetenv("NUMBERS_TIME_BETWEEN_REPORTS")
	require.NoError(t, err)

	// check that all the environment variables are being used correctly
	namesAndValues := [][2]string{
		{
			"NUMBERS_HOST", "10.10.10.10",
		},
		{
			"NUMBERS_PORT", "4444",
		},
		{
			"NUMBERS_CONCURRENT_CLIENTS", "3",
		},
		{
			"NUMBERS_TIME_BETWEEN_REPORTS", "25",
		},
	}

	for _, nameAndValue := range namesAndValues {
		err = os.Setenv(nameAndValue[0], nameAndValue[1])
		require.NoError(t, err)
	}

	cfg, err = config.New()
	require.NoError(t, err)

	require.Equal(t, "10.10.10.10", cfg.Host)
	require.Equal(t, "4444", cfg.Port)
	require.Equal(t, 3, cfg.ConcurrentClients)
	require.Equal(t, 25, cfg.TimeBetweenReportsInSeconds)
}
