package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Rows interface {
	pgx.Rows
}

type Row interface {
	pgx.Row
}

type CommandTag interface {
	RowsAffected() int64
	String() string
	Insert() bool
	Update() bool
	Delete() bool
	Select() bool
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	PoolSize string
	SSLMode  string
}

func (c Config) PostgresDSN() string {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", c.User, c.Password, c.Host, c.Port, c.Database)
	if c.SSLMode == "disable" {
		connStr += "?sslmode=disable"
	}

	return connStr
}

type Database struct {
	connection *pgxpool.Pool
}

func (d *Database) Query(ctx context.Context, sql string, args ...any) (Rows, error) {
	return d.connection.Query(ctx, sql, args...)
}

func (d *Database) QueryRow(ctx context.Context, sql string, args ...any) Row {
	return d.connection.QueryRow(ctx, sql, args...)
}

func (d *Database) Exec(ctx context.Context, sql string, args ...any) (CommandTag, error) {
	return d.connection.Exec(ctx, sql, args...)
}

func New(config Config) (*Database, error) {
	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_max_conns=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.SSLMode,
		config.PoolSize,
	)

	pgxConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, err
	}

	return &Database{
		connection: conn,
	}, nil
}
