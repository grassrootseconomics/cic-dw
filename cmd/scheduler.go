package main

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

var scheduler *asynq.Scheduler

func runScheduler(app *App) {
	scheduler = asynq.NewScheduler(app.rClient, nil)

	tokenTask := asynq.NewTask("token:sync", nil)

	_, err := scheduler.Register(conf.String("schedule.token"), tokenTask)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register token syncer")
	}
	log.Info().Msg("successfully registered token syncer")

	if err := scheduler.Run(); err != nil {
		log.Fatal().Err(err).Msg("could not start asynq scheduler")
	}
}
