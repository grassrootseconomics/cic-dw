package syncer

import (
	cic_net "github.com/grassrootseconomics/cic-go/net"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nleof/goyesql"
)

type Syncer struct {
	db           *pgxpool.Pool
	rClient      asynq.RedisConnOpt
	cicnetClient *cic_net.CicNet
	queries      goyesql.Queries
}

func New(db *pgxpool.Pool, rClient asynq.RedisConnOpt, cicnetClient *cic_net.CicNet, queries goyesql.Queries) *Syncer {
	return &Syncer{
		db:           db,
		rClient:      rClient,
		cicnetClient: cicnetClient,
		queries:      queries,
	}
}
