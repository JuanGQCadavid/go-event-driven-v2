package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type PaymentCompleted struct {
	PaymentID   string `json:"payment_id"`
	OrderID     string `json:"order_id"`
	CompletedAt string `json:"completed_at"`
}

type Confirmed struct {
	OrderID     string `json:"order_id"`
	ConfirmedAt string `json:"confirmed_at"`
}

func main() {
	logger := watermill.NewSlogLogger(nil)

	router := message.NewDefaultRouter(logger)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	sub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	pub, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	router.AddHandler(
		"02-marshalling",
		"payment-completed",
		sub,
		"order-confirmed",
		pub,
		func(msg *message.Message) ([]*message.Message, error) {
			var (
				payload *PaymentCompleted = &PaymentCompleted{}
			)

			if err := json.Unmarshal(msg.Payload, payload); err != nil {
				fmt.Println("Error while unmarshaling ", err.Error())
				return nil, err
			}

			confirmed := &Confirmed{
				ConfirmedAt: payload.CompletedAt,
				OrderID:     payload.OrderID,
			}

			newPayload, err := json.Marshal(confirmed)

			if err != nil {
				fmt.Println("Error while marshaling ", err.Error())
				return nil, err
			}

			return []*message.Message{
				{
					UUID:    msg.UUID,
					Payload: newPayload,
				},
			}, nil

		},
	)

	err = router.Run(context.Background())
	if err != nil {
		panic(err)
	}
}
