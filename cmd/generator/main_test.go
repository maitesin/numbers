package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomString(t *testing.T) {
	random := RandomString(9)
	require.Len(t, random, 9)
	for _, r := range random {
		require.True(t, r <= '9')
		require.True(t, r >= '0')
	}
}
