package main

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
	CursorPos string `db:"cursor_pos"`
}

func tokenSyncer(ctx context.Context, t *asynq.Task) error {
	var lastCursor tokenCursor

	if err := pgxscan.Get(ctx, db, &lastCursor, queries["cursor-pos"], 3); err != nil {
		return err
	}
	latestChainIdx, err := cicnetClient.EntryCount(ctx)
	if err != nil {
		return err
	}

	lastCursorPos, err := strconv.ParseInt(lastCursor.CursorPos, 10, 64)
	if err != nil {
		return err
	}

	latestChainPos := latestChainIdx.Int64() - 1
	log.Info().Msgf("=> %d tokens synced", lastCursorPos)
	if latestChainPos >= lastCursorPos {
		batch := &pgx.Batch{}

		for i := lastCursorPos; i <= latestChainPos; i++ {
			nextTokenAddress, err := cicnetClient.AddressAtIndex(ctx, big.NewInt(i))
			if err != nil {
				return err
			}
			tokenInfo, err := cicnetClient.TokenInfo(ctx, w3.A(nextTokenAddress))
			if err != nil {
				return err
			}

			batch.Queue(queries["insert-token-data"], nextTokenAddress[2:], tokenInfo.Name, tokenInfo.Symbol, tokenInfo.Decimals.Int64())
		}

		res := db.SendBatch(ctx, batch)
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

		_, err = db.Exec(ctx, queries["update-cursor"], strconv.FormatInt(latestChainIdx.Int64(), 10), 3)
		if err != nil {
			return err
		}
	}

	return nil
}
