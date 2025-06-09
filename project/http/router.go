package http

import (
	"tickets/worker"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	workAgent *worker.Worker,
) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		workAgent: workAgent,
	}

	e.POST("/tickets-confirmation", handler.PostTicketsConfirmation)

	return e
}
