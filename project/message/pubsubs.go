package message

import (
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
	publishers map[string]*redisstream.Publisher
}

func NewPubSub() *PubSub {
	logger := watermill.NewSlogLogger(nil)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	return &PubSub{
		rdb:        rdb,
		logger:     logger,
		publishers: make(map[string]*redisstream.Publisher),
	}
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
		if err = ps.publishers[msg.Topic].Publish(msg.Topic, &message.Message{
			UUID:    watermill.NewUUID(),
			Payload: []byte(msg.Payload),
		}); err != nil {
			panic("Dude, we could not publish a message... " + err.Error())
		}
	}
}
