package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	cic_net "github.com/grassrootseconomics/cic-go/net"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/nleof/goyesql"
	"strings"
)

type config struct {
	Db struct {
		Postgres string `koanf:"postgres"`
		Redis    string `koanf:"redis"`
	}
	Server struct {
		Address string   `koanf:"address"`
		Cors    []string `koanf:"cors"`
	}
	Chain struct {
		RpcProvider   string `koanf:"rpc"`
		TokenRegistry string `koanf:"index"`
	}
	Syncers map[string]string `koanf:"syncers"`
}

type queries struct {
	core      goyesql.Queries
	dashboard goyesql.Queries
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

	cicnetClient, err = cic_net.NewCicNet(rpcProvider, tokenIndex)
	if err != nil {
		return err
	}

	return nil
}

func loadQueries(sqlFilesPath string) error {
	coreQueries, err := goyesql.ParseFile(fmt.Sprintf("%s/core.sql", sqlFilesPath))
	if err != nil {
		return err
	}

	dashboardQueries, err := goyesql.ParseFile(fmt.Sprintf("%s/dashboard.sql", sqlFilesPath))
	if err != nil {
		return err
	}

	preparedQueries = &queries{
		core:      coreQueries,
		dashboard: dashboardQueries,
	}

	return nil
}
