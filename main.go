package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/jeanphorn/log4go"
	"github.com/mvkvl/modbus"
	"github.com/mvkvl/modbus-http/controller"
	"github.com/mvkvl/modbus-http/model"
	"github.com/mvkvl/modbus-http/service"
	"net/http"
	"os"
	"time"
)

const (
	gateway = "localhost:20108"
)

func main() {

	log.LoadConfiguration("./conf/logger.json")
	defer log.Close()

	//config, err := readConfig("./conf/channels.json")
	//if nil != err {
	//	log.Error("error reading config: %s", err)
	//	return
	//}

	//printConfig(config)
	//pollerTest(config)
	//startServer(config)

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
		log.Info("title: %s, conn: %s, mode: %s, cpause: %d, rpause: %d",
			c.Title, c.Connection, c.Mode, c.GetCyclePause(), c.GetRegisterPause())
		for _, d := range c.Devices {
			log.Info("\t%s:%d", d.Title, d.SlaveId)
			for _, r := range d.Registers {
				log.Info("\t\taddr: %4d, size: %2d, type: %7s, mode: %s", r.Address, r.Size, r.Type, r.Mode)
			}
		}
	}
}

func pollerTest(config *model.Config) {
	//poller := service.NewModbusPoller(modbusHandlerFactory, config)
	//log.Println("poller start 1")
	//poller.Start()
	//time.Sleep(35 * time.Second)

	//log.Println("poller stop 1")
	//poller.Stop()
	//time.Sleep(10 * time.Second)
	//log.Println("poller start 2")
	//poller.Start()
	//time.Sleep(15 * time.Second)
	//log.Println("poller stop 2")
	//poller.Stop()
	//time.Sleep(10 * time.Second)

	//v, e := poller.Read("wb-mge-01:msw-b:noise")
	//if nil != e {
	//	fmt.Printf("MSW-B NOISE: %s\n", e)
	//} else {
	//	fmt.Printf("MSW-B NOISE: % x\n", v)
	//}
	//poller.Start()
	//defer poller.Stop()

}

func startServer(config *model.Config) {

	poller := service.NewModbusPoller(modbusHandlerFactory, config)
	poller.Start()

	directModbusService := controller.NewDirectModbusClient(modbus.NewClient(modbusHandlerFactory(gateway, model.ENC)))
	cachedModbusService := controller.NewCachedModbusClient(poller)

	r := mux.NewRouter()
	r.HandleFunc("/direct/c/{slaveId}/{address}", directModbusService.ReadCoil).Methods("GET")
	r.HandleFunc("/direct/d/{slaveId}/{address}", directModbusService.ReadDiscrete).Methods("GET")
	r.HandleFunc("/direct/i/{slaveId}/{address}", directModbusService.ReadInput).Methods("GET")
	r.HandleFunc("/direct/h/{slaveId}/{address}", directModbusService.ReadHolding).Methods("GET")

	r.HandleFunc("/direct/c/{slaveId}/{address}", directModbusService.WriteCoil).Methods("POST")
	r.HandleFunc("/direct/h/{slaveId}/{address}", directModbusService.WriteHolding).Methods("POST")

	r.HandleFunc("/cached/start", cachedModbusService.Start).Methods("POST")
	r.HandleFunc("/cached/stop", cachedModbusService.Stop).Methods("POST")
	r.HandleFunc("/cached/{metric}", cachedModbusService.Get).Methods("GET")

	// Bind to a port and pass our router in
	log.Error(http.ListenAndServe(":8080", r))
}

func modbusHandlerFactory(connection string, mode model.Mode) modbus.ClientHandler {
	switch mode {
	case model.ENC:
		_handler := modbus.NewEncClientHandler(connection)
		_handler.IdleTimeout = 2 * time.Second
		_handler.Timeout = 1 * time.Second
		//_handler.Logger = log.New(os.Stdout, fmt.Sprintf("[%s]: ", connection), log.LstdFlags|log.Lmicroseconds)
		return _handler
	case model.TCP:
		_handler := modbus.NewTCPClientHandler(connection)
		//_handler.Logger = log.New(os.Stdout, fmt.Sprintf("[%s]: ", connection), log.LstdFlags|log.Lmicroseconds)
		return _handler
	case model.RTU:
		_handler := modbus.NewRTUClientHandler(connection)
		//_handler.Logger = log.New(os.Stdout, fmt.Sprintf("[%s]: ", connection), log.LstdFlags|log.Lmicroseconds)
		return _handler
	}
	return nil
}
