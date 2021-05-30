package app

import (
	"bufio"
	"context"
	"io"

	"github.com/maitesin/numbers/internals/domain"
)

const stopServerCmd = "terminate"

func ClientHandler(ctx context.Context, cancelCtx context.CancelFunc, readCloser io.ReadCloser, reporter *Reporter, done chan<- struct{}) {
	defer func() {
		done <- struct{}{}
		readCloser.Close()
	}()

	bufferedReader := bufio.NewReader(readCloser)
	line, err := bufferedReader.ReadString('\n')
	if err != nil {
		return
	}

	line = line[:len(line)-1] // Remove trailing end line character

	if line == stopServerCmd {
		cancelCtx()
		return
	}

	number, err := domain.NewNumber(line)
	if err != nil {
		return
	}

	err = reporter.Record(number)
	if err != nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			line, err := bufferedReader.ReadString('\n')
			if err != nil {
				return
			}

			line = line[:len(line)-1] // Remove trailing end line character

			number, err := domain.NewNumber(line)
			if err != nil {
				return
			}

			err = reporter.Record(number)
			if err != nil {
				return
			}
		}
	}
}
