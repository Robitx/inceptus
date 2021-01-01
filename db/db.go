package db

import (
	"time"

	sqlx "github.com/jmoiron/sqlx"

	pgx "github.com/jackc/pgx/v4"
	pgx_zerolog "github.com/jackc/pgx/v4/log/zerologadapter"
	pgx_std "github.com/jackc/pgx/v4/stdlib"

	log "github.com/robitx/inceptus/log"
)

// Pool of connections for interaction with Database
type Pool struct {
	*sqlx.DB
}

// New creates Pool of tcp connections to database
// 
// DSN for TCP conn is specified as:
// fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s", dbTcpHost, dbUser, dbPwd, dbPort, dbName)
// 
// DSN for Unix socket is specified as:
// fmt.Sprintf("user=%s password=%s database=%s host=%s/%s", dbUser, dbPwd, dbName, socketDir, instanceConnectionName)
func New(
	dsn string,
	logger *log.Logger,
	connMaxIdle time.Duration,
	connMaxLife time.Duration,
	connMaxOpenIdle int,
	connMaxOpen int,
	) (*Pool, error){

	// prep config
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// prep logger
	level, err := pgx.LogLevelFromString(logger.GetLevel().String()) 
	if err != nil {
		level = pgx.LogLevelInfo
	}
	config.LogLevel = level
	config.Logger = pgx_zerolog.NewLogger(logger.Logger)

	// make connections
	conn, err := sqlx.Open("pgx", pgx_std.RegisterConnConfig(config))
	if err != nil {
		return nil, err
	}

	pool := &Pool{
		DB: conn,
	}

	// configure pool limits
	pool.SetConnMaxIdleTime(connMaxIdle)
	pool.SetConnMaxLifetime(connMaxLife)
	pool.SetMaxIdleConns(connMaxOpenIdle)
	pool.SetMaxOpenConns(connMaxOpen)

	return pool, nil
}