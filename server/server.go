package server

import (
	"encoding/json"
	"fmt"
	"makedb/datastore"
	"makedb/global"
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	Status string `json:"status"`
	Key    string `json:"key,omitempty"`
	Value  string `json:"value,omitempty"`
}

type Server struct {
	Ds *datastore.DataStore
}

func NewServer(path string) *Server {
	ds, err := datastore.NewDataStore(path)
	if err != nil {
		panic("init datastore error.")
	}
	return &Server{Ds: ds}
}

func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello World!")
}

func (s *Server) PutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value := vars["value"]

	if key == "" || value == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: "fail"})
		return
	}

	err := s.Ds.Put([]byte(key), []byte(value))
	if err != nil {
		global.MAKEDB_LOG.Error(fmt.Sprint(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: "fail"})
		return
	}
	global.MAKEDB_LOG.Info(fmt.Sprintf("put key: %s, value: %s", key, value))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Status: "success"})
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: "fail"})
		return
	}

	value, err := s.Ds.Get([]byte(key))
	if err != nil {
		global.MAKEDB_LOG.Error(fmt.Sprintf("get key: %s - %s", key, err.Error()))
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{Status: "fail"})
		return
	}
	global.MAKEDB_LOG.Info(fmt.Sprintf("get key: %s, value: %s", key, string(value)))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Status: "success", Key: key, Value: string(value)})
}

func StartServer() *Server {
	s := NewServer(global.MAKEDB_CONFIG.Server.DataPath)
	r := mux.NewRouter()
	r.HandleFunc("/", s.HomeHandler).Methods("GET")
	r.HandleFunc("/{key}", s.GetHandler).Methods("GET")
	r.HandleFunc("/{key}/{value}", s.PutHandler).Methods("PUT")
	http.Handle("/", r)
	global.MAKEDB_LOG.Info(fmt.Sprintf("server listen on %s", global.MAKEDB_CONFIG.Server.HttpPort))
	go func() {
		global.MAKEDB_LOG.Error(fmt.Sprint(http.ListenAndServe(fmt.Sprintf(":%s", global.MAKEDB_CONFIG.Server.HttpPort), nil)))
	}()
	return s
}
