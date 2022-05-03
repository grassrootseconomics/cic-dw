package main

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func runProcessor(app *App) {
	processorServer := asynq.NewServer(
		app.rClient,
		asynq.Config{
			Concurrency: 10,
		},
	)

	mux := asynq.NewServeMux()
	mux.Handle("token:sync", newTokenSyncer(app))
	mux.Handle("cache:sync", newCacheSyncer(app))
	mux.Handle("ussd:sync", newUssdSyncer(app))

	if err := processorServer.Run(mux); err != nil {
		log.Fatal().Err(err).Msg("failed to start job processor")
	}
}
