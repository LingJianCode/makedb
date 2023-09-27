package httpserver

import (
	"fmt"
	"log"
	"makedb/datastore"
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
	err := s.Ds.Put([]byte(vars["key"]), []byte(vars["value"]))
	w.WriteHeader(http.StatusOK)
	if err != nil {
		fmt.Fprintf(w, STATUS_FORMAT, FAIL)
	} else {
		fmt.Fprintf(w, STATUS_FORMAT, SUCCESS)
	}
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// fmt.Println(vars["key"])
	value, err := s.Ds.Get([]byte(vars["key"]))
	w.WriteHeader(http.StatusOK)
	if err != nil {
		fmt.Fprintf(w, STATUS_FORMAT, FAIL)
	} else {
		fmt.Fprintf(w, STATUS_DATA_FORMAT, SUCCESS, vars["key"], string(value))
	}
}

func StartServer() {
	s := NewServer("./data")
	r := mux.NewRouter()
	r.HandleFunc("/", s.HomeHandler).Methods("GET")
	r.HandleFunc("/{key}", s.GetHandler).Methods("GET")
	r.HandleFunc("/{key}/{value}", s.PutHandler).Methods("PUT")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
