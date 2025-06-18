package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/redis/go-redis/v9"
)

func main() {

	logger := watermill.NewSlogLogger(nil)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	subs, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)

	if err != nil {
		logger.Error("Error while creating the subs", err, nil)
	}

	messages, err := subs.Subscribe(context.Background(), "progress")

	if err != nil {
		logger.Error("Error while subscribing", err, nil)
		panic("PANIC!")
	}

	for msg := range messages {
		status := string(msg.Payload)
		logger.Info(fmt.Sprintf("Message ID: %s - %s", msg.UUID, status), nil)
		// logger.Info(msg.UUID+" "+status, nil)
		msg.Ack()
	}

}
