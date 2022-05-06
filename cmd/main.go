package main

import (
	"github.com/grassrootseconomics/cic_go/cic_net"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/knadh/koanf"
	"github.com/lmittmann/w3"
	"github.com/nleof/goyesql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
)

var (
	k = koanf.New(".")

	queries      goyesql.Queries
	conf         config
	db           *pgxpool.Pool
	cicnetClient *cic_net.CicNet
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := loadConfig("config.toml", k); err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	if err := loadQueries("queries.sql"); err != nil {
		log.Fatal().Err(err).Msg("failed to load sql file")
	}

	if err := connectDb(conf.Db.Postgres); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}

	if err := connectCicNet(conf.Chain.RpcProvider, w3.A(conf.Chain.TokenRegistry)); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}
}

func main() {
	rClient, err := parseRedis(conf.Db.Redis)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse redis connection string")
	}

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
	log.Info().Msg("gracefully shutdown processor and scheduler")
}
