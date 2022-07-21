package core

import (
	"errors"
	"log"
)

var ErrorNoSuchKey = errors.New("no such key")

type EventType byte

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

type KeyValueStore struct {
	m        map[string]string
	transact TransactionLogger
}

func NewKeyValueStore(transact TransactionLogger) *KeyValueStore {
	return &KeyValueStore{
		m:        make(map[string]string),
		transact: transact,
	}
}
func (store *KeyValueStore) Restore() error {
	var err error
	events, errors := store.transact.ReadEvents()
	count, ok, e := 0, true, Event{}
	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				delete(store.m, e.Key)
				count++
			case EventPut:
				store.m[e.Key] = e.Value
				count++
			}

		}
	}
	log.Printf("%d events replayed\n", count)
	store.transact.Run()
	go func() {
		for err := range store.transact.Err() {
			log.Println(err)
		}
	}()
	return err
}
func (store *KeyValueStore) Get(key string) (string, error) {
	value := store.m[key]
	return value, nil
}
func (store *KeyValueStore) Put(key string, value string) error {
	store.m[key] = value
	store.transact.WritePut(key, value)
	return nil
}
func (store *KeyValueStore) Delete(key string) error {
	delete(store.m, key)
	store.transact.WriteDelete(key)
	return nil
}
