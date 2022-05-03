package main

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

var scheduler *asynq.Scheduler

func runScheduler(app *App) {
	scheduler = asynq.NewScheduler(app.rClient, nil)

	// TODO: Refactor boilerplate and pull enabled tasks from koanf
	tokenTask := asynq.NewTask("token:sync", nil)
	cacheTask := asynq.NewTask("cache:sync", nil)
	ussdTask := asynq.NewTask("ussd:sync", nil)

	_, err := scheduler.Register(conf.String("token.schedule"), tokenTask)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register token syncer")
	}
	log.Info().Msg("successfully registered token syncer")

	_, err = scheduler.Register(conf.String("cache.schedule"), cacheTask)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register cache syncer")
	}
	log.Info().Msg("successfully registered cache syncer")

	_, err = scheduler.Register(conf.String("ussd.schedule"), ussdTask)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register ussd syncer")
	}
	log.Info().Msg("successfully registered ussd syncer")

	if err := scheduler.Run(); err != nil {
		log.Fatal().Err(err).Msg("could not start asynq scheduler")
	}
}
