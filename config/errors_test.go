package config_test

import (
	"testing"

	"github.com/maitesin/numbers/config"
	"github.com/stretchr/testify/require"
)

func TestInvalidConfigError(t *testing.T) {
	t.Parallel()

	err := config.NewInvalidConfigError("-wololo", "42")
	require.Equal(t, `invalid configuration in field "-wololo" with value "42"`, err.Error())
}
