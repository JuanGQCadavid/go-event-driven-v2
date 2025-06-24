package worker

import (
	"context"
	"encoding/json"
	"tickets/entities"
	"tickets/message"
)

type Task int

const (
	TaskIssueReceipt Task = iota
	TaskAppendToTracker
)

// type Message struct {
// 	Task     Task
// 	TicketID string
// }

type Worker struct {
	// queue chan Message

	spreadsheetsAPI SpreadsheetsAPI
	receiptsService ReceiptsService

	pubSub *message.PubSub
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) error
}

func NewWorker(
	spreadsheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
	pubSub *message.PubSub,
) *Worker {
	return &Worker{
		spreadsheetsAPI: spreadsheetsAPI,
		receiptsService: receiptsService,
		pubSub:          pubSub,
	}
}

func (w *Worker) Init() {

	w.pubSub.Subscribe("TicketBookingConfirmed", "issue-receipt", w.HandleReceipt)
	w.pubSub.Subscribe("TicketBookingConfirmed", "append-to-tracker", w.HandleSpread)
	w.pubSub.Subscribe("TicketBookingCanceled", "append-to-refund", w.HandleRefound)
}
func (w *Worker) castPaylaod(payload string) (*entities.TicketBookingConfirmed, error) {
	var (
		ticket *entities.TicketBookingConfirmed = &entities.TicketBookingConfirmed{}
	)

	if err := json.Unmarshal([]byte(payload), ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (w *Worker) HandleReceipt(ctx context.Context, payload string) error {

	ticket, err := w.castPaylaod(payload)

	if err != nil {
		return err
	}

	return w.receiptsService.IssueReceipt(ctx, entities.IssueReceiptRequest{
		TicketID: ticket.TicketID,
		Price:    ticket.Price,
	})
}

func (w *Worker) HandleRefound(ctx context.Context, payload string) error {
	ticket, err := w.castPaylaod(payload)

	if err != nil {
		return err
	}

	return w.spreadsheetsAPI.AppendRow(ctx, "tickets-to-refund", []string{ticket.TicketID, ticket.CustomerEmail, ticket.Price.Amount, ticket.Price.Currency})
}

func (w *Worker) HandleSpread(ctx context.Context, payload string) error {
	ticket, err := w.castPaylaod(payload)

	if err != nil {
		return err
	}

	return w.spreadsheetsAPI.AppendRow(ctx, "tickets-to-print", []string{ticket.TicketID, ticket.CustomerEmail, ticket.Price.Amount, ticket.Price.Currency})
}
