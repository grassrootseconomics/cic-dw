package main

import (
	"cic-dw/pkg/cicnet"
	"context"
	"github.com/ethereum/go-ethereum/common"
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

type config struct {
	Db struct {
		Postgres string `koanf:"postgres"`
		Redis    string `koanf:"redis"`
	}
	Chain struct {
		RpcProvider   string `koanf:"rpc"`
		TokenRegistry string `koanf:"index"`
	}
	Syncers map[string]string `koanf:"syncers"`
}

func loadConfig(configFilePath string, k *koanf.Koanf) error {
	confFile := file.Provider(configFilePath)
	if err := k.Load(confFile, toml.Parser()); err != nil {
		return err
	}
	if err := k.Load(env.Provider("", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, "")), "_", ".")
	}), nil); err != nil {
		return err
	}

	err := k.UnmarshalWithConf("", &conf, koanf.UnmarshalConf{Tag: "koanf"})
	if err != nil {
		return err
	}

	return nil
}

func connectDb(dsn string) error {
	var err error
	db, err = pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return err
	}

	return nil
}

func parseRedis(dsn string) (asynq.RedisConnOpt, error) {
	rconn, err := asynq.ParseRedisURI(dsn)
	if err != nil {
		return nil, err
	}

	return rconn, nil
}

func connectCicNet(rpcProvider string, tokenIndex common.Address) error {
	var err error

	cicnetClient, err = cicnet.NewCicNet(rpcProvider, tokenIndex)
	if err != nil {
		return err
	}

	return nil
}

func loadQueries(sqlFile string) error {
	var err error
	queries, err = goyesql.ParseFile(sqlFile)
	if err != nil {
		return err
	}

	return nil
}

func bootstrapScheduler(redis asynq.RedisConnOpt) (*asynq.Scheduler, error) {
	scheduler := asynq.NewScheduler(redis, nil)

	for k, v := range conf.Syncers {
		task := asynq.NewTask(k, nil)

		_, err := scheduler.Register(v, task)
		if err != nil {
			return nil, err
		}

		log.Info().Msgf("successfully registered %s syncer", k)
	}

	return scheduler, nil
}

func bootstrapProcessor(redis asynq.RedisConnOpt) (*asynq.Server, *asynq.ServeMux) {
	processorServer := asynq.NewServer(
		redis,
		asynq.Config{
			Concurrency: 5,
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("token", tokenSyncer)
	mux.HandleFunc("cache", cacheSyncer)
	mux.HandleFunc("ussd", ussdSyncer)

	return processorServer, mux
}
