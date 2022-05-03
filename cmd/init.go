package main

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/nleof/goyesql"
	"github.com/rs/zerolog/log"
	"strings"
)

func loadConfig(configFilePath string, envOverridePrefix string, conf *koanf.Koanf) error {
	// assumed to always be at the root folder
	confFile := file.Provider(configFilePath)
	if err := conf.Load(confFile, toml.Parser()); err != nil {
		return err
	}
	// override with env variables
	if err := conf.Load(env.Provider(envOverridePrefix, ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, envOverridePrefix)), "_", ".")
	}), nil); err != nil {
		return err
	}

	return nil
}

func connectDb(dsn string) *pgxpool.Pool {
	conn, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}

	return conn
}

func loadQueries(sqlFile string) goyesql.Queries {
	q, err := goyesql.ParseFile(sqlFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse sql queries")
	}

	return q
}

func connectQueue(dsn string) asynq.RedisClientOpt {
	rClient := asynq.RedisClientOpt{Addr: dsn}

	return rClient
}
