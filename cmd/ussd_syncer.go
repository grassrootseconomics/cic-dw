package main

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func ussdSyncer(ctx context.Context, t *asynq.Task) error {
	_, err := db.Exec(ctx, queries["ussd-syncer"])
	if err != nil {
		return asynq.SkipRetry
	}

	var count tableCount
	if err := pgxscan.Get(ctx, db, &count, "SELECT COUNT(*) from users"); err != nil {
		return asynq.SkipRetry
	}

	log.Info().Msgf("=> %d users synced", count.Count)

	return nil
}
