package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"golang.org/x/sync/errgroup"

	"tickets/adapters"
	"tickets/message"
	"tickets/service"
	"tickets/worker"
)

func main() {
	log.Init(slog.LevelInfo)

	apiClients, err := clients.NewClients(os.Getenv("GATEWAY_ADDR"), nil)
	if err != nil {
		panic(err)
	}

	spreadsheetsAPI := adapters.NewSpreadsheetsAPIClient(apiClients)
	receiptsService := adapters.NewReceiptsServiceClient(apiClients)

	// Context
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	// Pub SUb
	pubSub := message.NewPubSub()

	g.Go(func() error {
		return pubSub.Init(ctx)
	})

	fmt.Println("Waiting")
	pubSub.WaitUntilReady()
	fmt.Println("Done")

	// Worker
	var workerAgent = worker.NewWorker(spreadsheetsAPI, receiptsService, pubSub)
	workerAgent.Init()

	// Service

	srv := service.New(
		pubSub,
	)

	g.Go(func() error {
		return srv.Run(ctx)
	})

	g.Go(func() error {
		<-ctx.Done()
		return srv.ShutDown(ctx)
	})

	err = g.Wait()

	if err != nil {
		panic(err)
	}

	// err = service.New(
	// 	pubSub,
	// ).Run(context.Background())
	// if err != nil {
	// 	panic(err)
	// }
}
