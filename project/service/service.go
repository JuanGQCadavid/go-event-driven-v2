package service

import (
	"context"
	"errors"
	stdHTTP "net/http"

	"github.com/labstack/echo/v4"

	ticketsHttp "tickets/http"
	"tickets/worker"
)

type Service struct {
	echoRouter *echo.Echo
}

func New(
	spreadsheetsAPI ticketsHttp.SpreadsheetsAPI,
	receiptsService ticketsHttp.ReceiptsService,
) Service {

	// Initialazing and starting worker
	var workerAgent = worker.NewWorker(spreadsheetsAPI, receiptsService)
	go workerAgent.Run(context.Background())

	echoRouter := ticketsHttp.NewHttpRouter(workerAgent)

	return Service{
		echoRouter: echoRouter,
	}
}

func (s Service) Run(ctx context.Context) error {
	err := s.echoRouter.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}

	return nil
}
