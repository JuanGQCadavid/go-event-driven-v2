package main

import (
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := watermill.NewSlogLogger(nil)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)

	if err != nil {
		logger.Error("Error while creating the publisher", err, nil)
	}

	messages := []string{"50", "100"}

	for _, msg := range messages {
		if err := publisher.Publish("progress", message.NewMessage(watermill.NewUUID(), []byte(msg))); err != nil {
			logger.Error("Error while publishing "+msg, err, nil)
		}
	}

}
