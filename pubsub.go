package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/pubsub"
)

func startPubSubEmulator(ctx context.Context, port int) error {
	cmd := exec.CommandContext(ctx, "gcloud", "beta", "emulators", "pubsub", "start", fmt.Sprintf("--host-port=0.0.0.0:%d", port))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func parseTopicAndSubscriptions(configPath string) (map[string][]string, error) {
	config := map[string][]string{}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return config, nil
}

// https://github.com/prep/pubsubc/blob/acd31b169239c1a8f0ed4d2a45d8a9e6f813a4a0/main.go
func registerTopicAndSubscriptions(ctx context.Context, projectID string, pubsAndSubs map[string][]string) error {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("Unable to create client to project %q: %s", projectID, err)
	}
	defer client.Close()

	log.Printf("Project %q", projectID)

	for topicID, subscriptions := range pubsAndSubs {
		log.Printf("  Topic %q", topicID)
		topic, err := client.CreateTopic(ctx, topicID)
		if err != nil {
			return fmt.Errorf("creating topic %q for project %q: %w", topicID, projectID, err)
		}

		for _, subscriptionID := range subscriptions {
			log.Printf("    Subscription %q", subscriptionID)
			_, err = client.CreateSubscription(ctx, subscriptionID, pubsub.SubscriptionConfig{Topic: topic})
			if err != nil {
				return fmt.Errorf("creating subscription %q on topic %q for project %q: %w", subscriptionID, topicID, projectID, err)
			}
		}
	}
	return nil
}
