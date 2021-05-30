package infra

import (
	"context"
	"net"

	"github.com/maitesin/numbers/internals/app"
	"golang.org/x/sync/semaphore"
)

//go:generate moq -out zmock_accepter_test.go -pkg infra_test . Accepter

type Accepter interface {
	Accept() (net.Conn, error)
}

type ClientManager struct {
	sem      *semaphore.Weighted
	accepter Accepter
}

func NewClientManager(numberOfClients int, accepter Accepter) *ClientManager {
	return &ClientManager{
		sem:      semaphore.NewWeighted(int64(numberOfClients)),
		accepter: accepter,
	}
}

func (cm *ClientManager) Start(ctx context.Context, cancel context.CancelFunc, reporter *app.Reporter) {
	defer func() {
		cancel()
	}()

	done := make(chan struct{})

	go func() {
		<-done
		cm.sem.Release(1)
	}()

	for {
		conn, err := cm.accepter.Accept()
		if err != nil {
			return
		}

		go app.ClientHandler(ctx, cancel, conn, reporter, done)
		err = cm.sem.Acquire(ctx, 1)
		if err != nil {
			return
		}
	}
}
