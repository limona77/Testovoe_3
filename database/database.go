package database

import (
	"Testovoe_3/config"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPool interface {
	Close()
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	CopyFrom(
		ctx context.Context,
		tableName pgx.Identifier,
		columnNames []string,
		rowSrc pgx.CopyFromSource) (int64, error)
	Ping(ctx context.Context) error
}
type DB struct {
	Pool PgxPool
}

func Connect() *DB {
	cfg := config.NewConfig()
	db := &DB{}
	var err error
	db.Pool, err = pgxpool.New(context.Background(), cfg.URL)
	if err != nil {
		panic("can't connect to Postgres")
	}
	_, err = pgx.Connect(context.Background(), cfg.URL)
	if err != nil {
		panic("can't connect to Postgres")
	}
	return db
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
