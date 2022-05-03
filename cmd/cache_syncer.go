package main

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type cacheSyncer struct {
	app *App
}

type tableCount struct {
	Count int `db:"count"`
}

func newCacheSyncer(app *App) *cacheSyncer {
	return &cacheSyncer{
		app: app,
	}
}

func (s *cacheSyncer) ProcessTask(ctx context.Context, t *asynq.Task) error {
	_, err := s.app.db.Exec(ctx, s.app.queries["cache-syncer"])
	if err != nil {
		return asynq.SkipRetry
	}

	var count tableCount
	if err := pgxscan.Get(ctx, s.app.db, &count, "SELECT COUNT(*) from transactions"); err != nil {
		return asynq.SkipRetry
	}

	log.Info().Msgf("=> %d transactions synced", count.Count)

	return nil
}
