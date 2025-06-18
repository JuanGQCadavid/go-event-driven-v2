package http

import (
	"context"
	"tickets/message"
)

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, ticketID string) error
}

type Handler struct {
	pubSub *message.PubSub
}
