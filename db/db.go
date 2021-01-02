// Package db provides Pool for communicationg with Postgre database.
//
// Example usage:
//
//   // Preparing database connection pool
//   DBPool, err := db.New(
//   	fmt.Sprintf("host=%s port=%s user=%s password=%s database=%s",
//   		config.Database.Host,
//   		config.Database.Port,
//   		config.Database.User,
//   		config.Database.Password,
//   		config.Database.Database),
//   	logger,
//   	config.Database.Connections.MaxIdle,
//   	config.Database.Connections.MaxLife,
//   	config.Database.Connections.MaxOpenIdle,
//   	config.Database.Connections.MaxOpen,
//   )
//
//   // Creating table
//   _, err := DBPool.ExecContext(context.Background(), "CREATE TABLE IF NOT EXISTS dummy(name TEXT);")
//   if err != nil {
//   	fmt.Fprintf(os.Stderr, "table creation failed: %v\n", err)
//   	os.Exit(1)
//   }
//
//   // Adding data to table
//   _, err := DBPool.ExecContext(context.Background(), "insert into dummy(name) values($1)", "app item")
//   if err != nil {
//   	fmt.Fprintf(os.Stderr, "Insert failed: %v\n", err)
//   	os.Exit(1)
//   }
//
//   // Selecting data from table
//   rows, _ := DBPool.QueryContext(context.Background(), "select * from dummy")
//   for rows.Next() {
//   	var name string
//   	err := rows.Scan(&name)
//   	if err != nil {
//   		fmt.Fprintf(os.Stderr, "Scan failed: %v\n", err)
//   	}
//   	fmt.Printf("%s\n", name)
//   }
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
) (*Pool, error) {
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
