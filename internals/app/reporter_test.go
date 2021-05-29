package app_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/maitesin/numbers/internals/app"
	"github.com/maitesin/numbers/internals/domain"
	"github.com/stretchr/testify/require"
)

type numberMutator func(domain.Number) domain.Number

func noopNumberMutator(number domain.Number) domain.Number { return number }

func validNumber(t *testing.T) domain.Number {
	t.Helper()

	number, err := domain.NewNumber("123456789")
	require.NoError(t, err)

	return number
}

var errFromFailingWriter = errors.New("something went wrong while writing")

type failingWriter struct{}

func (wm failingWriter) Write([]byte) (int, error) {
	return 0, errFromFailingWriter
}

func TestReporter_Record(t *testing.T) {
	fixtureNumber, err := domain.NewNumber("987654321")
	require.NoError(t, err)

	tests := []struct {
		name           string
		numberMutator  numberMutator
		expectedOutput string
		expectedErr    error
	}{
		{
			name: `Given a number that has not been recorded,
                   when the record method is called,
                   then no error is returned and the number is written in the unique writer`,
			numberMutator:  noopNumberMutator,
			expectedOutput: "123456789\n",
		},
		{
			name: `Given a number that has been recorded,
                   when the record method is called,
                   then no error is returned and the number is not written in the unique writer`,
			numberMutator: func(domain.Number) domain.Number {
				return fixtureNumber
			},
			expectedOutput: "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uniqueWriter := &bytes.Buffer{}
			statsWriter := &bytes.Buffer{}

			r := app.NewReporter(uniqueWriter, statsWriter)

			err := r.Record(fixtureNumber)
			require.NoError(t, err)

			uniqueWriter.Reset()

			err = r.Record(tt.numberMutator(validNumber(t)))
			if tt.expectedErr != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedOutput, uniqueWriter.String())
				require.Equal(t, "", statsWriter.String())
			}
		})
	}

	t.Run(`Given a failing unique writer,
		when the record method is called,
		then an error is returned`,
		func(t *testing.T) {
			t.Parallel()

			r := app.NewReporter(&failingWriter{}, &bytes.Buffer{})
			err := r.Record(fixtureNumber)
			require.ErrorIs(t, err, errFromFailingWriter)
		})
}

func TestReporter_Report(t *testing.T) {
	t.Run(`Given a working stats writer,
                 when the report method is called consecutive times the unique and duplicate counters are reset,
                 then no error is returned`,
		func(t *testing.T) {
			t.Parallel()

			statsWriter := &bytes.Buffer{}

			r := app.NewReporter(&bytes.Buffer{}, statsWriter)

			err := r.Record(validNumber(t))
			require.NoError(t, err)
			err = r.Record(validNumber(t))
			require.NoError(t, err)

			err = r.Report()
			require.NoError(t, err)
			require.Equal(t, "Received 1 unique numbers, 1 duplicates. Unique total: 2\n", statsWriter.String())

			statsWriter.Reset()

			err = r.Report()
			require.NoError(t, err)
			require.Equal(t, "Received 0 unique numbers, 0 duplicates. Unique total: 2\n", statsWriter.String())
		})

	t.Run(`Given a failing stats writer,
		when the report method is called,
		then an error is returned`,
		func(t *testing.T) {
			t.Parallel()

			r := app.NewReporter(&bytes.Buffer{}, &failingWriter{})
			err := r.Report()
			require.ErrorIs(t, err, errFromFailingWriter)
		})
}
