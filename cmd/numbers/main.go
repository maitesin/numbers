package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/maitesin/numbers/internals/app"
	"github.com/maitesin/numbers/internals/infra"
)

const (
	exitStatusFailedOpeningNumbersFile int = iota + 1
	exitStatusFailedListen
)

const (
	maxNumberOfClients          = 5
	timeBetweenReportsInSeconds = 10
)

func main() {
	numbersFile, err := os.OpenFile("numbers.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitStatusFailedOpeningNumbersFile)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:4000")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitStatusFailedListen)
	}

	reporter := app.NewReporter(numbersFile, os.Stdout)
	clientManager := infra.NewClientManager(maxNumberOfClients, ln)

	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(timeBetweenReportsInSeconds * time.Second)
	finish := make(chan struct{})

	go app.CallReportAtEveryTick(reporter, ticker, finish)
	go clientManager.Start(ctx, cancel, reporter)
	<-ctx.Done()
	finish <- struct{}{}
}
