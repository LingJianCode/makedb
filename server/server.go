package server

import (
	"fmt"
	"makedb/datastore"
	"makedb/global"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	SUCCESS            = "success"
	FAIL               = "FAIL"
	STATUS_FORMAT      = `{"status":"%s"}`
	STATUS_DATA_FORMAT = `{"status":"%s","key":"%s","value":"%s"}`
)

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
	// fmt.Println(vars["key"], vars["value"])
	key := vars["key"]
	value := vars["value"]
	err := s.Ds.Put([]byte(key), []byte(value))
	if err != nil {
		global.MAKEDB_LOG.Error(fmt.Sprint(err))
		fmt.Fprintf(w, STATUS_FORMAT, FAIL)
		return
	}
	global.MAKEDB_LOG.Info(fmt.Sprintf("put key: %s, value: %s", key, value))
	w.WriteHeader(http.StatusOK)
	if err != nil {
		global.MAKEDB_LOG.Error(fmt.Sprint(err))
		fmt.Fprintf(w, STATUS_FORMAT, FAIL)
	} else {
		fmt.Fprintf(w, STATUS_FORMAT, SUCCESS)
	}
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// fmt.Println(vars["key"])
	key := vars["key"]
	value, err := s.Ds.Get([]byte(key))
	if err != nil {
		global.MAKEDB_LOG.Error(fmt.Sprintf("get key: %s - %s", key, err.Error()))
		fmt.Fprintf(w, STATUS_FORMAT, FAIL)
		return
	}
	global.MAKEDB_LOG.Info(fmt.Sprintf("get key: %s, value: %s", key, string(value)))
	w.WriteHeader(http.StatusOK)
	if err != nil {
		global.MAKEDB_LOG.Error(fmt.Sprint(err))
		fmt.Fprintf(w, STATUS_FORMAT, FAIL)
	} else {
		fmt.Fprintf(w, STATUS_DATA_FORMAT, SUCCESS, vars["key"], string(value))
	}
}

func StartServer() {
	s := NewServer(global.MAKEDB_CONFIG.Server.DataPath)
	r := mux.NewRouter()
	r.HandleFunc("/", s.HomeHandler).Methods("GET")
	r.HandleFunc("/{key}", s.GetHandler).Methods("GET")
	r.HandleFunc("/{key}/{value}", s.PutHandler).Methods("PUT")
	http.Handle("/", r)
	global.MAKEDB_LOG.Info(fmt.Sprintf("server listen on %s", global.MAKEDB_CONFIG.Server.HttpPort))
	global.MAKEDB_LOG.Error(fmt.Sprint(http.ListenAndServe(fmt.Sprintf(":%s", global.MAKEDB_CONFIG.Server.HttpPort), nil)))
}
