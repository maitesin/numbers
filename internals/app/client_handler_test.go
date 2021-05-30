package app_test

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/maitesin/numbers/internals/app"
	"github.com/stretchr/testify/require"
)

type ctxGenerator func() (context.Context, context.CancelFunc)

func validCtxGenerator() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func canceledCtxGenerator() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx, cancel
}

type secondCallFailingWriter struct {
	calls int
}

func (fw *secondCallFailingWriter) Write([]byte) (int, error) {
	if fw.calls == 0 {
		fw.calls++
		return 0, nil
	}
	return 0, errFromFailingWriter
}

func TestClientHandler(t *testing.T) {
	tests := []struct {
		name           string
		ctxGenerator   ctxGenerator
		reader         io.ReadCloser
		reporter       *app.Reporter
		expectedCtxErr error
	}{
		{
			name: `Given a non canceled context, a reader containing the exit command, and a working reporter,
                   when the client handler is executed,
                   then the context cancel function is called`,
			ctxGenerator:   validCtxGenerator,
			reader:         ioutil.NopCloser(strings.NewReader("terminate\n")),
			reporter:       app.NewReporter(&bytes.Buffer{}, &bytes.Buffer{}),
			expectedCtxErr: context.Canceled,
		},
		{
			name: `Given a non canceled context, an empty reader, and a working reporter,
                   when the client handler is executed,
                   then nothing is reported to the reporter`,
			ctxGenerator: validCtxGenerator,
			reader:       ioutil.NopCloser(strings.NewReader("")),
			reporter:     app.NewReporter(&bytes.Buffer{}, &bytes.Buffer{}),
		},
		{
			name: `Given a non canceled context, a reader with an invalid number, and a working reporter,
                   when the client handler is executed,
                   then nothing is reported to the reporter`,
			ctxGenerator: validCtxGenerator,
			reader:       ioutil.NopCloser(strings.NewReader("wololo\n")),
			reporter:     app.NewReporter(&bytes.Buffer{}, &bytes.Buffer{}),
		},
		{
			name: `Given a non canceled context, a reader with a valid number, and a working reporter,
                   when the client handler is executed,
                   then the number is reported to the reporter`,
			ctxGenerator: validCtxGenerator,
			reader:       ioutil.NopCloser(strings.NewReader("123456789\n")),
			reporter:     app.NewReporter(&bytes.Buffer{}, &bytes.Buffer{}),
		},
		{
			name: `Given a non canceled context, a reader with multiple valid numbers, and a working reporter,
                   when the client handler is executed,
                   then the numbers are reported to the reporter`,
			ctxGenerator: validCtxGenerator,
			reader:       ioutil.NopCloser(strings.NewReader("123456789\n987654321\n")),
			reporter:     app.NewReporter(&bytes.Buffer{}, &bytes.Buffer{}),
		},
		{
			name: `Given a non canceled context, a reader with a valid number and an invalid one, and a working reporter,
                   when the client handler is executed,
                   then the number is reported to the reporter`,
			ctxGenerator: validCtxGenerator,
			reader:       ioutil.NopCloser(strings.NewReader("123456789\nwololo\n")),
			reporter:     app.NewReporter(&bytes.Buffer{}, &bytes.Buffer{}),
		},
		{
			name: `Given a non canceled context, a reader with a valid number, and a failing in the first call reporter,
                   when the client handler is executed,
                   then nothing is reported by the reporter`,
			ctxGenerator: validCtxGenerator,
			reader:       ioutil.NopCloser(strings.NewReader("123456789\n")),
			reporter:     app.NewReporter(&failingWriter{}, &bytes.Buffer{}),
		},
		{
			name: `Given a non canceled context, a reader with multiple valid numbers, and a failing in the second call reporter,
                   when the client handler is executed,
                   then the first number is successfully reported`,
			ctxGenerator: validCtxGenerator,
			reader:       ioutil.NopCloser(strings.NewReader("123456789\n987654321\n")),
			reporter:     app.NewReporter(&secondCallFailingWriter{}, &bytes.Buffer{}),
		},
		{
			name: `Given a canceled context, a reader with a valid number, and a working reporter,
                   when the client handler is executed,
                   then the number is successfully reported and the context is returned`,
			ctxGenerator:   canceledCtxGenerator,
			reader:         ioutil.NopCloser(strings.NewReader("123456789\n")),
			reporter:       app.NewReporter(&bytes.Buffer{}, &bytes.Buffer{}),
			expectedCtxErr: context.Canceled,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := tt.ctxGenerator()
			done := make(chan struct{})
			go app.ClientHandler(ctx, cancel, tt.reader, tt.reporter, done)
			<-done
			if tt.expectedCtxErr != nil {
				require.ErrorIs(t, ctx.Err(), tt.expectedCtxErr)
			} else {
				require.NoError(t, ctx.Err())
			}
		})
	}
}
