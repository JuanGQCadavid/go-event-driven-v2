package http

import (
	"encoding/json"
	"net/http"
	"tickets/entities"
	"tickets/message"

	"github.com/labstack/echo/v4"
)

func (h Handler) PostTicketsConfirmation(c echo.Context) error {
	var request entities.TicketsStatusRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		var topic string = "TicketBookingConfirmed"

		if ticket.Status != "confirmed" {
			// TODO - They use an other struct TicketBookingCanceled
			topic = "TicketBookingCanceled"
		}

		h.pubSub.SendMessages(
			message.Msg{
				Payload: castTicketBookingConfirmed(ticket),
				Topic:   topic,
			},
		)
	}

	return c.NoContent(http.StatusOK)
}

func castTicketBookingConfirmed(ticket entities.TicketStatusRequest) string {
	structPayload := &entities.TicketBookingConfirmed{
		Header:        entities.NewEventHeader(),
		TicketID:      ticket.TicketID,
		CustomerEmail: ticket.CustomerEmail,
		Price:         ticket.Price,
	}

	if payload, err := json.Marshal(structPayload); err != nil {
		return ticket.TicketID
	} else {
		return string(payload)
	}
}
