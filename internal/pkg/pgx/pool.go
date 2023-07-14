package pgx

import (
	"context"
	"log"
	"time"

	"tech-tsarka/internal/config"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPool(ctx context.Context, cfg config.PgxPool, dsn string) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	poolCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		cache := conn.StatementCache()
		prefix := "/internal/pkg/pgx/pool.go:poolConfig.AfterConnect(): Statement cache mode:"
		if cache.Mode() == 0 {
			log.Println(prefix, "prepare")
		} else {
			log.Println(prefix, "describe")
		}
		return nil
	}

	poolCfg.MaxConnIdleTime = time.Second * time.Duration(cfg.MaxConnIdleTime)
	poolCfg.MaxConnLifetime = time.Second * time.Duration(cfg.MaxConnLifetime)
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConns = cfg.MaxConns

	pool, err := pgxpool.ConnectConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
