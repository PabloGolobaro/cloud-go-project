package transact

import (
	"cloud-go-project/hexarch/core"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type PostgresTransactionLogger struct {
	events chan<- core.Event
	errors <-chan error
	db     *sql.DB
}

type PostgresDBParams struct {
	dbName   string
	host     string
	user     string
	password string
}

func NewPostgresTransactionLogger(config PostgresDBParams) (core.TransactionLogger, error) {
	connStr := fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s sslmode=disable",
		config.host,
		config.dbName,
		config.user,
		config.password,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failde to open db: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failde to open db conncetion: %w", err)
	}
	logger := &PostgresTransactionLogger{db: db}
	err = logger.verifyTableExists()
	if err != nil {
		if err := logger.createTable(); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}

	}
	return logger, nil

}
func (l *PostgresTransactionLogger) Run() {
	events := make(chan core.Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		query := `INSERT INTO transactions 
				(type,key,value)
				VALUES ($1,$2,$3)`
		for e := range events {
			_, err := l.db.Exec(query, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
			}
		}
	}()
}
func (l *PostgresTransactionLogger) ReadEvents() (<-chan core.Event, <-chan error) {
	outEvent := make(chan core.Event)
	outError := make(chan error)
	go func() {
		defer close(outEvent)
		defer close(outError)

		query := `SELECT * FROM transactions ORDER BY id`
		rows, err := l.db.Query(query)
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}
		defer rows.Close()
		e := core.Event{}
		for rows.Next() {
			err := rows.Scan(&e.Sequence, &e.EventType, &e.Key, &e.Value)
			if err != nil {
				outError <- fmt.Errorf("sql query error reading: %w", err)
				return
			}
			outEvent <- e
		}
		err = rows.Err()
		if err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()
	return outEvent, outError
}
func (l *PostgresTransactionLogger) WritePut(key string, value string) {
	l.events <- core.Event{EventType: core.EventPut, Key: key, Value: value}
}
func (l *PostgresTransactionLogger) WriteDelete(key string) {
	l.events <- core.Event{EventType: core.EventDelete, Key: key}
}
func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *PostgresTransactionLogger) verifyTableExists() error {
	query := `
	SELECT table_name
	FROM INFORMATION_SCHEMA.TABLES
	WHERE TABLE_TYPE = 'BASE TABLE'
  	AND TABLE_NAME = 'transactions';`
	row := l.db.QueryRow(query)
	var table_name string
	err := row.Scan(&table_name)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error: %w", err)
	}
	log.Println(table_name)
	log.Print(err)
	return nil

}

func (l *PostgresTransactionLogger) createTable() error {
	query := `
	CREATE TABLE transactions (
	id SERIAL PRIMARY KEY,
	type VARCHAR ( 50 ) NOT NULL,
	key VARCHAR ( 50 ) NOT NULL,
	value VARCHAR ( 255 ));`
	_, err := l.db.Exec(query)
	return err
}
