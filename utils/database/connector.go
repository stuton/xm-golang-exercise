package database

import (
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

const driverName = "pgx"

type newDBOptions struct {
	logger           *zapadapter.Logger
	maxIdleConns     int
	maxIdleConnsTime time.Duration
	maxOpenConns     int
	maxOpenConnsTime time.Duration
}

type NewDBOption func(o *newDBOptions)

func WithLogger(logger *zapadapter.Logger) NewDBOption {
	return func(o *newDBOptions) {
		o.logger = logger
	}
}

func WithMaxOpenConnections(max int) NewDBOption {
	return func(o *newDBOptions) {
		o.maxOpenConns = max
	}
}

func WithMaxOpenConnectionsTime(max time.Duration) NewDBOption {
	return func(o *newDBOptions) {
		o.maxOpenConnsTime = max
	}
}

func WithMaxIdleConnections(max int) NewDBOption {
	return func(o *newDBOptions) {
		o.maxIdleConns = max
	}
}

func WithMaxIdleConnectionsTime(max time.Duration) NewDBOption {
	return func(o *newDBOptions) {
		o.maxIdleConnsTime = max
	}
}

func NewDb(masterURI string, opts ...NewDBOption) (*sqlx.DB, error) {
	o := &newDBOptions{
		maxIdleConns:     10,
		maxIdleConnsTime: time.Duration(10),
		maxOpenConns:     10,
		maxOpenConnsTime: time.Duration(10),
	}

	for _, opt := range opts {
		opt(o)
	}

	masterConnect, err := makeConnect(masterURI, o)
	if err != nil {
		return nil, err
	}
	if err := masterConnect.Ping(); err != nil {
		return nil, err
	}

	return masterConnect, nil
}

func makeConnect(uri string, o *newDBOptions) (*sqlx.DB, error) {
	var config, err = pgx.ParseConfig(uri)
	if err != nil {
		return nil, err
	}

	config.PreferSimpleProtocol = true
	config.Logger = o.logger

	var connectionPool = stdlib.OpenDB(*config)

	connectionPool.SetMaxOpenConns(o.maxOpenConns)
	connectionPool.SetMaxIdleConns(o.maxIdleConns)
	connectionPool.SetConnMaxLifetime(o.maxOpenConnsTime)
	connectionPool.SetConnMaxIdleTime(o.maxIdleConnsTime)

	return sqlx.NewDb(connectionPool, driverName), nil
}
