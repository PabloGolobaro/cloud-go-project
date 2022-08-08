package main

import (
	"cloud-go-project/hexarch/core"
	"cloud-go-project/hexarch/frontend"
	"cloud-go-project/hexarch/transact"
	"fmt"
	"log"
	"os"
)

func main() {

	//*********************************************************
	os.Setenv("TLOG_TYPE", os.Args[1])
	fmt.Println(os.Args[1])
	os.Setenv("FRONTEND_TYPE", os.Args[2])
	fmt.Println(os.Args[2])
	logger, err := transact.NewTransactionLogger(os.Getenv("TLOG_TYPE"))
	if err != nil {
		log.Fatal(err)
	}
	store := core.NewKeyValueStore(logger)
	store.Restore()
	fe, err := frontend.NewFrontend(os.Getenv("FRONTEND_TYPE"))
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(fe.Start(store))
}
