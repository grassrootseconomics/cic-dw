package main

import (
	"cic-dw/pkg/cicnet"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/knadh/koanf"
	"github.com/lmittmann/w3"
	"github.com/nleof/goyesql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

type App struct {
	db           *pgxpool.Pool
	queries      goyesql.Queries
	rClient      asynq.RedisClientOpt
	cicnetClient *cicnet.CicNet
	sigChan      chan os.Signal
}

const (
	confEnvOverridePrefix = ""
)

var (
	conf         = koanf.New(".")
	db           *pgxpool.Pool
	queries      goyesql.Queries
	redisConn    asynq.RedisClientOpt
	cicnetClient *cicnet.CicNet
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := loadConfig("config.toml", confEnvOverridePrefix, conf); err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	db = connectDb(conf.String("db.dsn"))
	queries = loadQueries("queries.sql")
	redisConn = connectQueue(conf.String("redis.dsn"))
	cicnetClient = cicnet.NewCicNet(conf.String("chain.rpc"), w3.A(conf.String("chain.registry")))
}

func main() {
	// TODO: Graceful shutdown of go routines (handle SIG INT/TERM)
	var wg sync.WaitGroup

	app := &App{
		db:           db,
		queries:      queries,
		rClient:      redisConn,
		cicnetClient: cicnetClient,
	}

	wg.Add(2)
	go runScheduler(app)
	go runProcessor(app)
	wg.Wait()
}
