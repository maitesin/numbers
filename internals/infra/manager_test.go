package infra_test

import (
	"bytes"
	"context"
	"errors"
	"net"
	"testing"

	"github.com/maitesin/numbers/internals/app"
	"github.com/maitesin/numbers/internals/infra"
)

func TestClientManager_Start(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	accepter := &AccepterMock{
		AcceptFunc: func() (net.Conn, error) {
			return nil, errors.New("something went wrong")
		},
	}

	cm := infra.NewClientManager(5, accepter)
	cm.Start(ctx, cancel, app.NewReporter(&bytes.Buffer{}, &bytes.Buffer{}))
	<-ctx.Done()
}
