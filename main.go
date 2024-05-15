package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

const (
	pubsubPort = 8681
	readyPort  = 8682

	topicsAndSubscriptionsConfigPath = "topics_and_subscriptions.json"
)

var (
	pubsubProjects = []string{"development", "test"}
)

func run() error {
	// 1. Start pubsub emulator
	// 2. Wait for port 8681 to be ready
	// 3. Register topics and subscriptions
	// 4. Listen and discard on 8682
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	defer func(wg *sync.WaitGroup) {
		cancel()
		log.Printf("waiting for goroutines to exit")
		wg.Wait()
		log.Printf("all goroutines done")
	}(&wg)

	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		if err := startPubSubEmulator(ctx, pubsubPort); err != nil {
			log.Printf("failed to start pubsub emulator: %v", err)
		}
	}(ctx, &wg)

	if err := waitFor(ctx, pubsubPort, 15*time.Second); err != nil {
		log.Printf("failed to start pubsub emulator: %v", err)
		return err
	}

	pubsAndSubs, err := parseTopicAndSubscriptions(topicsAndSubscriptionsConfigPath)
	if err != nil {
		log.Printf("failed to parse topics and subscriptions from '%s': %v", topicsAndSubscriptionsConfigPath, err)
		return err
	}

	for _, project := range pubsubProjects {
		if err := registerTopicAndSubscriptions(ctx, project, pubsAndSubs); err != nil {
			log.Printf("failed to register topics and subscriptions: %v", err)
			return err
		}
	}

	serv := &http.Server{
		Addr: fmt.Sprintf(":%d", readyPort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	}

	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		if err := serv.ListenAndServe(); err != nil {
			log.Printf("failed to listen and serve: %v", err)
		}
	}(ctx, &wg)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case <-interrupt:
		log.Printf("interrupted")
		serv.Shutdown(context.Background())
		cancel()
		return nil
	case <-ctx.Done():
		serv.Shutdown(context.Background())
		log.Printf("context done")
	}
	return nil
}
