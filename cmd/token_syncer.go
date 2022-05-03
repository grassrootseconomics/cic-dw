package main

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v4"
	"github.com/lmittmann/w3"
	"github.com/rs/zerolog/log"
	"math/big"
	"strconv"
)

type tokenSyncer struct {
	app *App
}

type tokenCursor struct {
	CursorPos string `db:"cursor_pos"`
}

func newTokenSyncer(app *App) *tokenSyncer {
	return &tokenSyncer{
		app: app,
	}
}

func (s *tokenSyncer) ProcessTask(ctx context.Context, t *asynq.Task) error {
	log.Info().Msgf("running task type: %s", t.Type())
	var lastCursor tokenCursor

	if err := pgxscan.Get(ctx, s.app.db, &lastCursor, s.app.queries["token-cursor-pos"]); err != nil {
		return err
	}
	latestChainIdx, err := s.app.cicnetClient.EntryCount(ctx)
	if err != nil {
		return err
	}

	lastCursorPos, err := strconv.ParseInt(lastCursor.CursorPos, 10, 64)
	if err != nil {
		return err
	}

	latestChainPos := latestChainIdx.Int64() - 1
	log.Info().Msgf("current db cursor: %s, latest chain pos: %d", lastCursor.CursorPos, latestChainPos)
	if latestChainPos >= lastCursorPos {
		batch := &pgx.Batch{}

		for i := lastCursorPos; i <= latestChainPos; i++ {
			nextTokenAddress, err := s.app.cicnetClient.AddressAtIndex(ctx, big.NewInt(i))
			if err != nil {
				return err
			}

			tokenInfo, err := s.app.cicnetClient.TokenInfo(ctx, w3.A(fmt.Sprintf("0x%s", nextTokenAddress)))
			if err != nil {
				return err
			}

			batch.Queue(s.app.queries["insert-token-data"], nextTokenAddress, tokenInfo.Name, tokenInfo.Symbol, tokenInfo.Decimals.Int64())
		}

		res := s.app.db.SendBatch(ctx, batch)
		log.Info().Msgf("inserting %d new records", batch.Len())
		for i := 0; i < batch.Len(); i++ {
			_, err := res.Exec()
			if err != nil {
				return err
			}
		}
		err := res.Close()
		if err != nil {
			return err
		}

		_, err = s.app.db.Exec(ctx, s.app.queries["update-token-cursor"], strconv.FormatInt(latestChainIdx.Int64(), 10))
		if err != nil {
			return err
		}
	}

	return nil
}
