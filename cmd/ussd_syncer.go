package main

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type ussdSyncer struct {
	app *App
}

func newUssdSyncer(app *App) *ussdSyncer {
	return &ussdSyncer{
		app: app,
	}
}

func (s *ussdSyncer) ProcessTask(ctx context.Context, t *asynq.Task) error {
	_, err := s.app.db.Exec(ctx, s.app.queries["ussd-syncer"])
	if err != nil {
		return asynq.SkipRetry
	}

	var count tableCount
	if err := pgxscan.Get(ctx, s.app.db, &count, "SELECT COUNT(*) from users"); err != nil {
		return asynq.SkipRetry
	}

	log.Info().Msgf("=> %d users synced", count.Count)

	return nil
}
