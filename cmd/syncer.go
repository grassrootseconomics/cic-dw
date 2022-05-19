package main

import (
	"cic-dw/internal/syncer"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func bootstrapScheduler(redis asynq.RedisConnOpt) (*asynq.Scheduler, error) {
	scheduler := asynq.NewScheduler(redis, nil)

	for k, v := range conf.Syncers {
		task := asynq.NewTask(k, nil)

		_, err := scheduler.Register(v, task)
		if err != nil {
			return nil, err
		}

		log.Info().Msgf("successfully registered %s syncer", k)
	}

	return scheduler, nil
}

func bootstrapProcessor(redis asynq.RedisConnOpt) (*asynq.Server, *asynq.ServeMux) {
	processorServer := asynq.NewServer(
		redis,
		asynq.Config{
			Concurrency: 5,
		},
	)

	syncer := syncer.New(db, redis, cicnetClient, preparedQueries.core)

	mux := asynq.NewServeMux()
	mux.HandleFunc("token", syncer.TokenSyncer)
	mux.HandleFunc("cache", syncer.CacheSyncer)
	mux.HandleFunc("ussd", syncer.UssdSyncer)

	return processorServer, mux
}
