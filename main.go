package main

import (
	"github.com/gorilla/mux"
	"github.com/mvkvl/modbus"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	gateway = "localhost:20108"
)

func main() {

	handler := modbus.NewEncClientHandler(gateway)
	handler.IdleTimeout = 2 * time.Second
	handler.Timeout = 1 * time.Second
	handler.Logger = log.New(os.Stdout, "tcp: ", log.LstdFlags|log.Lmicroseconds)
	defer handler.Close()

	server := &Server{
		client: modbus.NewClient(handler),
	}

	r := mux.NewRouter()
	r.HandleFunc("/c/{slaveId}/{address}", server.ReadCoil).Methods("GET")
	r.HandleFunc("/d/{slaveId}/{address}", server.ReadDiscrete).Methods("GET")
	r.HandleFunc("/i/{slaveId}/{address}", server.ReadInput).Methods("GET")
	r.HandleFunc("/h/{slaveId}/{address}", server.ReadHolding).Methods("GET")

	r.HandleFunc("/c/{slaveId}/{address}", server.WriteCoil).Methods("POST")
	r.HandleFunc("/h/{slaveId}/{address}", server.WriteHolding).Methods("POST")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":80", r))
}
