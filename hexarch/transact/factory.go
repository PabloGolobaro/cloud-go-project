package transact

import (
	"cloud-go-project/hexarch/core"
	"fmt"
	_ "github.com/lib/pq"
)

func NewTransactionLogger(logger string) (core.TransactionLogger, error) {
	switch logger {
	case "file":
		return NewFileTransactionLogger("hexarch/transaction.log")
	case "postgres":
		return NewPostgresTransactionLogger(PostgresDBParams{
			dbName:   "postgres",
			host:     "localhost",
			user:     "postgres",
			password: "pacan334",
		})
	case "":
		return nil, fmt.Errorf("Transation logger type not defined")
	default:
		return nil, fmt.Errorf("no such transaction logger %s", logger)
	}

}
