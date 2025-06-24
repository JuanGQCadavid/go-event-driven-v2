package http

import (
	"context"
	"tickets/entities"
	"tickets/message"
)

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) error
}

type Handler struct {
	pubSub *message.PubSub
}
