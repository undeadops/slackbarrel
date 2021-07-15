package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"barrel/pkg/agent"
)

var (
	serverUrl  string
	configFile string
	interval   int
	app        string
)

func main() {
	flag.StringVar(&configFile, "config", "./agent.conf", "Agent Config File")
	flag.StringVar(&serverUrl, "server", "http://localhost:5000", "Config Management Server to connect too")
	flag.IntVar(&interval, "interval", 25, "Checkin interval with Management Server")
	flag.Parse()

	// Setup Connection timeouts
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Capture Sig Terms, send to sub goroutine
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGQUIT)
	//signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	agent := agent.NewAgent(configFile, serverUrl, interval)

	// Checkin on a timer, interruptable between runs with the Interrupt channel
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-interrupt:
			os.Exit(0)
		case <-ticker.C:
			agent.CheckIn(interrupt)
		}
	}
}
