package cmd

import (
	"fmt"
)

type EventType byte

var logger TransactionLogger

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)

type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     string
}
type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key string, value string)
	Err() <-chan error

	ReadEvents() (<-chan Event, <-chan error)
	Run()
}

func InitializeFileTransactionLog() error {
	var err error
	logger, err = NewFileTransactionLogger("transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}
	events, errors := logger.ReadEvents()
	e, ok := Event{}, true
	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				err = Delete(e.Key)
			case EventPut:
				err = Put(e.Key, e.Value)

			}

		}
	}
	logger.Run()
	return err
}
func InitializePostgresTransactionLog() error {
	var err error
	logger, err = NewPostgresTransactionLogger(PostgresDBParams{
		dbName:   "postgres",
		host:     "localhost",
		user:     "postgres",
		password: "pacan334",
	})
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}
	events, errors := logger.ReadEvents()
	e, ok := Event{}, true
	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				err = Delete(e.Key)
			case EventPut:
				err = Put(e.Key, e.Value)

			}

		}
	}
	logger.Run()
	return err
}
