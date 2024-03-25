package pgs

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"route256.ozon.ru/project/loms/internal/config"
	"route256.ozon.ru/project/loms/internal/repository/pgs/queries"
	"strings"
	"sync/atomic"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectToPgsDb(ctx context.Context, config *config.Config, masterOnly bool) (*DB, error) {
	connString := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s;",
		config.LomsDbUser,
		config.LomsDbPass,
		config.LomsDbHost,
		config.LomsDbName,
	)
	if masterOnly {
		connString += connString
		connString = strings.Trim(connString, ";")
	} else {
		connString += fmt.Sprintf(
			"postgresql://%s:%s@%s/%s",
			config.LomsDbSlaveUser,
			config.LomsDbSlavePass,
			config.LomsDbSlaveHost,
			config.LomsDbName,
		)
	}
	return ConnectByDataSourceNames(ctx, connString)
}

// https://github.com/tsenart/nap Взял отсюда и переделал
type DB struct {
	pdbs  []*pgxpool.Pool
	count uint64
}

func ConnectByDataSourceNames(ctx context.Context, dataSourceNames string) (*DB, error) {
	conns := strings.Split(dataSourceNames, ";")
	db := &DB{pdbs: make([]*pgxpool.Pool, len(conns))}

	err := scatter(len(db.pdbs), func(i int) (err error) {
		cfg, err := pgxpool.ParseConfig(conns[i])
		if err != nil {
			return fmt.Errorf("parse config: %w", err)
		}

		db.pdbs[i], err = pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			return fmt.Errorf("create connection pool: %w", err)
		}
		err = db.pdbs[i].Ping(ctx)
		if err != nil {
			return fmt.Errorf("cant ping connection pool: %w", err)
		}
		log.Printf("Connected to postgres db number %d", i)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) Close() {
	scatter(len(db.pdbs), func(i int) error {
		db.pdbs[i].Close()
		return nil
	})
}

func (db *DB) GetSelectQueries(ctx context.Context) queries.Querier {
	return queries.New(db.Slave())
}

func (db *DB) GetUpdateQueries(ctx context.Context) queries.Querier {
	return queries.New(db.Master())
}

func (db *DB) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.Master().Begin(ctx)
}

func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.Master().Exec(ctx, query, args...)
}

func (db *DB) Ping(ctx context.Context) error {
	return scatter(len(db.pdbs), func(i int) error {
		return db.pdbs[i].Ping(ctx)
	})
}

func (db *DB) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return db.Slave().Query(ctx, query, args...)
}

func (db *DB) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.Slave().QueryRow(ctx, query, args...)
}

func (db *DB) Slave() *pgxpool.Pool {
	return db.pdbs[db.slave(len(db.pdbs))]
}

func (db *DB) Master() *pgxpool.Pool {
	return db.pdbs[0]
}

func (db *DB) slave(n int) int {
	if n <= 1 {
		return 0
	}
	return int(1 + (atomic.AddUint64(&db.count, 1) % uint64(n-1)))
}

func scatter(n int, fn func(i int) error) error {
	errors := make(chan error, n)

	var i int
	for i = 0; i < n; i++ {
		go func(i int) { errors <- fn(i) }(i)
	}

	var err, innerErr error
	for i = 0; i < cap(errors); i++ {
		if innerErr = <-errors; innerErr != nil {
			err = innerErr
		}
	}

	return err
}
