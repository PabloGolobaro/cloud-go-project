package frontend

import (
	"cloud-go-project/hexarch/core"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type FrontEnd interface {
	Start(kv *core.KeyValueStore) error
}

func NewFrontend(frontend string) (FrontEnd, error) {
	switch frontend {
	case "rest":
		return &restFrontEnd{}, nil
	case "grpc":
		return &grpcFrontEnd{}, nil
	case "":
		return nil, fmt.Errorf("Frontend type not defined")
	default:
		return nil, fmt.Errorf("no such frontend %s", frontend)
	}
}

type restFrontEnd struct {
	store *core.KeyValueStore
}

func (f *restFrontEnd) Start(store *core.KeyValueStore) error {
	f.store = store
	r := mux.NewRouter()
	r.HandleFunc("/v1/{key}", f.keyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", f.keyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", f.keyValueDeleteHandler).Methods("DELETE")

	return http.ListenAndServe(":8080", r)

}

func (f *restFrontEnd) keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	err := f.store.Delete(key)
	if errors.Is(err, core.ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Deleted"))

}
func (f *restFrontEnd) keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = f.store.Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

}
func (f *restFrontEnd) keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value, err := f.store.Get(key)
	if errors.Is(err, core.ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(value))

}
