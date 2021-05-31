package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/maitesin/numbers/config"
	"github.com/maitesin/numbers/internals/app"
	"github.com/maitesin/numbers/internals/infra"
)

const (
	exitStatusFailedOpeningNumbersFile int = iota + 1
	exitStatusFailedConfiguration
	exitStatusFailedListen
)

func main() {
	numbersFile, err := os.OpenFile("numbers.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitStatusFailedOpeningNumbersFile)
	}

	cfg, err := config.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitStatusFailedConfiguration)
	}

	ln, err := net.Listen("tcp", strings.Join([]string{cfg.Host, cfg.Port}, ":"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitStatusFailedListen)
	}

	reporter := app.NewReporter(numbersFile, os.Stdout)
	clientManager := infra.NewClientManager(cfg.ConcurrentClients, ln)

	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(time.Duration(cfg.TimeBetweenReportsInSeconds) * time.Second)
	finish := make(chan struct{})

	go app.CallReportAtEveryTick(reporter, ticker, finish)
	go clientManager.Start(ctx, cancel, reporter)
	<-ctx.Done()
	finish <- struct{}{}
}
