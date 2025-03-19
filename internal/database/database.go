package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// minimal interfaces to work with postgres/sqlx  properly
type QueryContext interface {
	Exec(query string, params ...interface{}) (sql.Result, error)
	Query(query string, params ...interface{}) (*sql.Rows, error)
}

type Transaction interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	QueryContext
}

type Transactioner interface {
	Transaction
	Commit() error
	Rollback() error
}

type Database interface {
	BeginTx(context.Context) (Transactioner, error)
	Transaction
}

type PostgresDatabase struct {
	db *sqlx.DB
}

type PostgresTransaction struct {
	tx *sqlx.Tx
}

func (t *PostgresTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *PostgresTransaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *PostgresTransaction) Exec(query string, params ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, params...)
}

func (t *PostgresTransaction) Query(query string, params ...interface{}) (*sql.Rows, error) {
	return t.tx.Query(query, params...)
}

func (t *PostgresTransaction) Get(dest interface{}, query string, args ...interface{}) error {
	return t.tx.Get(dest, query, args...)
}

func (t *PostgresTransaction) Select(dest interface{}, query string, args ...interface{}) error {
	return t.tx.Select(dest, query, args...)
}

func NewPostgresDatabase(ctx context.Context, connString string) (Database, error) {
	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	return &PostgresDatabase{db: db}, nil
}

func (d *PostgresDatabase) BeginTx(ctx context.Context) (Transactioner, error) {
	tx, err := d.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	return &PostgresTransaction{tx: tx}, nil
}

func (d *PostgresDatabase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}

func (d *PostgresDatabase) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(query, args...)
}

func (d *PostgresDatabase) Get(dest interface{}, query string, args ...interface{}) error {
	return d.db.Get(dest, query, args...)
}

func (d *PostgresDatabase) Select(dest interface{}, query string, args ...interface{}) error {
	return d.db.Select(dest, query, args...)
}
