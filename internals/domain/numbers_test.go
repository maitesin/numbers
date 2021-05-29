package domain_test

import (
	"testing"

	"github.com/maitesin/numbers/internals/domain"
	"github.com/stretchr/testify/require"
)

type lineMutator func(string) string

func validLine() string { return "012345678" }

func noopLineMutator(line string) string { return line }

func TestNewNumber(t *testing.T) {
	tests := []struct {
		name           string
		lineMutator    lineMutator
		expectedNumber domain.Number
		expectedErr    error
	}{
		{
			name: `Given a line containing 9 characters and all of them are numbers,
                   when the new number function is called with that line,
                   then a valid number is returned`,
			lineMutator:    noopLineMutator,
			expectedNumber: domain.Number{Value: "012345678"},
		},
		{
			name: `Given a line containing 10 characters and all of them are numbers,
                   when the new number function is called with that line,
                   then an invalid number error is returned`,
			lineMutator: func(string) string { return "0123456789" },
			expectedErr: domain.ErrInvalidNumber,
		},
		{
			name: `Given a line containing 9 characters and not all of them are numbers,
                   when the new number function is called with that line,
                   then an invalid number error is returned`,
			lineMutator: func(string) string { return "0123A5678" },
			expectedErr: domain.ErrInvalidNumber,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := domain.NewNumber(tt.lineMutator(validLine()))
			if tt.expectedErr != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedNumber, got)
			}
		})
	}
}
