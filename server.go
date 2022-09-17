package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mvkvl/modbus"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	client modbus.Client
}

func (s *Server) ReadCoil(w http.ResponseWriter, r *http.Request) {
	read(w, r, s.client.ReadCoils)
}
func (s *Server) ReadDiscrete(w http.ResponseWriter, r *http.Request) {
	read(w, r, s.client.ReadDiscreteInputs)
}
func (s *Server) ReadInput(w http.ResponseWriter, r *http.Request) {
	read(w, r, s.client.ReadInputRegisters)
}
func (s *Server) ReadHolding(w http.ResponseWriter, r *http.Request) {
	read(w, r, s.client.ReadHoldingRegisters)
}

func (s *Server) WriteCoil(w http.ResponseWriter, r *http.Request) {
	write(w, r, s.client.WriteSingleCoil, true)
}
func (s *Server) WriteHolding(w http.ResponseWriter, r *http.Request) {
	write(w, r, s.client.WriteSingleRegister, false)
}

func read(w http.ResponseWriter, r *http.Request, reader func(uint8, uint16, uint16) ([]byte, error)) {
	vars := mux.Vars(r)
	slaveId, _ := strconv.Atoi(vars["slaveId"])
	address, _ := strconv.Atoi(vars["address"])
	result, err := reader(uint8(slaveId), uint16(address), 1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}
	w.Write([]byte(fmt.Sprintf("% x\n", result)))
}
func write(w http.ResponseWriter, r *http.Request, writer func(uint8, uint16, uint16) ([]byte, error), coil bool) {
	vars := mux.Vars(r)
	slaveId, _ := strconv.Atoi(vars["slaveId"])
	address, _ := strconv.Atoi(vars["address"])
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var v int
	if len(body) > 0 {
		v, _ = strconv.Atoi(fmt.Sprintf("%c", body[:]))
	} else {
		v = 0
	}
	if v > 0 && coil {
		v = 0xFF00
	}
	_, err = writer(uint8(slaveId), uint16(address), uint16(v))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
