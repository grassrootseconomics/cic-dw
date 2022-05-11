package syncer

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type tableCount struct {
	Count int `db:"count"`
}

func (s *Syncer) CacheSyncer(ctx context.Context, t *asynq.Task) error {
	_, err := s.db.Exec(ctx, s.queries["cache-syncer"])
	if err != nil {
		log.Err(err).Msg("cache syncer task failed")
		return asynq.SkipRetry
	}

	var table tableCount
	if err := pgxscan.Get(ctx, s.db, &table, "SELECT COUNT(*) from transactions"); err != nil {
		log.Err(err).Msg("cache syncer task failed")
		return asynq.SkipRetry
	}

	log.Info().Msgf("=> %d transactions synced", table.Count)

	return nil
}
