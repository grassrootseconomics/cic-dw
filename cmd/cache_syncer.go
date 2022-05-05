package main

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type tableCount struct {
	Count int `db:"count"`
}

func cacheSyncer(ctx context.Context, t *asynq.Task) error {
	_, err := db.Exec(ctx, queries["cache-syncer"])
	if err != nil {
		return asynq.SkipRetry
	}

	var count tableCount
	if err := pgxscan.Get(ctx, db, &count, "SELECT COUNT(*) from transactions"); err != nil {
		return asynq.SkipRetry
	}

	log.Info().Msgf("=> %d transactions synced", count.Count)

	return nil
}
