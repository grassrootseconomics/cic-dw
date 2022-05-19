package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"time"

	batch_balance "github.com/grassrootseconomics/cic-go/batch_balance"
	cic_net "github.com/grassrootseconomics/cic-go/net"
	"github.com/grassrootseconomics/cic-go/provider"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/knadh/koanf"
	"github.com/lmittmann/w3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sys/unix"
)

var (
	k = koanf.New(".")

	preparedQueries *queries
	conf            config
	db              *pgxpool.Pool
	rpcProvider     *provider.Provider
	cicnetClient    *cic_net.CicNet
	batchBalance    *batch_balance.BatchBalance
	rClient         asynq.RedisConnOpt
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := loadConfig("config.toml", k); err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	if err := loadQueries("queries"); err != nil {
		log.Fatal().Err(err).Msg("failed to load sql file")
	}

	if err := connectDb(conf.Db.Postgres); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}

	if err := loadProvider(conf.Chain.RpcProvider); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}

	if err := loadCicNet(w3.A(conf.Chain.TokenRegistry)); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}

	if err := loadBatchBalance(w3.A(conf.Chain.BalanceResolver)); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}

	if err := parseRedis(conf.Db.Redis); err != nil {

		log.Fatal().Err(err).Msg("could not parse redis connection string")
	}
}

func main() {
	scheduler, err := bootstrapScheduler(rClient)
	if err != nil {
		log.Fatal().Err(err).Msg("could not bootstrap scheduler")
	}

	go func() {
		if err := scheduler.Run(); err != nil {
			log.Fatal().Err(err).Msg("could not start scheduler")
		}
	}()

	processor, mux := bootstrapProcessor(rClient)
	go func() {
		if err := processor.Run(mux); err != nil {
			log.Fatal().Err(err).Msg("failed to start job processor")
		}
	}()

	server := initHTTPServer()
	go func() {
		if err := server.Start(conf.Server.Address); err != nil {
			if strings.Contains(err.Error(), "Server closed") {
				log.Info().Msg("shutting down server")
			} else {
				log.Fatal().Err(err).Msg("could not start server")
			}
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT, unix.SIGTSTP)
	for {
		s := <-sigs
		if s == unix.SIGTSTP {
			processor.Stop()
			scheduler.Shutdown()
			continue
		}
		break
	}

	processor.Shutdown()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("could not shut down server")
	}
	log.Info().Msg("gracefully shutdown processor, scheduler and server")
}
