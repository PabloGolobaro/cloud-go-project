package main

import (
	"cloud-go-project/cmd"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	err := cmd.InitializePostgresTransactionLog()
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/v1/{key}", cmd.KeyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", cmd.KeyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", cmd.KeyValueDeleteHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}
