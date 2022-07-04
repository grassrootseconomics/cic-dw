package syncer

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func (s *Syncer) UssdSyncer(ctx context.Context, _ *asynq.Task) error {
	_, err := s.db.Exec(ctx, s.queries["ussd-syncer"])
	if err != nil {
		log.Err(err).Msg("ussd syncer task failed")
		return asynq.SkipRetry
	}

	var table tableCount
	if err := pgxscan.Get(ctx, s.db, &table, "SELECT COUNT(*) from users"); err != nil {
		log.Err(err).Msg("ussd syncer task failed")
		return asynq.SkipRetry
	}

	log.Info().Msgf("=> %d users synced", table.Count)

	return nil
}
