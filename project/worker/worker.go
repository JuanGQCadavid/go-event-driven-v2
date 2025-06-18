package worker

import (
	"context"
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
	IssueReceipt(ctx context.Context, ticketID string) error
}

func NewWorker(
	spreadsheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
	pubSub *message.PubSub,
) *Worker {
	return &Worker{
		// queue:           make(chan Message, 100),
		spreadsheetsAPI: spreadsheetsAPI,
		receiptsService: receiptsService,
		pubSub:          pubSub,
	}
}

// func (w *Worker) Send(msgs ...Message) {
// 	for _, msg := range msgs {
// 		w.queue <- msg
// 	}
// }

func (w *Worker) Init() {

	w.pubSub.Subscribe("issue-receipt", "issue-receipt", w.HandleReceipt)
	w.pubSub.Subscribe("append-to-tracker", "append-to-tracker", w.HandleSpread)
}

func (w *Worker) HandleReceipt(ctx context.Context, payload string) error {
	return w.receiptsService.IssueReceipt(ctx, payload)
}

func (w *Worker) HandleSpread(ctx context.Context, payload string) error {
	return w.spreadsheetsAPI.AppendRow(ctx, "tickets-to-print", []string{payload})
}

// func (w *Worker) Run(ctx context.Context) {
// 	for msg := range w.queue {
// 		var err error = nil
// 		switch msg.Task {
// 		case TaskIssueReceipt:
// 			err = w.receiptsService.IssueReceipt(ctx, msg.TicketID)
// 			if err != nil {
// 				slog.With("error", err).Error("failed to issue the receipt")
// 			}
// 		case TaskAppendToTracker:
// 			err = w.spreadsheetsAPI.AppendRow(ctx, "tickets-to-print", []string{msg.TicketID})
// 			if err != nil {
// 				slog.With("error", err).Error("failed to append to tracker")
// 			}
// 		}
// 		if err != nil {
// 			w.Send(msg)
// 		}
// 	}
// }
