package postgresql

import (
	"BusServer/config"
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	Pool *pgxpool.Pool
)

func Initialize() {
	uri := fmt.Sprintf("postgresql://%s", strings.Join(config.PostgresAddresses, ","))
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	q := u.Query()
	q.Add("target_session_attrs", "read-write")
	u.RawQuery = q.Encode()
	uri = u.String()

	cfg, err := pgxpool.ParseConfig(uri)
	if err != nil {
		panic(err)
	}

	cfg.MaxConns = 50
	cfg.MaxConnLifetime = 60 * time.Second
	cfg.MaxConnIdleTime = 60 * time.Second
	cfg.ConnConfig.ConnectTimeout = 10 * time.Second
	cfg.ConnConfig.Database = config.PostgresDBName
	if config.PostgresUser != "" && config.PostgresPassword != "" {
		cfg.ConnConfig.User = config.PostgresUser
		cfg.ConnConfig.Password = config.PostgresPassword
	}

	Pool = setupPool(cfg)
}

func setupPool(cfg *pgxpool.Config) *pgxpool.Pool {
	log.Printf("[Postgresql] %v", cfg.ConnConfig.ConnString())
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		panic(err)
	}
	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	return pool
}

func Dispose() {
	Pool.Close()
}
