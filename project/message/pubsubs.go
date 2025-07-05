package message

import (
	"context"
	"fmt"
	"os"
	"tickets/message/mid"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type Msg struct {
	Message *message.Message
	Topic   string
}

type PubSub struct {
	rdb           *redis.Client
	logger        watermill.LoggerAdapter
	router        *message.Router
	publishers    map[string]*redisstream.Publisher
	globalContext context.Context
}

func NewPubSub() *PubSub {
	logger := watermill.NewSlogLogger(nil)
	router := message.NewDefaultRouter(logger)

	router.AddMiddleware(
		mid.EnsureCorrelationId,
	)
	router.AddMiddleware(
		mid.LoggingMiddleware{
			Logger: logger,
		}.Middleware,
	)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	return &PubSub{
		rdb:           rdb,
		logger:        logger,
		router:        router,
		publishers:    make(map[string]*redisstream.Publisher),
		globalContext: context.Background(),
	}
}

func (ps *PubSub) Init(ctx context.Context) error {
	ps.globalContext = ctx
	return ps.router.Run(ctx)
}

func (ps *PubSub) WaitUntilReady() {
	<-ps.router.Running()
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

	if err := ps.router.RunHandlers(ps.globalContext); err != nil {
		panic("We could not run handlers " + err.Error())
	}
	fmt.Println("Run handers done!")

}

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
		if err = ps.publishers[msg.Topic].Publish(msg.Topic, msg.Message); err != nil {
			panic("Dude, we could not publish a message... " + err.Error())
		}
	}
}
