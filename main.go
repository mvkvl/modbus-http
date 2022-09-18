package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mvkvl/modbus"
	"github.com/mvkvl/modbus-http/model"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	gateway = "localhost:20108"
)

func main() {
	config, err := readConfig("./conf/channels.json")
	if nil != err {
		log.Fatalf("error reading config: %s", err)
		return
	}
	printConfig(config)
}

func readConfig(path string) (*model.Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config model.Config
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
func printConfig(config *model.Config) {
	for _, c := range config.Channels {
		log.Printf("title: %s, conn: %s, mode: %s\n", c.Title, c.Connection, c.Mode)
		for _, d := range c.Devices {
			log.Printf("\t%s:%d\n", d.Title, d.SlaveId)
			for _, r := range d.Registers {
				log.Printf("\t\taddr: %4d, size: %2d, type: %7s, mode: %s\n", r.Address, r.Size, r.Type, r.Mode)
			}
		}
	}
}

func startServer() {
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
