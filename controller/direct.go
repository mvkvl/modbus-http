package controller

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mvkvl/modbus"
	"io"
	"log"
	"net/http"
	"strconv"
)

// region - service

type directModbusController struct {
	client modbus.Client
}

func NewDirectModbusClient(client modbus.Client) DirectModbusAPI {
	return &directModbusController{
		client: client,
	}
}

// endregion

// region - controller API

type DirectModbusAPI interface {
	ReadCoil(w http.ResponseWriter, r *http.Request)
	ReadDiscrete(w http.ResponseWriter, r *http.Request)
	ReadInput(w http.ResponseWriter, r *http.Request)
	ReadHolding(w http.ResponseWriter, r *http.Request)
	WriteCoil(w http.ResponseWriter, r *http.Request)
	WriteHolding(w http.ResponseWriter, r *http.Request)
}

func (s *directModbusController) ReadCoil(w http.ResponseWriter, r *http.Request) {
	read(w, r, s.client.ReadCoils)
}
func (s *directModbusController) ReadDiscrete(w http.ResponseWriter, r *http.Request) {
	read(w, r, s.client.ReadDiscreteInputs)
}
func (s *directModbusController) ReadInput(w http.ResponseWriter, r *http.Request) {
	read(w, r, s.client.ReadInputRegisters)
}
func (s *directModbusController) ReadHolding(w http.ResponseWriter, r *http.Request) {
	read(w, r, s.client.ReadHoldingRegisters)
}
func (s *directModbusController) WriteCoil(w http.ResponseWriter, r *http.Request) {
	write(w, r, s.client.WriteSingleCoil, true)
}
func (s *directModbusController) WriteHolding(w http.ResponseWriter, r *http.Request) {
	write(w, r, s.client.WriteSingleRegister, false)
}

// endregion

// region - private methods

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

// endregion
