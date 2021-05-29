package app

import (
	"fmt"
	"io"
	"sync"

	"github.com/maitesin/numbers/internals/domain"
)

type Reporter struct {
	written          map[domain.Number]struct{}
	uniqueCounter    int
	duplicateCounter int
	totalCounter     int
	statsWriter      io.Writer
	uniqueWriter     io.Writer
	m                sync.Mutex
}

func NewReporter(uniqueWrite, statsWriter io.Writer) *Reporter {
	return &Reporter{
		written:      map[domain.Number]struct{}{},
		uniqueWriter: uniqueWrite,
		statsWriter:  statsWriter,
	}
}

func (r *Reporter) Record(number domain.Number) error {
	r.m.Lock()
	defer r.m.Unlock()

	_, ok := r.written[number]
	if !ok {
		r.written[number] = struct{}{}
		r.uniqueCounter++
		_, err := fmt.Fprintln(r.uniqueWriter, number.Value)
		if err != nil {
			return err
		}
	} else {
		r.duplicateCounter++
	}
	r.totalCounter++

	return nil
}

func (r *Reporter) Report() error {
	r.m.Lock()
	defer r.m.Unlock()

	_, err := fmt.Fprintf(
		r.statsWriter,
		"Received %d unique numbers, %d duplicates. Unique total: %d\n",
		r.uniqueCounter,
		r.duplicateCounter,
		r.totalCounter,
	)
	if err != nil {
		return err
	}

	r.uniqueCounter = 0
	r.duplicateCounter = 0
	return nil
}
