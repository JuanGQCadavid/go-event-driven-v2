package service

import (
	"context"
	"errors"
	stdHTTP "net/http"

	"github.com/labstack/echo/v4"

	ticketsHttp "tickets/http"
	"tickets/message"
)

type Service struct {
	echoRouter *echo.Echo
}

func New(
	pubSub *message.PubSub,
) Service {
	echoRouter := ticketsHttp.NewHttpRouter(pubSub)

	return Service{
		echoRouter: echoRouter,
	}
}

func (s Service) ShutDown(ctx context.Context) error {
	return s.echoRouter.Shutdown(ctx)
}

func (s Service) Run(ctx context.Context) error {
	err := s.echoRouter.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}

	return nil
}
