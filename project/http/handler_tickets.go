package http

import (
	"net/http"
	"tickets/message"

	"github.com/labstack/echo/v4"
)

type ticketsConfirmationRequest struct {
	Tickets []string `json:"tickets"`
}

func (h Handler) PostTicketsConfirmation(c echo.Context) error {
	var request ticketsConfirmationRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		h.pubSub.SendMessages(
			message.Msg{
				Payload: ticket,
				Topic:   "issue-receipt",
			},
			message.Msg{
				Payload: ticket,
				Topic:   "append-to-tracker",
			},
		)
	}

	return c.NoContent(http.StatusOK)
}
