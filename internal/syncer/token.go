package syncer

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v4"
	"github.com/lmittmann/w3"
	"github.com/rs/zerolog/log"
	"math/big"
	"strconv"
)

type tokenCursor struct {
	cursorPos string `db:"cursor_pos"`
}

func (s *Syncer) TokenSyncer(ctx context.Context, t *asynq.Task) error {
	var lastCursor tokenCursor

	if err := pgxscan.Get(ctx, s.db, &lastCursor, s.queries["cursor-pos"], 3); err != nil {
		return err
	}
	latestChainIdx, err := s.cicnetClient.EntryCount(ctx)
	if err != nil {
		log.Err(err).Msg("token syncer task failed")
		return err
	}

	lastCursorPos, err := strconv.ParseInt(lastCursor.cursorPos, 10, 64)
	if err != nil {
		log.Err(err).Msg("token syncer task failed")
		return err
	}

	latestChainPos := latestChainIdx.Int64() - 1
	log.Info().Msgf("=> %d tokens synced", lastCursorPos)
	if latestChainPos >= lastCursorPos {
		batch := &pgx.Batch{}

		for i := lastCursorPos; i <= latestChainPos; i++ {
			nextTokenAddress, err := s.cicnetClient.AddressAtIndex(ctx, big.NewInt(i))
			if err != nil {
				log.Err(err).Msg("token syncer task failed")
				return err
			}
			tokenInfo, err := s.cicnetClient.ERC20TokenInfo(ctx, w3.A(nextTokenAddress))
			if err != nil {
				log.Err(err).Msg("token syncer task failed")
				return err
			}

			batch.Queue(s.queries["insert-token-data"], nextTokenAddress[2:], tokenInfo.Name, tokenInfo.Symbol, tokenInfo.Decimals.Int64())
		}

		res := s.db.SendBatch(ctx, batch)
		for i := 0; i < batch.Len(); i++ {
			_, err := res.Exec()
			if err != nil {
				log.Err(err).Msg("token syncer task failed")
				return err
			}
		}
		err := res.Close()
		if err != nil {
			log.Err(err).Msg("token syncer task failed")
			return err
		}

		_, err = s.db.Exec(ctx, s.queries["update-cursor"], strconv.FormatInt(latestChainIdx.Int64(), 10), 3)
		if err != nil {
			log.Err(err).Msg("token syncer task failed")
			return err
		}
	}

	return nil
}
