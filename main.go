package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mvkvl/modbus"
	"github.com/mvkvl/modbus-http/device"
	"github.com/mvkvl/modbus-http/model"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	gateway = "mge:20108"
)

func main() {
	config, err := readConfig("./conf/channels.json")
	if nil != err {
		log.Fatalf("error reading config: %s", err)
		return
	}
	//printConfig(config)
	modbusClient := createModbusClient()
	readRegister(&modbusClient, config, "wb-mge-01:msw-k:temperature")
	readRegister(&modbusClient, config, "wb-mge-01:msw-k:Temperature")
	readRegister(&modbusClient, config, "wb-mge-01:msw-k:humidity")
	readRegister(&modbusClient, config, "wb-mge-01:msw-k:Humidity")
	readRegister(&modbusClient, config, "wb-mge-01:msw-k:noise")
	readRegister(&modbusClient, config, "wb-mge-01:msw-k:CO2")
	readRegister(&modbusClient, config, "wb-mge-01:msw-k:air_quality")
	log.Println("-----------------------------------------------")
	readRegister(&modbusClient, config, "wb-mge-01:msw-b:temperature")
	readRegister(&modbusClient, config, "wb-mge-01:msw-b:Temperature")
	readRegister(&modbusClient, config, "wb-mge-01:msw-b:humidity")
	readRegister(&modbusClient, config, "wb-mge-01:msw-b:Humidity")
	readRegister(&modbusClient, config, "wb-mge-01:msw-b:noise")
	readRegister(&modbusClient, config, "wb-mge-01:msw-b:CO2")
	readRegister(&modbusClient, config, "wb-mge-01:msw-b:air_quality")
}

func readRegister(client *modbus.Client, config *model.Config, reference string) {
	v, t, err := device.ReadFloat(client, config, reference)
	if nil != err {
		log.Fatalf("%s\n", err)
	} else {
		log.Printf("%-11s: %.2f", t, v)
	}
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
	// add back-reference from register to device
	for i := 0; i < len(config.Channels); i++ {
		c := config.Channels[i]
		for j := 0; j < len(c.Devices); j++ {
			d := c.Devices[j]
			for k := 0; k < len(d.Registers); k++ {
				d.Registers[k].Device = &d
			}
		}
	}
	return &config, nil
}
func printConfig(config *model.Config) {
	for _, c := range config.Channels {
		log.Printf("title: %s, conn: %s, mode: %s\n", c.Title, c.Connection, c.Mode)
		for _, d := range c.Devices {
			log.Printf("\t%s:%d\n", d.Title, d.SlaveId)
			for _, r := range d.Registers {
				log.Printf(
					"\t\taddr: %4d, size: %2d, type: %7s, mode: %s, factor: %f, dev: %s\n",
					r.Address, r.Size, r.Type, r.Mode, r.Factor, r.Device.Title)
			}
		}
	}
}
func createModbusClient() modbus.Client {
	handler := modbus.NewEncClientHandler(gateway)
	handler.IdleTimeout = 2 * time.Second
	handler.Timeout = 1 * time.Second
	//handler.Logger = log.New(os.Stdout, "tcp: ", log.LstdFlags|log.Lmicroseconds)
	defer handler.Close()
	client := modbus.NewClient(handler)
	return client
}
func startServer(client *modbus.Client) {

	server := &Server{
		client: *client,
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

//v, err := device.ReadFloatRegister(&modbusClient, &config.Channels[0].Devices[2].Registers[0])
//res, _ := modbusClient.ReadHoldingRegisters(12, 0, 1)
//val := binary.BigEndian.Uint16(res)
//fmt.Printf("% x => %d\n", res, val)
//startServer(&modbusClient)
//eval := goval.NewEvaluator()
