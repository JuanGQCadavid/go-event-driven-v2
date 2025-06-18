package http

import (
	"tickets/message"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	pubSub *message.PubSub,
) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		pubSub: pubSub,
	}

	e.POST("/tickets-confirmation", handler.PostTicketsConfirmation)

	return e
}
