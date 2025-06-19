package message

import (
	"context"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type Msg struct {
	Payload string
	Topic   string
}

type PubSub struct {
	rdb        *redis.Client
	logger     watermill.LoggerAdapter
	router     *message.Router
	publishers map[string]*redisstream.Publisher
}

func NewPubSub() *PubSub {
	logger := watermill.NewSlogLogger(nil)
	router := message.NewDefaultRouter(logger)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	return &PubSub{
		rdb:        rdb,
		logger:     logger,
		router:     router,
		publishers: make(map[string]*redisstream.Publisher),
	}
}

func (ps *PubSub) Init(ctx context.Context) error {
	return ps.router.Run(ctx)

	// go func() {
	// 	defer log.Println("Done")
	// 	log.Println("Dude")
	// 	if err := ps.router.Run(context.Background()); err != nil {
	// 		panic(err.Error())
	// 	}

	// }()
}

func (ps *PubSub) WaitUntilReady() {
	<-ps.router.Running()
	// for !ps.router.IsRunning() {
	// 	time.Sleep(10 * time.Millisecond)
	// }
}

func (ps *PubSub) Subscribe(topic, group string, callback func(context.Context, string) error) {
	sub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        ps.rdb,
		ConsumerGroup: group,
	}, ps.logger)
	if err != nil {
		panic("Could not subscribe! " + err.Error())
	}

	ps.router.AddNoPublisherHandler(
		fmt.Sprintf("%s-%s", topic, group),
		topic,
		sub,
		func(msg *message.Message) error {
			return callback(msg.Context(), string(msg.Payload))
		},
	)

	if err := ps.router.RunHandlers(context.Background()); err != nil {
		panic("We could not run handlers " + err.Error())
	}
	fmt.Println("Run handers done!")

	// go processMessages(topic, sub, callback)
}

// func processMessages(topic string, sub message.Subscriber, action func(context.Context, string) error) {
// 	messages, err := sub.Subscribe(context.Background(), topic)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for msg := range messages {
// 		orderID := string(msg.Payload)
// 		err := action(msg.Context(), orderID)
// 		if err != nil {
// 			msg.Nack()
// 		} else {
// 			msg.Ack()
// 		}
// 	}
// }

func (ps *PubSub) SendMessages(messages ...Msg) {
	for _, msg := range messages {
		var err error

		if ps.publishers[msg.Topic] == nil {
			if ps.publishers[msg.Topic], err = redisstream.NewPublisher(redisstream.PublisherConfig{
				Client: ps.rdb,
			}, ps.logger); err != nil {
				panic("Dude, wer could not create the publisher... " + err.Error())
			}
		}
		if err = ps.publishers[msg.Topic].Publish(msg.Topic, &message.Message{
			UUID:    watermill.NewUUID(),
			Payload: []byte(msg.Payload),
		}); err != nil {
			panic("Dude, we could not publish a message... " + err.Error())
		}
	}
}
