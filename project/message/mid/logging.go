package mid

import (
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type LoggingMiddleware struct {
	Logger watermill.LoggerAdapter
}

func (md LoggingMiddleware) Middleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		logger := log.FromContext(msg.Context()).With("message_id", msg.UUID)
		logger.Info("Handling a message")
		return next(msg)
	}
}
